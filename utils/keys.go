package utils

import (
	"crypto/ecdsa"
	"github.com/patrickmao1/gosig/crypto"
)

func GenTestKeyPairs(n int) (pubkeys, privkeys [][]byte) {
	for i := 0; i < n; i++ {
		sk, pk := crypto.GenKeyPairBytesFromSeed(int64(i))
		pubkeys = append(pubkeys, pk)
		privkeys = append(privkeys, sk)
	}
	return pubkeys, privkeys
}

var TestECDSAKeys = []*ecdsa.PrivateKey{
	crypto.UnmarshalHexECDSA("c71e183d51e9fae1d4fc410ca16a17a3a89da8e105b0e108576e2a77133f87b0"),
	crypto.UnmarshalHexECDSA("26c65dc72d016ebe50a5751c258d8ff3ddc3da40b5dcf7ac638619e041119b71"),
	crypto.UnmarshalHexECDSA("6a7e2b2ee79a8444d489c900f0c32bd944c88530882c9348b0d477b825773956"),
	crypto.UnmarshalHexECDSA("64d691d9af74ff28b23f38e49bedbae5aa5298933c477d60173a3990eb263481"),
	crypto.UnmarshalHexECDSA("376bb541ff3c913ea6b07cc0c405b991d354e140b4c3b8908884b84ef1475984"),
}
