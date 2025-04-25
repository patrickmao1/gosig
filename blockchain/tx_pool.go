package blockchain

import (
	"fmt"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"sync"
)

type TxPool struct {
	mu  sync.Mutex
	txs map[[32]byte]*types.SignedTransaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		txs: make(map[[32]byte]*types.SignedTransaction),
	}
}

func (p *TxPool) AddTransaction(tx *types.SignedTransaction) error {
	txBytes, err := proto.Marshal(tx.Tx)
	if err != nil {
		return err
	}
	if !crypto.VerifySigBytes(tx.Tx.From, txBytes, tx.Sig) {
		return fmt.Errorf("invalid sig")
	}

	var txHash [32]byte
	copy(txHash[:], tx.Tx.Hash())

	p.mu.Lock()
	p.txs[txHash] = tx
	p.mu.Unlock()

	return nil
}

func (p *TxPool) AddTransactions(txs []*types.SignedTransaction) error {
	for _, tx := range txs {
		txBytes, err := proto.Marshal(tx.Tx)
		if err != nil {
			return err
		}
		if !crypto.VerifySigBytes(tx.Tx.From, txBytes, tx.Sig) {
			return fmt.Errorf("invalid sig")
		}
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	for _, tx := range txs {
		var txHash [32]byte
		copy(txHash[:], tx.Tx.Hash())

		p.txs[txHash] = tx
	}

	return nil
}

func (p *TxPool) List(n int) []*types.SignedTransaction {
	p.mu.Lock()
	defer p.mu.Unlock()

	var ret []*types.SignedTransaction
	count := 0
	for _, tx := range p.txs {
		count++
		ret = append(ret, tx)
		if count >= n {
			break
		}
	}
	return ret
}

func (p *TxPool) BatchGet(txHashes [][]byte) ([]*types.SignedTransaction, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var ret []*types.SignedTransaction
	for _, hash := range txHashes {
		var key [32]byte
		copy(key[:], hash)
		tx, ok := p.txs[key]
		if !ok {
			return nil, fmt.Errorf("tx %x not found in tx pool", key)
		}
		ret = append(ret, tx)
	}
	return ret, nil
}

func (p *TxPool) BatchDelete(txHashes [][]byte) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, hash := range txHashes {
		var key [32]byte
		copy(key[:], hash)
		delete(p.txs, key)
	}
	log.Infof("deleted %d tx hashes", len(txHashes))
}
