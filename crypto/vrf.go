package crypto

import (
	"golang.org/x/crypto/blake2b"
	"math/big"
)

func VRF(privKey, seed []byte) (rng uint32, proof []byte) {
	proof = SignBytes(privKey, seed)
	return RngFromProof(proof), proof
}

func RngFromProof(proof []byte) uint32 {
	rng := blake2b.Sum256(proof)
	return uint32(new(big.Int).SetBytes(rng[:32]).Uint64())
}
