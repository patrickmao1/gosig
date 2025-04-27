package blockchain

import (
	"fmt"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	"github.com/patrickmao1/gosig/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

type TxPool struct {
	txs             map[[32]byte]*types.SignedTransaction
	mu              sync.RWMutex
	recentlyDeleted map[[32]byte]struct{}
	mu2             sync.RWMutex
}

func NewTxPool() *TxPool {
	return &TxPool{
		txs:             make(map[[32]byte]*types.SignedTransaction),
		recentlyDeleted: make(map[[32]byte]struct{}),
	}
}

func (p *TxPool) CleanRecentlyDeleted() {
	p.mu2.Lock()
	p.recentlyDeleted = make(map[[32]byte]struct{})
	p.mu2.Unlock()
}

func (p *TxPool) AddTransaction(tx *types.SignedTransaction) error {
	var txHash [32]byte
	copy(txHash[:], tx.Tx.Hash())

	p.mu2.RLock()
	defer p.mu2.RUnlock()
	if _, ok := p.recentlyDeleted[txHash]; ok {
		return nil
	}

	txBytes, err := proto.Marshal(tx.Tx)
	if err != nil {
		return err
	}
	if !crypto.VerifySigBytes(tx.Tx.From, txBytes, tx.Sig) {
		return fmt.Errorf("invalid sig")
	}

	p.mu.Lock()
	p.txs[txHash] = tx
	p.mu.Unlock()

	return nil
}

func (p *TxPool) AddTransactions(txs []*types.SignedTransaction) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	pass := checkTxs(txs)
	if !pass {
		return fmt.Errorf("tx check failed")
	}
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

func (p *TxPool) BatchDelete(txHashes [][]byte) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.mu2.Lock()
	defer p.mu2.Unlock()

	for _, hash := range txHashes {
		var key [32]byte
		copy(key[:], hash)
		delete(p.txs, key)
		p.recentlyDeleted[key] = struct{}{}
	}
	log.Infof("deleted %d tx hashes", len(txHashes))
}

func checkTxs(txs []*types.SignedTransaction) bool {
	defer utils.LogExecTime(time.Now(), "checkTxs")
	for _, tx := range txs {
		bs, err := proto.Marshal(tx.Tx)
		if err != nil {
			log.Error(err)
			return false
		}
		pass := crypto.VerifySigBytes(tx.Tx.From, bs, tx.Sig)
		if !pass {
			log.Errorf("tx %s invalid signature", tx.ToString())
			return false
		}
	}
	return true
}
