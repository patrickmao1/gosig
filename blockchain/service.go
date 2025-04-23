package blockchain

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/crypto/blake2b"
	"google.golang.org/protobuf/proto"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

func init() {
	//log.SetReportCaller(true)
}

type Service struct {
	db     *DB
	txPool *TxPool
	nw     *Network

	// initial configs
	cfg              *NodeConfig
	genesis          *GenesisConfig
	genesisBlockHash []byte

	// current round states
	rmu       sync.RWMutex
	round     *atomic.Uint32
	rseed     []byte                          // round concatenated with seed
	proposals map[uint32]*types.BlockProposal // key is val index
	head      *types.BlockHeader

	// gossip messages
	outMsgs *OutboundMsgBuffer
	inMsgs  *InboundMsgBuffer

	// catch up routine state
	catchUpRunning *atomic.Bool
}

func NewService(cfg *NodeConfig, genesis *GenesisConfig, db *DB, txPool *TxPool) *Service {

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
		cfg:            cfg,
		genesis:        genesis,
		db:             db,
		txPool:         txPool,
		nw:             nw,
		outMsgs:        outbound,
		inMsgs:         inbound,
		proposals:      make(map[uint32]*types.BlockProposal),
		round:          new(atomic.Uint32),
		catchUpRunning: new(atomic.Bool),
	}
}

func (s *Service) Start() {
	log.Infoln("Starting blockchain service")

	go s.nw.StartGossip()
	go s.startProcessingMsgs()

	nextRoundTime := s.nextRoundTime()
	log.Infof("genesis time %s, current round %d, now %s, next round time %s",
		s.genesis.GenesisTime(), s.getCurrentRound(), time.Now(), nextRoundTime)

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
			log.Fatal(err)
		}

		// check if I'm a proposer and propose
		err = s.proposeIfChosen()
		if err != nil {
			log.Errorln(err)
			// we don't continue the loop here as we still want to participate the next stage
		}

		// block until the end of the proposal stage
		<-proposalStageEnd
		// reset round timer
		nextRoundTime = s.nextRoundTime()
		t.Reset(time.Until(nextRoundTime))

		/*
		 * Stage 2: BFT agreement
		 */

		decidedBlock, err := s.decideBlock()
		if err != nil {
			log.Errorln("cannot decide block", err)
			continue
		}

		if decidedBlock != nil {
			err := s.db.PutBlockHeader(decidedBlock)
			if err != nil {
				log.Errorln("failed to put decided block", err)
				continue
			}
			err = s.prepare(decidedBlock.Hash())
			if err != nil {
				log.Errorln("failed to prepare", err)
				continue
			}
		}
	}
}

func (s *Service) initRoundState() error {
	s.rmu.Lock()
	defer s.rmu.Unlock()

	s.proposals = make(map[uint32]*types.BlockProposal)
	head, err := s.db.GetHeadBlock()
	if err != nil && errors.Is(err, leveldb.ErrNotFound) {
		log.Infof("no head block found, performing genesis")
		genesisBlock := NewGenesisBlock(s.genesis.InitialSeed)
		err := s.db.PutBlockHeader(genesisBlock)
		if err != nil {
			return err
		}
		log.Infof("genesis block saved to db: %s", genesisBlock)
		s.genesisBlockHash = genesisBlock.Hash()
		err = s.db.PutHeadBlock(s.genesisBlockHash)
		if err != nil {
			return err
		}
		head = genesisBlock
	} else if err != nil {
		return err
	}
	s.round.Store(s.getCurrentRound())
	s.head = head

	// compute the seed for the round
	seed := blake2b.Sum256(head.ProposerProof)
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, s.round.Load())
	s.rseed = slices.Concat(b, seed[:])

	log.Infof("new round %d, head %x, rseed %x", s.round.Load(), s.head.Hash(), s.rseed)

	return nil
}

func (s *Service) proposeIfChosen() error {
	// determine if I should propose
	proposalScore, proposerProof := crypto.VRF(s.cfg.MyPrivKey(), s.rseed) // L^r = SIG_{l^h}(r^h,Q^{h-1})
	log.WithField("round", s.round.Load()).
		Debugf("my proposer score %d, threshold %d", proposalScore, s.cfg.ProposalThreshold)
	if proposalScore >= s.cfg.ProposalThreshold {
		return nil
	}
	log.WithField("round", s.round.Load()).Infof("I am a proposer!")
	log.Debugf("RNG info: pubkey %x.., rseed %x, proof %x..",
		s.cfg.Validators.PubKeys()[s.cfg.MyValidatorIndex()][:8], s.rseed, proposerProof[:8])

	// build a block
	headHash := s.head.Hash()
	header := &types.BlockHeader{
		Height:        s.head.Height + 1,
		ParentHash:    headHash,
		ProposerProof: proposerProof,
	}
	tcBlock, pCert, err := s.db.GetTcState()
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		return fmt.Errorf("leveldb error: %s", err.Error())
	}

	var txs []*types.SignedTransaction
	var cert *types.Certificate

	if errors.Is(err, leveldb.ErrNotFound) {
		txs = s.txPool.Dump()
		if len(txs) == 0 {
			log.Infof("skip proposing: no tx in pool")
			return nil
		}

		header.TxRoot = computeTxRoot(txs)
		log.Infof("proposing a new block with %d txs, txRoot %x", len(txs), header.TxRoot)

		if !bytes.Equal(s.genesisBlockHash, headHash) {
			// skip this if I'm still in genesis mode because the genesis block does not have a tc cert
			headTcCert, err := s.db.GetTcCert(headHash)
			if err != nil {
				return fmt.Errorf("failed to GetTcCert for head block %x: %s", headHash, err.Error())
			}
			cert = headTcCert
		}
	} else { // TC block found.
		log.Infof("tc block found, proposing tc block %s", tcBlock.Hash())
		// Sanity check. If these aren't the same then I wouldn't have TC'd that block
		if !bytes.Equal(tcBlock.ParentHash, s.head.Hash()) {
			log.Panicf("inconsistent parent hashes! tcBlock %v, headBlock %v", tcBlock, s.head)
		}
		// I have previously TC'd a block, but I never committed it. Re-propose the block.
		txs, header.TxRoot, err = s.db.GetBlockTxs(tcBlock.Hash())
		if err != nil {
			return err
		}
		cert = pCert
	}
	if len(txs) == 0 {
		return nil
	}

	// build & sign proposal
	proposal := &types.BlockProposal{
		Round:         s.round.Load(),
		BlockHeader:   header,
		Cert:          cert,
		ProposerIndex: s.cfg.MyValidatorIndex(),
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
	s.rmu.RLock()
	defer s.rmu.RUnlock()

	// no proposals received, leaving the round empty
	if len(s.proposals) == 0 {
		log.Infof("decideBlock: no block proposals for round %d", s.round.Load())
		return nil, nil
	}

	log.WithField("round", s.round.Load()).
		Infof("decide block: num proposals: %d", len(s.proposals))

	// find the proposal with the lowest score
	var minProposal *types.BlockProposal
	for _, proposal := range s.proposals {
		if minProposal == nil || proposal.Score() < minProposal.Score() {
			minProposal = proposal
		}
	}

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
	if minProposal.Cert.Round > pCert.Round {
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
			proposal.Cert.Round >= pCert.Round {
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
	err := s.db.PutHeadBlock(blockHash)
	if err != nil {
		return err
	}

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
		Round:     s.round.Load(),
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

func (s *Service) getCurrentRound() uint32 {
	return uint32(time.Since(s.genesis.GenesisTime()) / s.cfg.RoundDuration())
}

func (s *Service) roundProposalEndTime() time.Time {
	return s.nextRoundTime().Add(-s.cfg.AgreementStateDuration())
}

func (s *Service) nextRoundTime() time.Time {
	round := s.getCurrentRound()
	sinceGenesis := time.Duration(round+1) * s.cfg.RoundDuration()
	return s.genesis.GenesisTime().Add(sinceGenesis)
}

func computeTxRoot(signedTxs []*types.SignedTransaction) []byte {
	var txHashes []byte
	for _, tx := range signedTxs {
		txHashes = append(txHashes, tx.Tx.Hash()...)
	}
	root := blake2b.Sum256(txHashes)
	return root[:]
}

func (s *Service) proposalKey(vi uint32) string {
	return fmt.Sprintf("proposal-%d-%d", s.round.Load(), vi)
}
func (s *Service) prepareKey() string { return fmt.Sprintf("prepare-%d", s.round.Load()) }
func (s *Service) tcKey() string      { return fmt.Sprintf("tc-%d", s.round.Load()) }
