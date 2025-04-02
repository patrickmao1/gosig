package blockchain

import "time"

type GenesisConfig struct {
	Round         int
	BlockNum      int
	GenesisTimeMs int64
	InitialSeed   []byte
}

func (c *GenesisConfig) GenesisTime() time.Time {
	return time.UnixMilli(c.GenesisTimeMs)
}
