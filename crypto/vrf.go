package crypto

import (
	bls12381 "github.com/kilic/bls12-381"
	"golang.org/x/crypto/sha3"
	"math/big"
)

func GenRNG(privKey *bls12381.Fr, seed []byte) (rng uint32, proof []byte) {
	sig := Sign(privKey, seed)
	num := sha3.Sum256(proof)
	rng = uint32(new(big.Int).SetBytes(num[:32]).Int64())
	return rng, sig
}

func VerifyRNG(pubKey *bls12381.PointG1, proof, signData []byte) bool {
	return false
}
