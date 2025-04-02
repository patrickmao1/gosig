package blockchain

import (
	bls12381 "github.com/kilic/bls12-381"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	"google.golang.org/protobuf/proto"
	"sync"
)

type TxPool struct {
	mu  sync.Mutex
	txs []*types.SignedTransaction
}

func NewTxPool() *TxPool {
	return &TxPool{}
}

func (p *TxPool) AppendTx(tx *types.SignedTransaction) error {
	pub, err := bls12381.NewG1().FromCompressed(tx.Tx.From)
	if err != nil {
		return err
	}
	txBytes, err := proto.Marshal(tx.Tx)
	if err != nil {
		return err
	}
	sig, err := bls12381.NewG2().FromCompressed(tx.Sig)
	if err != nil {
		return err
	}
	crypto.VerifySig(pub, txBytes, sig)

	p.mu.Lock()
	p.txs = append(p.txs, tx)
	p.mu.Unlock()

	return nil
}

func (p *TxPool) Dump() []*types.SignedTransaction {
	p.mu.Lock()
	ret := p.txs
	p.txs = nil
	p.mu.Unlock()
	return ret
}
