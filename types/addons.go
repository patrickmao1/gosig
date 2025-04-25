package types

import (
	"fmt"
	"github.com/patrickmao1/gosig/crypto"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/blake2b"
	"google.golang.org/protobuf/proto"
	"slices"
)

func (tx *SignedTransaction) Hash() []byte {
	return tx.Tx.Hash()
}

func (tx *SignedTransaction) ToString() string {
	if tx == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SignedTransaction(%s, Sig %x..)", tx.Tx.ToString(), tx.Sig[:8])
}

func (tx *Transaction) Hash() []byte {
	bs, err := proto.Marshal(tx)
	if err != nil {
		panic(err)
	}
	hash := blake2b.Sum256(bs)
	return hash[:]
}

func (tx *Transaction) ToString() string {
	if tx == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Transaction(hash %x.., From %x.., To %x.., Amount %d)", tx.Hash()[:8], tx.From[:8], tx.To[:8], tx.Amount)
}

func (txs *TransactionHashes) Root() []byte {
	if txs == nil {
		return nil
	}
	root := blake2b.Sum256(slices.Concat(txs.TxHashes...))
	return root[:]
}

func (h *BlockHeader) Hash() []byte {
	bs, err := proto.Marshal(h)
	if err != nil {
		panic(err)
	}
	hash := blake2b.Sum256(bs)
	return hash[:]
}

func (h *BlockHeader) ToString() string {
	if h == nil {
		return "<nil>"
	}
	return fmt.Sprintf("BlockHeader(Height %d, ParentHash %x.., ProposerProof %x.., TxRoot %x..)",
		h.Height, h.ParentHash[:8], h.ProposerProof[:8], h.TxRoot[:8])
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

func (c *Certificate) ToString() string {
	if c == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Certificate(SigCounts %d, AggSig %x.., Round %d)", c.SigCounts, c.AggSig[:8], c.Round)
}

func (m *Message) ToString() string {
	if m == nil {
		return "<nil>"
	}
	switch m.Message.(type) {
	case *Message_Proposal:
		return m.GetProposal().ToString()
	case *Message_Prepare:
		return m.GetPrepare().ToString()
	case *Message_Tc:
		return m.GetTc().ToString()
	default:
		log.Panic("Unknown Message type %T", m.Message)
	}
	return "<nil>"
}

func (p *BlockProposal) Score() uint32 {
	return crypto.RngFromProof(p.BlockHeader.ProposerProof)
}

func (p *BlockProposal) ToString() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("BlockProposal(ProposerIndex %d, Round %d, Cert %s, BlockHeader %s)",
		p.ProposerIndex, p.Round, p.Cert.ToString(), p.BlockHeader.ToString())
}

func (p *PrepareCertificate) ToString() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("PrepareCertificate(Msg %s, Cert %s)", p.Msg.ToString(), p.Cert.ToString())
}

func (p *Prepare) ToString() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Prepare(BlockHash %x..)", p.BlockHash[:8])
}

func (t *TentativeCommitCertificate) ToString() string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf("TentativeCommitCertificate(Tc %s, Cert %s)", t.Msg.ToString(), t.Cert.ToString())
}

func (t *TentativeCommit) ToString() string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf("TentativeCommit(BlockHash %x..)", t.BlockHash[:8])
}
