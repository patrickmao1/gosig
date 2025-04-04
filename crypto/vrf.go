package crypto

import (
	"golang.org/x/crypto/sha3"
	"math/big"
)

func GenRNG(privKey, seed []byte) (rng uint32, proof []byte) {
	sig := SignBytes(privKey, seed)
	num := sha3.Sum256(proof)
	rng = uint32(new(big.Int).SetBytes(num[:32]).Int64())
	return rng, sig
}

func GenRNGWithProof(proof []byte) uint32 {
	rng := sha3.Sum256(proof)
	return uint32(new(big.Int).SetBytes(rng[:32]).Uint64())
}

func VerifyRNG(pubKey, proof, signData []byte) bool {
	return VerifySigBytes(pubKey, signData, proof)
}
