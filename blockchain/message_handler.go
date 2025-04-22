package blockchain

import (
	"bytes"
	"fmt"
	bls12381 "github.com/kilic/bls12-381"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	"google.golang.org/protobuf/proto"
)

func (s *Service) startProcessingMsgs() error {
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
	pass, err := s.checkSig(signedMsg)
	if err != nil {
		log.Error("handleMessage", err)
		return
	}
	if !pass {
		log.Warnf("handleMessage: incorrect signature for message %v", signedMsg)
		return
	}

	msg := signedMsg.Msg

	switch signedMsg.Msg.Message.(type) {
	case *types.Message_Proposal:
		err = s.handleProposal(signedMsg.ValidatorIndex, msg.GetProposal())
	case *types.Message_Prepare:
		err = s.handlePrepare(msg.GetPrepare())
	case *types.Message_Tc:
		err = s.handleTC(msg.GetTc())
	default:
		err = fmt.Errorf("handleMessage: unknown message type. msg %v", signedMsg)
	}
	log.Errorln("handleMessage", err)
}

func (s *Service) handleProposal(vi uint32, prop *types.BlockProposal) error {
	s.rmu.RLock()
	defer s.rmu.RUnlock()
	if prop.Round != s.round {
		log.Warnf("received proposal for round %d but current round is %d", prop.Round, s.round)
		return nil
	}
	if prop.BlockHeader.Height > s.head.Height+1 {
		log.Warnf("received proposal block (%d) > head block + 1 (%d)", prop.BlockHeader.Height, s.head.Height+1)
		go s.runCatchUpRoutine()
		// We don't put back the proposal to queue because we expect to receive it again from gossip
	}
	// check proposal proof
	if !s.isValidProposal(vi, prop) {
		log.Warnf("received invalid proposal %v", prop)
		return nil
	}
	s.rmu.Lock()
	s.proposals[vi] = prop
	s.rmu.Unlock()

	// relay the newly received proposal
	msg := &types.Message{Message: &types.Message_Proposal{Proposal: prop}}
	s.outMsgs.Put(s.proposalKey(vi), msg)
	return nil
}

func (s *Service) handlePrepare(incPrep *types.PrepareCertificate) error {
	s.rmu.RLock()
	if incPrep.Cert.Round != s.round {
		s.rmu.RUnlock()
		log.Warnf("ignoring prepare: prepare.round %d != local round %d", incPrep.Cert.Round, s.round)
		return nil
	}
	s.rmu.RUnlock()

	// Make sure the (partial) certificate is valid
	err := s.verifyCertificate(incPrep.Msg, incPrep.GetCert())
	if err != nil {
		return fmt.Errorf("failed to handle prepare msg: %s", err.Error())
	}

	return s.outMsgs.Tx(func(buf *OutboundMsgBuffer) error {
		myPrepMsg, ok := buf.Get(s.prepareKey())
		if !ok {
			// No prep found means either I'm not prepared for this round
			// or my round has ended and I have deleted it.
			// According to the paper, I should only handle prepare msgs if I'm prepared
			log.Debugf("no prepare found for %s", s.prepareKey())
			return nil
		}
		myPrep := myPrepMsg.GetPrepare()

		if !bytes.Equal(myPrep.Msg.BlockHash, incPrep.Msg.BlockHash) {
			return fmt.Errorf("got prepare different from mine: theirs %+v, mine %+v", incPrep, myPrep)
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
		buf.Put(s.prepareKey(), mergedMsg)

		if mergedPrep.Cert.NumSigned() >= s.cfg.Quorum() {
			// persist the prepare cert because if this block is later TC'd by me but never committed, I'll need to
			// re-propose this block with its P-cert as the proposal certificate.
			err := s.db.PutPCert(mergedPrep.Msg.BlockHash, mergedPrep.Cert)
			if err != nil {
				return err
			}

			// initiate a tc if I haven't yet
			tc, ok := buf.Get(s.tcKey())
			// I already have a TC msg. This could either mean that I have previously processed
			// P msgs that reached quorum so that I initiated a TC, or that I have received a TC
			// msg before I collect quorum sigs form P. In the former case, I have signed the TC
			// msg, no need to sign it again. In the latter case, I need to add my TC signature.
			if ok {
				if tc.GetTc().Cert.SigCounts[s.cfg.MyValidatorIndex()] == 0 {
					// TODO
				}
			}
			if !buf.Has(s.tcKey()) {
				err := s.tentativeCommit(mergedPrep.Msg.BlockHash)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *Service) handleTC(incTc *types.TentativeCommitCertificate) error {
	// merge tc with my tc of this round
	s.rmu.RLock()
	if incTc.Cert.Round != s.round {
		s.rmu.RUnlock()
		log.Warnf("ignoring TC msg: TC.round %d != local round %d", incTc.Cert.Round, s.round)
		return nil
	}
	s.rmu.RUnlock()

	// Make sure the (partial) certificate is valid
	err := s.verifyCertificate(incTc.Msg, incTc.GetCert())
	if err != nil {
		return fmt.Errorf("failed to handle TC msg: %s", err.Error())
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
			return fmt.Errorf("incTc block != my prepare block: %x != %x",
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
			myTcMsg = &types.Message{
				Message:  nil,
				Deadline: 0,
			}
		}

		if !bytes.Equal(myTc.Msg.BlockHash, incTc.Msg.BlockHash) {
			return fmt.Errorf("got TC different from mine: theirs %+v, mine %+v", incTc, myTc)
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
		buf.Put(s.tcKey(), mergedMsg)

		if mergedTc.Cert.NumSigned() >= s.cfg.Quorum() {
			return s.commit(mergedTc.Msg.BlockHash)
		}
		return nil
	})
}

func isSigSuperSet(a, b []uint32) bool {
	for i := range a {
		if b[i] > 0 && a[i] < b[i] {
			return false
		}
	}
	return true
}

func mergeCerts(a, b *types.Certificate) (*types.Certificate, error) {
	if a == nil && b == nil {
		return nil, fmt.Errorf("mergeCerts: nil certs")
	} else if a != nil && b == nil {
		return a, nil
	} else if a == nil {
		return b, nil
	}
	if a.Round != b.Round {
		return nil, fmt.Errorf("mergeCerts: round not equal: a %v, b %v", a, b)
	}

	// optimization: if one cert is a superset of another, we could just use the superset without merging
	if isSigSuperSet(a.SigCounts, b.SigCounts) {
		return a, nil
	}
	if isSigSuperSet(b.SigCounts, a.SigCounts) {
		return b, nil
	}

	g2 := bls12381.NewG2()
	aSig, err := g2.FromCompressed(a.AggSig)
	if err != nil {
		return nil, err
	}
	bSig, err := g2.FromCompressed(b.AggSig)
	if err != nil {
		return nil, err
	}
	sum := g2.Add(g2.New(), aSig, bSig)
	aggSig := g2.ToCompressed(sum)

	// sanity check. all nodes have hardcoded validator list. these should equal
	if len(a.SigCounts) != len(b.SigCounts) {
		log.Panic("mergeCerts: sigCounts len not the same")
	}
	retCounts := make([]uint32, len(a.SigCounts))
	for i := range a.SigCounts {
		retCounts[i] = a.SigCounts[i] + b.SigCounts[i]
	}
	return &types.Certificate{
		AggSig:    aggSig,
		SigCounts: retCounts,
		Round:     a.Round,
	}, nil
}

func (s *Service) isValidProposal(vi uint32, prop *types.BlockProposal) bool {
	s.rmu.RLock()
	defer s.rmu.RUnlock()
	// verify RNG is generated correctly with the proposer's private key
	rngValid := crypto.VerifySigBytes(s.cfg.Validators[vi].GetPubKey(), s.rseed, prop.BlockHeader.ProposerProof)
	if !rngValid {
		return false
	}
	// proposer's score must be lower than the threshold
	return crypto.VRF2(prop.BlockHeader.ProposerProof) < s.cfg.ProposalThreshold
}

func (s *Service) checkSig(signedMsg *types.Envelope) (bool, error) {
	pubKey := s.cfg.Validators[signedMsg.ValidatorIndex].GetPubKey()
	msgBytes, err := proto.Marshal(signedMsg.Msg)
	if err != nil {
		return false, err
	}
	return crypto.VerifySigBytes(pubKey, msgBytes, signedMsg.Sig), nil
}
