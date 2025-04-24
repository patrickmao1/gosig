package blockchain

import (
	"github.com/patrickmao1/gosig/types"
	"golang.org/x/crypto/blake2b"
	"time"
)

type GenesisConfig struct {
	GenesisTimeMs int64
	InitialSeed   []byte
}

func (c *GenesisConfig) GenesisTime() time.Time {
	return time.UnixMilli(c.GenesisTimeMs)
}

func DefaultGenesisConfig() *GenesisConfig {
	t := time.Now().Unix() / 3600 * 3600 * 1000 // the beginning of the current hour in ms
	initSeed := blake2b.Sum256([]byte("Gosig"))
	return &GenesisConfig{
		GenesisTimeMs: t,
		InitialSeed:   initSeed[:],
	}
}

func NewGenesisBlock(initSeed []byte) *types.BlockHeader {
	return &types.BlockHeader{
		Height:        0,
		ParentHash:    initSeed[:],
		ProposerProof: initSeed[:],
		TxRoot:        initSeed[:],
	}
}
