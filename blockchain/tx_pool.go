package blockchain

import (
	"fmt"
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
	txBytes, err := proto.Marshal(tx.Tx)
	if err != nil {
		return err
	}
	if crypto.VerifySigBytes(tx.Tx.From, txBytes, tx.Sig) {
		return fmt.Errorf("invalid sig")
	}

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
