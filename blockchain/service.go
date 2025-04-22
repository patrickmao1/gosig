package blockchain

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/crypto/sha3"
	"google.golang.org/protobuf/proto"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

var log = logrus.New()

type Service struct {
	db     *DB
	txPool *TxPool
	nw     *Network

	cfg           *NodeConfig
	genesis       *GenesisConfig
	roundInterval time.Duration

	// current round states
	rmu       sync.RWMutex
	round     uint32
	seed      []byte
	rseed     []byte                          // round concatenated with seed
	proposals map[uint32]*types.BlockProposal // key is val index
	head      *types.BlockHeader

	// gossip messages
	outMsgs *OutboundMsgBuffer
	inMsgs  *InboundMsgBuffer

	// catch up routine state
	catchUpRunning atomic.Bool
}

func NewService(cfg *NodeConfig, genesis *GenesisConfig, db *DB, txPool *TxPool) *Service {
	roundInterval := time.Duration(cfg.ProposalStageDurationMs + cfg.ProposalStageDurationMs)

	outbound := NewOutboundMsgBuffer(cfg.MyPrivKey(), cfg.MyValidatorIndex())
	inbound := NewInboundMsgBuffer(cfg.Validators)

	nw := NewNetwork(
		outbound,
		inbound,
		cfg.GossipDegree,
		time.Duration(cfg.GossipIntervalMs)*time.Millisecond,
		cfg.Validators,
	)
	return &Service{
		cfg:           cfg,
		genesis:       genesis,
		roundInterval: roundInterval,
		db:            db,
		txPool:        txPool,
		nw:            nw,
		outMsgs:       outbound,
		inMsgs:        inbound,
		proposals:     make(map[uint32]*types.BlockProposal),
	}
}

func (s *Service) Start() error {
	log.Infoln("Starting blockchain service")

	go s.nw.StartGossip()

	nextRoundTime := s.nextRoundTime()
	t := time.NewTicker(time.Until(nextRoundTime))
	for {
		<-t.C
		proposalStageEnd := time.After(s.cfg.ProposalStageDuration())

		/*
		 * Stage 1: VRF-based Proposal
		 */

		// initialize round number and seed
		err := s.initRoundState()
		if err != nil {
			log.Errorln(err)
		}

		// check if I'm a proposer and propose
		err = s.proposeIfChosen()
		if err != nil {
			log.Errorln(err)
		}

		// block until the end of the proposal stage
		<-proposalStageEnd

		/*
		 * Stage 2: BFT agreement
		 */

		decidedBlock, err := s.decideBlock()
		if err != nil {
			log.Errorln("cannot decide block", err)
			continue
		}

		if decidedBlock != nil {
			err = s.prepare(decidedBlock.Hash())
			if err != nil {
				log.Errorln("failed to prepare", err)
				continue
			}
		} else {
			log.Infof("no block decided for round %d", s.round)
		}
		// TODO: Remove proposal messages from msg buffer.

		// reset round timer
		nextRoundTime = s.nextRoundTime()
		t.Reset(time.Until(nextRoundTime))
	}
}

func (s *Service) initRoundState() error {
	s.rmu.Lock()

	s.round++
	s.proposals = nil
	head, err := s.db.GetHeadBlock()
	if err != nil {
		return err
	}
	s.head = head

	// compute the seed for the round
	seed := sha3.Sum256(head.ProposerProof)
	s.seed = seed[:]
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, s.round)
	s.rseed = slices.Concat(b, s.seed)

	s.rmu.Unlock()
	return nil
}

func (s *Service) proposeIfChosen() error {
	// determine if I should propose
	proposalScore, proposerProof := crypto.VRF(s.cfg.MyPrivKey(), s.rseed) // L^r = SIG_{l^h}(r^h,Q^{h-1})
	if proposalScore >= s.cfg.ProposalThreshold {
		return nil
	}

	// build a block
	headHash := s.head.Hash()
	header := &types.BlockHeader{
		Height:        s.head.Height + 1,
		ParentHash:    headHash,
		ProposerProof: proposerProof,
	}
	tcBlock, pCert, err := s.db.GetTcState()
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		return err
	}

	var txs []*types.SignedTransaction
	if errors.Is(err, leveldb.ErrNotFound) {
		// no tc block found, we propose a new block
		txs = s.txPool.Dump()
		header.TxRoot = computeTxRoot(txs)
		headTcCert, err := s.db.GetTcCert(headHash)
		if err != nil {
			return err
		}
		header.ProposalCert = headTcCert
	} else { // TC block found.
		// Sanity check. If these aren't the same then I wouldn't have TC'd that block
		if !bytes.Equal(tcBlock.ParentHash, s.head.Hash()) {
			log.Panicf("inconsistent parent hashes! tcBlock %v, headBlock %v", tcBlock, s.head)
		}
		// I have previously TC'd a block, but I never committed it. Re-propose the block.
		txs, header.TxRoot, err = s.db.GetBlockTxs(tcBlock.Hash())
		if err != nil {
			return err
		}
		header.ProposalCert = pCert
	}
	if len(txs) == 0 {
		return nil
	}

	// build & sign proposal
	proposal := &types.BlockProposal{
		Round:       s.round,
		BlockHeader: header,
	}
	msg := &types.Message{
		Message:  &types.Message_Proposal{Proposal: proposal},
		Deadline: s.roundProposalEndTime().UnixMilli(),
	}

	// broadcast
	s.outMsgs.Put(s.proposalKey(s.cfg.MyValidatorIndex()), msg)

	return nil
}

func (s *Service) decideBlock() (*types.BlockHeader, error) {
	// no proposals received, leaving the round empty
	if len(s.proposals) == 0 {
		return nil, nil
	}

	// find the proposal with the lowest score
	minProposalVi := uint32(0)
	for vi, proposal := range s.proposals {
		if proposal.Score() < s.proposals[minProposalVi].Score() {
			minProposalVi = vi
		}
	}
	minProposal := s.proposals[minProposalVi]

	tcBlock, pCert, err := s.db.GetTcState()
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		return nil, err
	}
	// if I have no TC'd block, then I can safely use the best proposal I see
	if errors.Is(err, leveldb.ErrNotFound) {
		return minProposal.BlockHeader, nil
	}

	// if I do have a TC'd block...

	// the incoming proposal's cert is either the prior block's TC cert, or a prepare cert of a
	// block TC'd by the proposer in the previous round. If the new proposal has a cert round >
	// my tc_round, then it implies that their view of the blockchain is newer than mine. Thus,
	// I can just use their proposal
	if minProposal.BlockHeader.ProposalCert.Round > pCert.Round {
		return minProposal.BlockHeader, nil
	}

	// if I do have a TC'd block and if no block has a higher rounded cert than my TC'd block, then
	// either the proposers are lagging or they also want to re-propose the block I TC'd. I need to
	// check every proposal and find if it's the latter case, then support that TC block. Otherwise,
	// I vote for nothing and wait for either a higher round block supported by the majority or for
	// my turn to propose the TC'd block.
	var proposalsWithTcBlock []*types.BlockProposal
	for _, proposal := range s.proposals {
		if bytes.Equal(proposal.BlockHeader.Hash(), tcBlock.Hash()) &&
			proposal.BlockHeader.ProposalCert.Round >= pCert.Round {
			proposalsWithTcBlock = append(proposalsWithTcBlock, proposal)
		}
	}
	// no one has the same tc block as me, I do nothing.
	if len(proposalsWithTcBlock) == 0 {
		return nil, nil
	}
	// otherwise I need to find the proposal with the min score so that nodes like me all agree on
	// the same proposal to prepare
	minProposal = proposalsWithTcBlock[0]
	for _, proposal := range proposalsWithTcBlock {
		if proposal.Score() < minProposal.Score() {
			minProposal = proposal
		}
	}
	return minProposal.BlockHeader, nil
}

func (s *Service) prepare(blockHash []byte) error {
	prep := &types.Prepare{BlockHash: blockHash}
	cert, err := s.signNewCertificate(prep)
	if err != nil {
		return err
	}
	prepCert := &types.PrepareCertificate{
		Msg:  prep,
		Cert: cert,
	}
	signedMsg := &types.Message{
		Message:  &types.Message_Prepare{Prepare: prepCert},
		Deadline: s.nextRoundTime().UnixMilli(),
	}
	s.outMsgs.Put(s.prepareKey(), signedMsg)
	return nil
}

func (s *Service) tentativeCommit(blockHash []byte) error {
	// TODO check transactions before TC

	tc := &types.TentativeCommit{BlockHash: blockHash}
	cert, err := s.signNewCertificate(tc)
	if err != nil {
		return err
	}
	tcWithCert := &types.TentativeCommitCertificate{
		Msg:  tc,
		Cert: cert,
	}
	signedMsg := &types.Message{
		Message:  &types.Message_Tc{Tc: tcWithCert},
		Deadline: s.nextRoundTime().UnixMilli(),
	}
	s.outMsgs.Put(s.tcKey(), signedMsg)

	err = s.db.PutTcCert(blockHash, cert)
	if err != nil {
		return err
	}
	// We also need to persist our state of what I have TC'd so that later if this block isn't
	// committed we can know its block hash. This state is cleared upon commit of this block.
	return s.db.Put(tcBlockKey, blockHash, nil)
}

func (s *Service) commit(blockHash []byte) error {
	log.Infof("committing block %x", blockHash)

	// TODO do commit stuff

	return s.db.Delete(tcBlockKey, nil)
}

func (s *Service) signNewCertificate(msg proto.Message) (*types.Certificate, error) {
	bs, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	sig := crypto.SignBytes(s.cfg.privKey, bs)
	sigCounts := make([]uint32, len(s.cfg.Validators))
	sigCounts[s.cfg.MyValidatorIndex()] = 1
	return &types.Certificate{
		AggSig:    sig,
		SigCounts: sigCounts,
		Round:     s.round,
	}, nil
}

func (s *Service) verifyCertificate(msg proto.Message, cert *types.Certificate) error {
	bs, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	pass := crypto.VerifyAggSigBytes(s.cfg.Validators.PubKeys(), cert.SigCounts, bs, cert.AggSig)
	if !pass {
		return fmt.Errorf("invalid cert: %+v %+v", msg, cert)
	}
	return nil
}

func (s *Service) runCatchUpRoutine() {
	if s.catchUpRunning.Load() {
		return
	}
	s.catchUpRunning.Store(true)
	log.Infof("running catch up routine")

	// TODO

	s.catchUpRunning.Store(false)
}

func (s *Service) getCurrentRound() int {
	return int(time.Since(s.genesis.GenesisTime()) / s.roundInterval)
}

func (s *Service) roundProposalEndTime() time.Time {
	return s.nextRoundTime().Add(-time.Duration(s.cfg.AgreementStateDurationMs))
}

func (s *Service) nextRoundTime() time.Time {
	round := s.getCurrentRound()
	sinceGenesis := time.Duration(round+1) * s.roundInterval
	return s.genesis.GenesisTime().Add(sinceGenesis)
}

func computeTxRoot(signedTxs []*types.SignedTransaction) []byte {
	var txHashes []byte
	for _, tx := range signedTxs {
		txHashes = append(txHashes, tx.Tx.Hash()...)
	}
	root := sha3.Sum256(txHashes)
	return root[:]
}

func (s *Service) proposalKey(vi uint32) string { return fmt.Sprintf("proposal-%d-%d", s.round, vi) }
func (s *Service) prepareKey() string           { return fmt.Sprintf("prepare-%d", s.round) }
func (s *Service) tcKey() string                { return fmt.Sprintf("tc-%d", s.round) }
