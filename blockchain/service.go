package blockchain

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	"github.com/patrickmao1/gosig/utils"
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
	rmu         sync.RWMutex
	round       *atomic.Uint32
	rseed       []byte // round concatenated with seed
	minProposal *types.BlockProposal
	head        *types.BlockHeader

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
		round:          new(atomic.Uint32),
		catchUpRunning: new(atomic.Bool),
	}
}

func (s *Service) Start() {
	log.Infoln("Starting blockchain service")

	go s.nw.StartGossip()
	go s.startProcessingMsgs()
	go s.StartRPC()

	nextRoundTime := s.nextRoundTime()
	log.Infof("genesis time %s, current round %d, now %s, next round time %s, proposal threshold %.2f",
		s.genesis.GenesisTime(), s.getCurrentRound(), time.Now(), nextRoundTime, s.cfg.ProposalThresholdPerc())

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
		txs, err := s.proposeIfChosen()
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

		if s.minProposal == nil {
			continue
		}
		err = s.db.PutBlock(s.minProposal.BlockHeader)
		if err != nil {
			log.Errorln("failed to put decided block", err)
			continue
		}
		err = s.prepare(s.minProposal.BlockHeader.Hash())
		if err != nil {
			log.Errorln("failed to prepare", err)
			continue
		}
		// I won. I should send out the transaction hashes (body of the block)
		if s.minProposal.ProposerIndex == s.cfg.MyValidatorIndex() {
			s.broadcastTxHashes(txs)
		}
	}
}

func (s *Service) initRoundState() error {
	s.rmu.Lock()
	defer s.rmu.Unlock()

	head, err := s.db.GetHeadBlock()
	if err != nil && errors.Is(err, leveldb.ErrNotFound) {
		log.Infof("no head block found, performing genesis")
		head, err = s.doGenesis()
	} else if err != nil {
		return err
	}

	s.round.Store(s.getCurrentRound())
	s.head = head
	s.minProposal = nil

	s.txPool.CleanRecentlyDeleted()

	// compute the seed for the round
	seed := blake2b.Sum256(head.ProposerProof)
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, s.round.Load())
	s.rseed = slices.Concat(b, seed[:])

	log.Infof("new round %d, head %x, rseed %x", s.round.Load(), s.head.Hash(), s.rseed)

	return nil
}

func (s *Service) doGenesis() (*types.BlockHeader, error) {
	genesisBlock := NewGenesisBlock(s.genesis.InitialSeed)
	err := s.db.PutBlock(genesisBlock)
	if err != nil {
		return nil, err
	}
	log.Infof("genesis block saved to db: %s", genesisBlock.ToString())
	s.genesisBlockHash = genesisBlock.Hash()
	err = s.db.PutHeadBlock(s.genesisBlockHash)
	if err != nil {
		return nil, err
	}
	return genesisBlock, nil
}

func (s *Service) proposeIfChosen() ([]*types.SignedTransaction, error) {
	txs := s.txPool.List(35_000)
	if len(txs) == 0 {
		log.Debugf("skip proposing: no tx in pool")
		return nil, nil
	}
	log.Infof("round %d tx in pool: %d", s.round.Load(), len(txs))

	// determine if I should propose
	proposalScore, proposerProof := crypto.VRF(s.cfg.MyPrivKey(), s.rseed) // L^r = SIG_{l^h}(r^h,Q^{h-1})
	log.WithField("round", s.round.Load()).
		Debugf("my proposer score %d, threshold %d", proposalScore, s.cfg.ProposalThreshold)
	if proposalScore >= s.cfg.ProposalThreshold {
		return nil, nil
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

	var cert *types.Certificate
	header.TxRoot = computeTxRoot(txs)
	log.Infof("proposing a new block with %d txs, txRoot %x", len(txs), header.TxRoot)

	// If I'm still in genesis mode then cert is nil because the genesis block does not have a tc cert
	if !bytes.Equal(s.genesisBlockHash, headHash) {
		headTcCert, err := s.db.GetTcCert(headHash)
		if err != nil {
			return nil, fmt.Errorf("failed to GetTcCert for head block %x: %s", headHash, err.Error())
		}
		cert = headTcCert
	}

	// build & sign proposal
	proposal := &types.BlockProposal{
		Round:         s.round.Load(),
		BlockHeader:   header,
		Cert:          cert,
		ProposerIndex: s.cfg.MyValidatorIndex(),
	}
	msg := &types.Message{
		Message: &types.Message_Proposal{Proposal: proposal},
	}

	// broadcast
	s.outMsgs.Put(s.proposalKey(), msg, s.proposalStageEndTime())

	log.Infof("proposed block %s", header.ToString())

	return txs, nil
}

func (s *Service) prepare(blockHash []byte) error {
	defer utils.LogExecTime(time.Now(), "prepare")

	log.Infof("prepare: block %x", blockHash)

	prep := &types.Prepare{BlockHash: blockHash}
	cert, err := s.signNewCert(prep)
	if err != nil {
		return err
	}
	prepCert := &types.PrepareCertificate{
		Msg:  prep,
		Cert: cert,
	}
	msg := &types.Message{
		Message: &types.Message_Prepare{Prepare: prepCert},
	}
	s.outMsgs.Put(s.prepareKey(), msg, s.prepareGossipEndTime(), s.nextRoundTime())
	return nil
}

func (s *Service) tentativeCommit(outMsgs *OutboundMsgBuffer, blockHash []byte) error {
	defer utils.LogExecTime(time.Now(), "tentativeCommit")

	txHashes, err := s.db.GetTxHashes(blockHash)
	if err != nil {
		return fmt.Errorf("cannot tc block %x..: txHashes not found", blockHash[:8])
	}
	log.Infof("tentativeCommit: block %x", blockHash)

	txs, missingTxs := s.txPool.BatchGet(txHashes.TxHashes)
	if len(missingTxs) > 0 {
		received, err := s.nw.QueryTxs(missingTxs)
		if err != nil {
			return err
		}
		err = s.txPool.AddTransactions(received)
		if err != nil {
			return err
		}
		pass := checkTxs(received)
		if !pass {
			return fmt.Errorf("cannot tc block %x..: tx check failed", blockHash[:8])
		}
		txs = append(txs, received...)
	}

	tc := &types.TentativeCommit{BlockHash: blockHash}
	cert, err := s.signNewCert(tc)
	if err != nil {
		return err
	}
	tcWithCert := &types.TentativeCommitCertificate{
		Msg:  tc,
		Cert: cert,
	}
	signedMsg := &types.Message{
		Message: &types.Message_Tc{Tc: tcWithCert},
	}
	outMsgs.Put(s.tcKey(), signedMsg, s.nextRoundTime())

	return s.db.PutTcCert(blockHash, cert)
}

func checkTxs(txs []*types.SignedTransaction) bool {
	defer utils.LogExecTime(time.Now(), "checkTxs")
	for _, tx := range txs {
		bs, err := proto.Marshal(tx.Tx)
		if err != nil {
			log.Error(err)
			return false
		}
		pass := crypto.VerifySigBytes(tx.Tx.From, bs, tx.Sig)
		if !pass {
			log.Errorf("tx %s invalid signature", tx.ToString())
			return false
		}
	}
	return true
}

func (s *Service) commit(blockHash []byte) error {
	defer utils.LogExecTime(time.Now(), "commit")

	s.rmu.Lock()
	defer s.rmu.Unlock()

	head, err := s.db.Get(headKey, nil)
	if err != nil {
		return fmt.Errorf("commit: failed to get head block: %s", err.Error())
	}
	if bytes.Equal(head, blockHash) {
		return nil
	}
	log.Infof("commit block %x", blockHash)

	err = s.db.PutHeadBlock(blockHash)
	if err != nil {
		return err
	}

	var txHashes *types.TransactionHashes

	txHashes, err = s.db.GetTxHashes(blockHash)
	if err != nil {
		return fmt.Errorf("commit: get txhashes for block %x..: %s", blockHash[:8], err.Error())
	}
	txs, missing := s.txPool.BatchGet(txHashes.TxHashes)
	if len(missing) > 0 {
		// Sanity check: if I TC'd the block then it means I have checked all txs. There should be no missing txs
		log.Panicf("commit: block %x.., should have %d txs but %d txs are still missing!",
			blockHash[:8], len(txHashes.TxHashes), len(missing))
	}
	err = s.commitTxs(txs)
	if err != nil {
		return err
	}
	s.txPool.BatchDelete(txHashes.TxHashes)
	return nil
}

func (s *Service) commitTxs(txs []*types.SignedTransaction) error {
	defer utils.LogExecTime(time.Now(), "commitTxs")

	log.Infof("commit %d txs", len(txs))
	dbtx, err := s.db.OpenTransaction()
	if err != nil {
		return err
	}

	balances := make(map[string]uint64)

	// populate current state first
	for _, tx := range txs {
		senderKey, receiverKey := accountKey(tx.Tx.From), accountKey(tx.Tx.To)
		senderBal, err := dbtx.Get(senderKey, nil)
		if err != nil && errors.Is(err, leveldb.ErrNotFound) {
			senderBal = make([]byte, 8)
		} else if err != nil {
			return err
		}
		senderBalance := binary.BigEndian.Uint64(senderBal)
		balances[string(senderKey)] = senderBalance

		receiverBal, err := dbtx.Get(receiverKey, nil)
		if err != nil && errors.Is(err, leveldb.ErrNotFound) {
			receiverBal = make([]byte, 8)
		} else if err != nil {
			return err
		}
		receiverBalance := binary.BigEndian.Uint64(receiverBal)
		balances[string(receiverKey)] = receiverBalance
	}

	// apply changes to state
	for _, tx := range txs {
		senderKey, receiverKey := accountKey(tx.Tx.From), accountKey(tx.Tx.To)
		if balances[string(senderKey)] < tx.Tx.Amount {
			// sanity check: this should already be checked before TC
			dbtx.Discard()
			return fmt.Errorf("sender balance %d < send amount %d", balances[string(senderKey)], tx.Tx.Amount)
		}
		balances[string(senderKey)] -= tx.Tx.Amount
		balances[string(receiverKey)] += tx.Tx.Amount
	}

	// save the new state to db
	batch := new(leveldb.Batch)
	for key, bal := range balances {
		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, bal)
		batch.Put([]byte(key), bs)
	}
	err = dbtx.Write(batch, nil)
	if err != nil {
		dbtx.Discard()
		return err
	}
	err = dbtx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) broadcastTxHashes(txs []*types.SignedTransaction) {
	hashes := make([][]byte, len(txs))
	for i, tx := range txs {
		hashes[i] = tx.Tx.Hash()
	}
	txHashes := &types.TransactionHashes{TxHashes: hashes}
	log.Infof("broadcastTxHashes: txs %d, root %x", len(txs), txHashes.Root())
	msg := &types.Message{
		Message: &types.Message_TxHashes{TxHashes: txHashes},
	}
	s.outMsgs.Put(s.txHashesKey(), msg, s.prepareGossipEndTime())
}

func (s *Service) getCurrentRound() uint32 {
	return uint32(time.Since(s.genesis.GenesisTime()) / s.cfg.RoundDuration())
}

func (s *Service) proposalStageEndTime() time.Time {
	return s.nextRoundTime().Add(-s.cfg.AgreementStageDuration())
}

func (s *Service) prepareGossipEndTime() time.Time {
	return s.proposalStageEndTime().Add(s.cfg.GossipDuration())
}

func (s *Service) nextRoundTime() time.Time {
	round := s.getCurrentRound()
	sinceGenesis := time.Duration(round+1) * s.cfg.RoundDuration()
	return s.genesis.GenesisTime().Add(sinceGenesis)
}

func (s *Service) gossipDDLFromNow() time.Time {
	return time.Now().Add(s.cfg.GossipDuration())
}

func computeTxRoot(signedTxs []*types.SignedTransaction) []byte {
	var txHashes []byte
	for _, tx := range signedTxs {
		txHashes = append(txHashes, tx.Tx.Hash()...)
	}
	root := blake2b.Sum256(txHashes)
	return root[:]
}

func (s *Service) proposalKey() string {
	return fmt.Sprintf("proposal-%d", s.round.Load())
}
func (s *Service) prepareKey() string         { return fmt.Sprintf("prepare-%d", s.round.Load()) }
func (s *Service) tcKey() string              { return fmt.Sprintf("tc-%d", s.round.Load()) }
func (s *Service) txHashesKey() string        { return fmt.Sprintf("txHashes-%d", s.round.Load()) }
func (s *Service) txKey(txHash []byte) string { return fmt.Sprintf("tx-%x", txHash) }
