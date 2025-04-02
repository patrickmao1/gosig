package blockchain

import (
	"github.com/patrickmao1/gosig/types"
	"github.com/syndtr/goleveldb/leveldb"
)

type Blockchain struct {
}

func NewBlockchain(db *leveldb.DB, genesis *GenesisConfig) *Blockchain {
	return nil
}

// Head returns the latest committed head
func (b *Blockchain) Head() *types.BlockHeader {
	return nil
}
