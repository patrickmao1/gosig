package types

import (
	"github.com/patrickmao1/gosig/crypto"
	"golang.org/x/crypto/sha3"
	"google.golang.org/protobuf/proto"
	"time"
)

func (tx *Transaction) Hash() []byte {
	bs, err := proto.Marshal(tx)
	if err != nil {
		panic(err)
	}
	hash := sha3.Sum256(bs)
	return hash[:]
}

func (h *BlockHeader) Hash() []byte {
	bs, err := proto.Marshal(h)
	if err != nil {
		panic(err)
	}
	hash := sha3.Sum256(bs)
	return hash[:]
}

func (m *SignedMessage) DDL() time.Time {
	return time.UnixMilli(m.Deadline)
}

// NumSigned returns the number of validators who has provided a sig
func (c *Certificate) NumSigned() int {
	cnt := 0
	for _, count := range c.SigCounts {
		if count > 0 {
			cnt++
		}
	}
	return cnt
}

func (p *BlockProposal) Score() uint32 {
	return crypto.GenRNGWithProof(p.BlockHeader.ProposerProof)
}
