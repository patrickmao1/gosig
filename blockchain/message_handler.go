package blockchain

import (
	"bytes"
	"fmt"
	bls12381 "github.com/kilic/bls12-381"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	"google.golang.org/protobuf/proto"
)

func (s *Service) handleMessage(msg *types.SignedMessage) {
	pass, err := s.checkSig(msg)
	if err != nil {
		log.Error("handleMessage", err)
		return
	}
	if !pass {
		log.Warnf("handleMessage: incorrect signature for message %v", msg)
		return
	}

	switch msg.Message.(type) {
	case *types.SignedMessage_Proposal:
		err = s.handleProposal(msg.ValidatorIndex, msg.GetProposal())
	case *types.SignedMessage_Prepare:
		err = s.handlePrepare(msg.GetPrepare())
	case *types.SignedMessage_Tc:
		err = s.handleTC(msg.GetTc())
	default:
		err = fmt.Errorf("handleMessage: unknown message type. msg %v", msg)
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
	if prop.BlockHeader.Height != s.head.Height {
		go s.runCatchUpRoutine()
	}
	// check proposal proof
	if !s.isValidProposal(vi, prop) {
		log.Warnf("received invalid proposal %v", prop)
		return nil
	}
	s.rmu.Lock()
	s.proposals = append(s.proposals, prop)
	s.rmu.Unlock()

	// relay the newly received proposal
	//s.msgBuf.Put() //TODO
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
	cert := incPrep.GetCert()
	bs, err := proto.Marshal(incPrep.Msg)
	if err != nil {
		return err
	}
	pass := crypto.VerifyAggSigBytes(s.cfg.Validators.PubKeys(), cert.SigCounts, bs, cert.AggSig)
	if !pass {
		return fmt.Errorf("got prepare with invalid cert: %+v", incPrep)
	}

	return s.msgBuf.Tx(func(buf *MsgBuffer) error {
		var err error
		mergedCert := incPrep.Cert
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
		mergedCert, err = mergeCerts(myPrep.Cert, incPrep.Cert)
		if err != nil {
			return err
		}
		mergedMsg := proto.CloneOf(myPrepMsg)
		mergedPrep := &types.PrepareCertificate{
			Msg:  myPrep.Msg,
			Cert: mergedCert,
		}
		mergedMsg.Message = &types.SignedMessage_Prepare{Prepare: mergedPrep}
		buf.Put(s.prepareKey(), mergedMsg)

		// TODO maybe move this outside of the tx
		if mergedPrep.Cert.NumSigned() >= s.cfg.Quorum() {
			// persist the prepare cert because if this block is later TC'd by me but never committed, I'll need to
			// re-propose this block with its P-cert as the proposal certificate.
			err := s.db.PutPCert(mergedPrep.Msg.BlockHash, mergedPrep.Cert)
			if err != nil {
				return err
			}

			// initiate a tc if I haven't yet
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
		log.Warnf("ignoring prepare: prepare.round %d != local round %d", incTc.Cert.Round, s.round)
		return nil
	}
	s.rmu.RUnlock()

	// Make sure the (partial) certificate is valid
	cert := incTc.GetCert()
	bs, err := proto.Marshal(incTc.Msg)
	if err != nil {
		return err
	}
	pass := crypto.VerifyAggSigBytes(s.cfg.Validators.PubKeys(), cert.SigCounts, bs, cert.AggSig)
	if !pass {
		return fmt.Errorf("got prepare with invalid cert: %+v", incTc)
	}
	return nil
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
	if a == nil || b == nil {
		return nil, fmt.Errorf("mergeCerts: nil cert")
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

func (s *Service) checkSig(msg *types.SignedMessage) (bool, error) {
	var message proto.Message
	switch msg.Message.(type) {
	case *types.SignedMessage_Proposal:
		message = msg.GetProposal()
	case *types.SignedMessage_Prepare:
		message = msg.GetPrepare()
	case *types.SignedMessage_Tc:
		message = msg.GetTc()
	default:
		return false, fmt.Errorf("unknown message type. msg %v", msg)
	}
	pubKey := s.cfg.Validators[msg.ValidatorIndex].GetPubKey()
	msgBytes, err := proto.Marshal(message)
	if err != nil {
		return false, err
	}
	return crypto.VerifySigBytes(pubKey, msgBytes, msg.Sig), nil
}
