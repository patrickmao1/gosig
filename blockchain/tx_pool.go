package blockchain

import (
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"sync"
)

type TxPool struct {
	txs      map[[32]byte]*types.SignedTransaction
	mu       sync.RWMutex
	toDelete [][32]byte
}

func NewTxPool() *TxPool {
	return &TxPool{
		txs: make(map[[32]byte]*types.SignedTransaction),
	}
}

func (p *TxPool) AddTransaction(tx *types.SignedTransaction) error {
	var txHash [32]byte
	copy(txHash[:], tx.Tx.Hash())

	p.mu.Lock()
	p.txs[txHash] = tx
	p.mu.Unlock()

	return nil
}

func (p *TxPool) AddTransactions(txs []*types.SignedTransaction) error {
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
	p.mu.RLock()
	defer p.mu.RUnlock()

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

func (p *TxPool) BatchGet(txHashes [][]byte) (found []*types.SignedTransaction, missing [][]byte) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, hash := range txHashes {
		var key [32]byte
		copy(key[:], hash)
		tx, ok := p.txs[key]
		if ok {
			found = append(found, tx)
		} else {
			missing = append(missing, hash)
		}
	}
	return
}

func (p *TxPool) DoDelete() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, key := range p.toDelete {
		delete(p.txs, key)
	}
	p.toDelete = nil
}

func (p *TxPool) ScheduleDelete(txHashes [][]byte) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, hash := range txHashes {
		var key [32]byte
		copy(key[:], hash)
		p.toDelete = append(p.toDelete, key)
	}
	log.Infof("deleted %d txs", len(txHashes))
}

func checkTxs(txs []*types.SignedTransaction) bool {
	for _, tx := range txs {
		bs, err := proto.Marshal(tx.Tx)
		if err != nil {
			log.Error(err)
			return false
		}
		pass := crypto.VerifyECDSABytes(tx.Tx.From, bs, tx.Sig)
		if !pass {
			log.Errorf("tx %s invalid signature", tx.ToString())
			return false
		}
	}
	return true
}
