package blockchain

import (
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

func (s *Service) handlePrepare(prep *types.PrepareCertificate) error {
	s.rmu.RLock()
	if prep.Cert.Round != s.round {
		s.rmu.RUnlock()
		log.Warnf("ignoring prepare: prepare.round %d != local round %d", prep.Cert.Round, s.round)
		return nil
	}
	s.rmu.RUnlock()

	// make sure the (partial) certificate is valid
	cert := prep.GetCert()
	bs, err := proto.Marshal(prep.Msg)
	if err != nil {
		return err
	}
	pass := crypto.VerifyAggSigBytes(s.cfg.Validators.PubKeys(), cert.SigCounts, bs, cert.AggSig)
	if !pass {
		return fmt.Errorf("got prepare with invalid cert: %+v", prep)
	}

	return s.msgBuf.Tx(func(buf *MsgBuffer) error {
		// merge prepare with my prepare state of this round
		var err error
		mergedCert := prep.Cert
		myPrep, ok := buf.Get(s.prepareKey())
		if !ok {
			// no prep found means either I'm not prepared for this round
			// or my round has ended and I have deleted it
			log.Debugf("no prepare found for %s", s.prepareKey())
			return nil
		}
		mergedCert, err = mergeCerts(myPrep.GetPrepare().Cert, prep.Cert)
		if err != nil {
			return err
		}
		clone := proto.CloneOf(myPrep)
		clone.Message = &types.SignedMessage_Prepare{
			Prepare: &types.PrepareCertificate{
				Msg:  prep.Msg,
				Cert: mergedCert,
			},
		}
		buf.Put(s.prepareKey(), clone)

		// TODO maybe move this outside of the tx
		if clone.GetPrepare().Cert.NumSigned() >= s.cfg.Quorum() {
			// initiate a tc if I haven't yet
			if !buf.Has(s.tcKey()) {
				tc := &types.TentativeCommit{
					BlockHash:   myPrep.GetPrepare().Msg.BlockHash,
					BlockHeight: myPrep.GetPrepare().Msg.BlockHeight,
				}
				bs, err := proto.Marshal(tc)
				if err != nil {
					return err
				}
				sig := crypto.SignBytes(s.cfg.privKey, bs)
				sigCounts := make([]uint32, len(s.cfg.Validators))
				tcWithCert := &types.TentativeCommitCertificate{
					Msg: tc,
					Cert: &types.Certificate{
						AggSig:    sig,
						SigCounts: sigCounts,
						Round:     myPrep.GetPrepare().Cert.Round,
					},
				}
				msg := &types.SignedMessage{
					Sig:            nil,
					ValidatorIndex: 0,
					Message:        &types.SignedMessage_Tc{Tc: tcWithCert},
					Deadline:       s.nextRoundTime().UnixMilli(),
				}
				buf.Put(s.tcKey(), msg)
			}
		}
		return nil
	})
}

func isSigSuperSet(a, b []uint32) bool {
	for i := range a {
		if a[i] == 0 && b[i] > 0 {
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

func (s *Service) handleTC(tc *types.TentativeCommitCertificate) error {
	// TODO
	// merge tc with my tc of this round
	return nil
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
	return crypto.GenRNGWithProof(prop.BlockHeader.ProposerProof) < s.cfg.ProposalThreshold
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
