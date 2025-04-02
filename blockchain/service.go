package blockchain

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/crypto/sha3"
	"google.golang.org/protobuf/proto"
	"slices"
	"sync"
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
	rseed     []byte // round concatenated with seed
	proposals []*types.BlockProposal
}

func NewService(cfg *NodeConfig, genesis *GenesisConfig, db *DB, txPool *TxPool, nw *Network) *Service {
	roundInterval := time.Duration(cfg.ProposalStageDurationMs + cfg.ProposalStageDurationMs)
	return &Service{
		cfg:           cfg,
		genesis:       genesis,
		roundInterval: roundInterval,
		db:            db,
		txPool:        txPool,
		nw:            nw,
	}
}

func (s *Service) Start() error {
	nextRoundTime := s.nextRoundTime()
	t := time.NewTicker(time.Until(nextRoundTime))
	for {
		<-t.C
		s.rmu.Lock()

		// moving to a new round
		s.round++
		s.proposals = nil

		// compute the seed for the round
		head, err := s.db.GetHeadBlock()
		if err != nil {
			log.Error("failed to get head block", err)
			return err
		}
		seed := sha3.Sum256(head.ProposerProof)
		s.seed = seed[:]
		rb := make([]byte, 4)
		binary.BigEndian.PutUint32(rb, s.round)
		s.rseed = slices.Concat(rb, s.seed)

		s.rmu.Unlock()

		// determine if I should propose
		proposalScore, proposerProof := crypto.GenRNG(s.cfg.MyPrivKey(), s.rseed) // L^r = SIG_{l^{h}}(r^{h},Q^{h-1})
		if proposalScore < s.cfg.ProposalThreshold {
			// build a block
			headHash := head.Hash()
			header := &types.BlockHeader{
				Height:        head.Height + 1,
				ParentHash:    headHash,
				ProposerProof: proposerProof,
			}
			tcBlock, _, pCert, err := s.db.GetTcState()
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
			} else {
				// Sanity check. If these aren't the same then I wouldn't have TC'd that block
				if !bytes.Equal(tcBlock.ParentHash, head.ParentHash) {
					log.Panicf("inconsistent parent hashes! tcBlock %v, headBlock %v", tcBlock, head)
				}
				// I have previously TC'd a block, but I never committed it. Re-propose the block.
				txs, header.TxRoot, err = s.db.GetBlockTxs(tcBlock.Hash())
				if err != nil {
					return err
				}
				header.ProposalCert = pCert
			}

			// build & sign proposal
			proposal := &types.BlockProposal{
				Round:       s.round,
				BlockHeader: header,
			}
			bs, err := proto.Marshal(proposal)
			if err != nil {
				return err
			}
			signed := &types.SignedMessage{
				Sig:            crypto.Sign(s.cfg.MyPrivKey(), bs),
				ValidatorIndex: s.cfg.MyValidatorIndex(),
				MessageTypes:   &types.SignedMessage_Proposal{Proposal: proposal},
			}

			// broadcast
			s.nw.Broadcast(signed, nextRoundTime)
		}
	}
}

func (s *Service) getCurrentRound() int {
	return int(time.Since(s.genesis.GenesisTime()) / s.roundInterval)
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
