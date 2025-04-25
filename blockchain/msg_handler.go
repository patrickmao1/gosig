package blockchain

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/protobuf/proto"
)

func (s *Service) startProcessingMsgs() {
	log.Infoln("Starting processing messages")
	for {
		msgs := s.inMsgs.DequeueAll() // blocking if no msg
		for _, msg := range msgs {
			// TODO maybe parallelize
			s.handleMessage(msg)
		}
	}
}

func (s *Service) handleMessage(signedMsg *types.Envelope) {
	msg := signedMsg.Msg
	var err error
	switch signedMsg.Msg.Message.(type) {
	case *types.Message_Proposal:
		err = s.handleProposal(msg.GetProposal())
	case *types.Message_Prepare:
		err = s.handlePrepare(msg.GetPrepare())
	case *types.Message_Tc:
		err = s.handleTC(msg.GetTc())
	case *types.Message_TxHashes:
		err = s.handleTxHashes(msg.GetTxHashes())
	case *types.Message_Tx:
		err = s.handleTransaction(msg.GetTx())
	default:
		err = fmt.Errorf("handleMessage: unknown message type. msg %v", signedMsg)
	}
	if err != nil {
		log.Errorf("handleMessage %s err: %s", msg.ToString(), err.Error())
	}
}

func (s *Service) handleProposal(prop *types.BlockProposal) error {
	log.Debugf("handleProposal: block hash %x.., block %s", prop.BlockHeader.Hash()[:8], prop.ToString())
	s.rmu.Lock()
	defer s.rmu.Unlock()
	if prop.Round != s.round.Load() {
		log.Warnf("received proposal for round %d but current round is %d", prop.Round, s.round)
		return nil
	}
	if prop.BlockHeader.Height > s.head.Height+1 {
		log.Warnf("received proposal block (%d) > head block + 1 (%d)", prop.BlockHeader.Height, s.head.Height+1)
		return nil
	}
	// check proposal proof
	if !s.isValidProposal(prop) {
		log.Warnf("received invalid proposal")
		return nil
	}

	if s.minProposal == nil || prop.Score() < s.minProposal.Score() {
		s.minProposal = prop
	}

	// relay the known best proposal
	msg := &types.Message{
		Message: &types.Message_Proposal{Proposal: s.minProposal},
	}
	s.outMsgs.Put(s.proposalKey(), msg, s.proposalStageEndTime())
	return nil
}

func (s *Service) handlePrepare(incPrep *types.PrepareCertificate) error {
	log.Debugf("handlePrepare %s", incPrep.ToString())
	if incPrep.Cert.Round != s.round.Load() {
		log.Warnf("ignoring prepare: prepare.round %d != local round %d", incPrep.Cert.Round, s.round)
		return nil
	}

	// Make sure the (partial) certificate is valid
	pass := s.checkCert(incPrep.Msg, incPrep.GetCert())
	if !pass {
		return fmt.Errorf("handlePrepare: failed to handle Prepare msg: invalid cert")
	}

	return s.outMsgs.Tx(func(buf *OutboundMsgBuffer) error {
		myPrepMsg, ok := buf.Get(s.prepareKey())
		if !ok {
			// No prep found means either I'm not prepared for this round
			// or my round has ended and I have deleted it.
			// According to the paper, I should only handle prepare msgs if I'm prepared
			log.Infof("no prepare found for %s", s.prepareKey())
			return nil
		}
		myPrep := myPrepMsg.GetPrepare()

		if !bytes.Equal(myPrep.Msg.BlockHash, incPrep.Msg.BlockHash) {
			return fmt.Errorf("got prepare different from mine: theirs %s, mine %s", incPrep.ToString(), myPrep.ToString())
		}

		// Merge prepare with my prepare state of this round
		mergedCert, err := mergeCerts(myPrep.Cert, incPrep.Cert)
		if err != nil {
			return err
		}
		mergedMsg := proto.CloneOf(myPrepMsg)
		mergedPrep := &types.PrepareCertificate{
			Msg:  myPrep.Msg,
			Cert: mergedCert,
		}
		mergedMsg.Message = &types.Message_Prepare{Prepare: mergedPrep}
		buf.Put(s.prepareKey(), mergedMsg, s.prepareGossipEndTime(), s.nextRoundTime())

		log.Debugf("saved merged prepare %s", mergedMsg.ToString())

		if mergedPrep.Cert.NumSigned() >= s.cfg.Quorum() {
			log.Debugf("handlePrepare: reached quorum %d/%d", mergedPrep.Cert.NumSigned(), s.cfg.Quorum())
			// persist the prepare cert because if this block is later TC'd by me but never committed, I'll need to
			// re-propose this block with its P-cert as the proposal certificate.
			err := s.db.PutPCert(mergedPrep.Msg.BlockHash, mergedPrep.Cert)
			if err != nil {
				return err
			}

			// initiate a tc if I haven't yet
			tcMsg, ok := buf.Get(s.tcKey())
			// I already have a TC msg. This could either mean that I have previously processed
			// P msgs that reached quorum so that I initiated a TC, or that I have received a TC
			// msg before I collect quorum sigs form P. In the former case, I have signed the TC
			// msg, no need to sign it again. In the latter case, I need to add my TC signature.
			if ok {
				if tcMsg.GetTc().Cert.SigCounts[s.cfg.MyValidatorIndex()] == 0 {
					tcCert := tcMsg.GetTc()
					err := s.addSigToCert(tcCert.Msg, tcCert.Cert)
					if err != nil {
						return err
					}
					buf.Put(s.tcKey(), tcMsg, s.nextRoundTime())
				}
			}
			_, err = s.db.GetTcCert(mergedPrep.Msg.BlockHash)
			if err != nil && errors.Is(err, leveldb.ErrNotFound) {
				err := s.tentativeCommit(buf, mergedPrep.Msg.BlockHash)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *Service) handleTC(incTc *types.TentativeCommitCertificate) error {
	log.Debugf("handleTC %s", incTc.ToString())
	// merge tc with my tc of this round
	if incTc.Cert.Round != s.round.Load() {
		log.Warnf("handleTC: ignoring TC msg: TC.round %d != local round %d", incTc.Cert.Round, s.round)
		return nil
	}

	// Make sure the (partial) certificate is valid
	pass := s.checkCert(incTc.Msg, incTc.GetCert())
	if !pass {
		return fmt.Errorf("handleTC: failed to handle TC msg: invalid cert")
	}

	return s.outMsgs.Tx(func(buf *OutboundMsgBuffer) error {
		// The paper specifies that TC msgs can be handled as soon as I have prepared something.
		// This is for pipelining the Prepare and TC stages. I should process a TC msg if the TC
		// msg's cert is correct and I have a matching P msg regardless whether I have received
		// enough P sigs to TC myself.
		myPrepMsg, ok := buf.Get(s.prepareKey())
		if !ok {
			return fmt.Errorf("failed to handle TC msg: no prepare found for %s", s.prepareKey())
		}
		myPrep := myPrepMsg.GetPrepare()
		if !bytes.Equal(myPrep.Msg.BlockHash, incTc.Msg.BlockHash) {
			return fmt.Errorf("handleTC: incTc block != my prepare block: %x != %x",
				incTc.Msg.BlockHash, myPrep.Msg.BlockHash)
		}
		myTc := &types.TentativeCommitCertificate{
			Msg:  &types.TentativeCommit{BlockHash: myPrep.Msg.BlockHash},
			Cert: nil,
		}
		myTcMsg, ok := buf.Get(s.tcKey())
		if ok {
			myTc = myTcMsg.GetTc()
		} else {
			myTcMsg = &types.Message{}
		}

		if !bytes.Equal(myTc.Msg.BlockHash, incTc.Msg.BlockHash) {
			return fmt.Errorf("handleTC: got TC different from mine: theirs %+v, mine %+v", incTc, myTc)
		}

		// Merge incoming tc with my tc of this round
		mergedCert, err := mergeCerts(myTc.Cert, incTc.Cert)
		if err != nil {
			return err
		}
		mergedMsg := proto.CloneOf(myTcMsg)
		mergedTc := &types.TentativeCommitCertificate{
			Msg:  incTc.Msg,
			Cert: mergedCert,
		}
		mergedMsg.Message = &types.Message_Tc{Tc: mergedTc}
		buf.Put(s.tcKey(), mergedMsg, s.nextRoundTime())

		if mergedTc.Cert.NumSigned() >= s.cfg.Quorum() {
			return s.commit(mergedTc.Msg.BlockHash)
		}
		return nil
	})
}

func (s *Service) handleTxHashes(txHashes *types.TransactionHashes) error {
	s.rmu.RLock()
	defer s.rmu.RUnlock()

	expected := s.minProposal.BlockHeader.TxRoot
	received := txHashes.Root()
	if !bytes.Equal(received, expected) {
		return fmt.Errorf("handleTxHashes: tx root mismatch: expected %x, got %x", expected, received)
	}

	blockHash := s.minProposal.BlockHeader.Hash()
	msg := &types.Message{
		Message: &types.Message_TxHashes{TxHashes: txHashes},
	}
	s.outMsgs.Put(s.txHashesKey(), msg, s.prepareGossipEndTime())
	return s.db.PutTxHashes(blockHash, txHashes)
}

func (s *Service) handleTransaction(tx *types.SignedTransaction) error {
	return s.txPool.AddTransaction(tx)
}

func (s *Service) isValidProposal(prop *types.BlockProposal) bool {
	// no need to lock here cuz the caller has the lock

	// Verify RNG is generated correctly with the proposer's private key
	log.Debugf("verifying RNG: pubkey %x.., rseed %x, proof %x..",
		s.cfg.Validators[prop.ProposerIndex].GetPubKey()[:8],
		s.rseed,
		prop.BlockHeader.ProposerProof[:8],
	)
	rngValid := crypto.VerifySigBytes(
		s.cfg.Validators[prop.ProposerIndex].GetPubKey(),
		s.rseed,
		prop.BlockHeader.ProposerProof,
	)
	if !rngValid {
		log.Errorf("proposal %s: invalid RNG", prop.ToString())
		return false
	}

	// Proposer's score must be lower than the threshold
	if prop.Score() >= s.cfg.ProposalThreshold {
		log.Errorf("proposal %s: score too high %d", prop.ToString(), prop.Score())
		return false
	}

	// Check proposal cert
	block := prop.BlockHeader
	// if the previous block is the genesis block then it doesn't have a proposal cert
	if bytes.Equal(block.ParentHash, s.genesisBlockHash) {
		return true
	}
	// Assuming it's a newly proposed block (not a previously tc'd block)
	// The cert is the TC cert of the previous block
	tc := &types.TentativeCommit{BlockHash: block.ParentHash}
	pass := s.checkCert(tc, prop.Cert)
	if !pass {
		log.Errorf("proposal %s: invalid cert sig", prop.ToString())
		return false
	}
	return true
}
