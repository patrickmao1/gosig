package utils

import "github.com/patrickmao1/gosig/crypto"

func GenTestKeyPairs(n int) (pubkeys, privkeys [][]byte) {
	for i := 0; i < n; i++ {
		sk, pk := crypto.GenKeyPairBytesFromSeed(int64(i))
		pubkeys = append(pubkeys, pk)
		privkeys = append(privkeys, sk)
	}
	return pubkeys, privkeys
}
