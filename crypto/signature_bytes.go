package crypto

import (
	"crypto/rand"
	bls12381 "github.com/kilic/bls12-381"
	log "github.com/sirupsen/logrus"
)

func mustBytesToG2(b []byte) *bls12381.PointG2 {
	ret, err := bls12381.NewG2().FromCompressed(b)
	if err != nil {
		log.Panicf("failed to construct G2 point from compressed bytes %x", b)
	}
	return ret
}

func mustBytesToG1(b []byte) *bls12381.PointG1 {
	ret, err := bls12381.NewG1().FromCompressed(b)
	if err != nil {
		log.Panicf("failed to construct G1 point from compressed bytes %x", b)
	}
	return ret
}

func bytes2Fr(b []byte) *bls12381.Fr {
	return bls12381.NewFr().FromBytes(b)
}

func GenKeyPairBytesFromSeed(seed int64) (priv, pub []byte) {
	privkey, pubkey := GenKeyPairFromSeed(seed)
	g1 := bls12381.NewG1()
	return privkey.ToBytes(), g1.ToCompressed(pubkey)
}

func GenKeyPairBytes() (priv, pub []byte) {
	g1 := bls12381.NewG1()
	fr := bls12381.NewFr()
	privkey := fr.Zero()
	if _, err := privkey.Rand(rand.Reader); err != nil {
		panic("failed to generate random secret key")
	}
	pubkey := g1.New()
	g1.MulScalar(pubkey, g1.One(), privkey)
	return privkey.ToBytes(), g1.ToCompressed(pubkey)
}

func SignBytes(privKey []byte, message []byte) []byte {
	sig := Sign(bytes2Fr(privKey), message)
	return bls12381.NewG2().ToCompressed(sig)
}

func VerifySigBytes(pubKey, msg, sig []byte) bool {
	return VerifySig(mustBytesToG1(pubKey), msg, mustBytesToG2(sig))
}

func AggPubKeysBytes(pubKeys [][]byte, multiplicities []uint32) []byte {
	var keys []*bls12381.PointG1
	g1 := bls12381.NewG1()
	for _, pubKey := range pubKeys {
		keys = append(keys, mustBytesToG1(pubKey))
	}
	return g1.ToCompressed(AggPubKeys(keys, multiplicities))
}

func VerifyAggSigBytes(pubKeys [][]byte, multiplicities []uint32, msg, aggSig []byte) bool {
	s := mustBytesToG2(aggSig)
	var keys []*bls12381.PointG1
	for _, pubKey := range pubKeys {
		keys = append(keys, mustBytesToG1(pubKey))
	}
	return VerifyAggSig(keys, multiplicities, msg, s)
}
