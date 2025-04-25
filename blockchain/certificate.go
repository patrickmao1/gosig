package blockchain

import (
	"fmt"
	bls12381 "github.com/kilic/bls12-381"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func (s *Service) checkCert(msg proto.Message, cert *types.Certificate) bool {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Errorf("checkCert: marshal msg err: %v", err)
		return false
	}
	return crypto.VerifyAggSigBytes(
		s.cfg.Validators.PubKeys(),
		cert.SigCounts,
		data,
		cert.AggSig,
	)
}

func (s *Service) signNewCert(msg proto.Message) (*types.Certificate, error) {
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

func (s *Service) addSigToCert(msg proto.Message, cert *types.Certificate) error {
	if cert.SigCounts[s.cfg.MyValidatorIndex()] > 0 {
		return nil
	}
	bs, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	sig := crypto.SignBytes(s.cfg.MyPrivKey(), bs)
	cert.SigCounts[s.cfg.MyValidatorIndex()] += 1
	cert.AggSig = crypto.AggSigsBytes([][]byte{cert.AggSig, sig})
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
