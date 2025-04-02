package blockchain

import (
	"fmt"
	bls12381 "github.com/kilic/bls12-381"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	"google.golang.org/protobuf/proto"
)

func (s *Service) HandleMessage(msg *types.SignedMessage) {
	pass, err := s.checkSig(msg)
	if err != nil {
		log.Error("HandleMessage", err)
		return
	}
	if !pass {
		log.Warnf("HandleMessage: incorrect signature for message %v", msg)
		return
	}

	switch msg.MessageTypes.(type) {
	case *types.SignedMessage_Proposal:
		err = s.handleProposal(msg.ValidatorIndex, msg.GetProposal())
	case *types.SignedMessage_Prepare:
		err = s.handlePrepare(msg.GetPrepare())
	case *types.SignedMessage_Tc:
		err = s.handleTC(msg.GetTc())
	default:
		err = fmt.Errorf("HandleMessage: unknown message type. msg %v", msg)
	}
	log.Errorln("HandleMessage", err)
}

func (s *Service) handleProposal(vi uint32, prop *types.BlockProposal) error {
	s.rmu.RLock()
	defer s.rmu.RUnlock()
	if prop.Round != s.round {
		log.Warnf("received proposal for round %d but current round is %d", prop.Round, s.round)
		return nil
	}
	// check proposal proof
	if !s.isValidProposal(vi, prop) {
		log.Warnf("received invalid proposal %v", prop)
		return nil
	}
	s.rmu.Lock()
	s.proposals = append(s.proposals, prop)
	s.rmu.Unlock()
	return nil
}

func (s *Service) handlePrepare(prep *types.Prepare) error {
	// TODO
	return nil
}

func (s *Service) handleTC(tc *types.TentativeCommit) error {
	// TODO
	return nil
}

func (s *Service) isValidProposal(vi uint32, prop *types.BlockProposal) bool {
	return crypto.VerifyRNG(s.cfg.Validators[vi].GetPubKey(), prop.BlockHeader.ProposerProof, s.rseed)
}

func (s *Service) checkSig(msg *types.SignedMessage) (bool, error) {
	var message proto.Message
	switch msg.MessageTypes.(type) {
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
	sig, err := bls12381.NewG2().FromCompressed(msg.Sig)
	if err != nil {
		return false, err
	}
	return crypto.VerifySig(pubKey, msgBytes, sig), nil
}
