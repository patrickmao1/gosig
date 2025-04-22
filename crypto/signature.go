package crypto

import (
	"crypto/rand"
	bls12381 "github.com/kilic/bls12-381"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/blake2b"
	"math/big"
	mathrand "math/rand"
)

func GenKeyPairFromSeed(seed int64) (priv *bls12381.Fr, pub *bls12381.PointG1) {
	g1 := bls12381.NewG1()
	fr := bls12381.NewFr()
	privKey := fr.Zero()
	rng := mathrand.New(mathrand.NewSource(seed))
	if _, err := privKey.Rand(rng); err != nil {
		panic("failed to generate random secret key")
	}
	pubKey := g1.New()
	g1.MulScalar(pubKey, g1.One(), privKey)
	return privKey, pubKey
}

func GenKeyPair() (priv *bls12381.Fr, pub *bls12381.PointG1) {
	g1 := bls12381.NewG1()
	fr := bls12381.NewFr()
	privKey := fr.Zero()
	if _, err := privKey.Rand(rand.Reader); err != nil {
		panic("failed to generate random secret key")
	}
	pubKey := g1.New()
	g1.MulScalar(pubKey, g1.One(), privKey)
	return privKey, pubKey
}

func hashToG2(message []byte) *bls12381.PointG2 {
	digest := blake2b.Sum256(message)
	dst := []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_NUL_")
	g2 := bls12381.NewG2()
	messagePoint, err := g2.HashToCurve(digest[:], dst)
	if err != nil {
		panic(err)
	}
	return messagePoint
}

func Sign(privKey *bls12381.Fr, message []byte) *bls12381.PointG2 {
	g2 := bls12381.NewG2()
	sig := g2.New()
	g2.MulScalar(sig, hashToG2(message), privKey)
	return sig
}

func VerifySig(pubKey *bls12381.PointG1, msg []byte, sig *bls12381.PointG2) bool {
	e := bls12381.NewEngine()
	g1 := bls12381.NewG1()
	g2 := bls12381.NewG2()

	// e(g1, signature) == e(publicKey, messagePoint)
	negSig := g2.New()
	g2.Neg(negSig, sig)

	e.AddPair(g1.One(), negSig)
	e.AddPair(pubKey, hashToG2(msg))
	result := e.Check()

	return result
}

func AggSigs(sigs []*bls12381.PointG2) *bls12381.PointG2 {
	g2 := bls12381.NewG2()
	aggSig := g2.Zero()
	for _, sig := range sigs {
		g2.Add(aggSig, aggSig, sig)
	}
	return aggSig
}

func AggPubKeys(pubKeys []*bls12381.PointG1, multiplicities []uint32) *bls12381.PointG1 {
	if len(pubKeys) == 0 {
		panic("no public keys to aggregate")
	}
	if len(pubKeys) != len(multiplicities) {
		log.Panicf("len(pubKeys) != len(multiplicities): %d != %d", len(pubKeys), len(multiplicities))
	}
	g1 := bls12381.NewG1()
	aggPubKey := g1.Zero()
	for i, pubKey := range pubKeys {
		multiPubKey := g1.New()
		multiplicity := new(big.Int).SetUint64(uint64(multiplicities[i]))
		g1.MulScalarBig(multiPubKey, pubKey, multiplicity)
		g1.Add(aggPubKey, aggPubKey, multiPubKey)
	}
	return aggPubKey
}

func VerifyAggSig(pubKeys []*bls12381.PointG1, multiplicities []uint32, msg []byte, aggSig *bls12381.PointG2) bool {
	if len(pubKeys) != len(multiplicities) {
		log.Panicf("len(pubKeys) != len(multiplicities): %d != %d", len(pubKeys), len(multiplicities))
	}
	aggPubKey := AggPubKeys(pubKeys, multiplicities)
	return VerifySig(aggPubKey, msg, aggSig)
}
