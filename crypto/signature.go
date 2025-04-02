package crypto

import (
	"crypto/rand"
	bls12381 "github.com/kilic/bls12-381"
	"golang.org/x/crypto/sha3"
)

// GenKeyPair creates a new BLS key pair
func GenKeyPair() (priv *bls12381.Fr, pub *bls12381.PointG1) {
	// Initialize the curve
	g1 := bls12381.NewG1()
	fr := bls12381.NewFr()
	// Generate a random secret key
	privKey := fr.Zero()
	if _, err := privKey.Rand(rand.Reader); err != nil {
		panic("failed to generate random secret key")
	}
	// Compute public key as g1 * secret key
	pubKey := g1.New()
	g1.MulScalar(pubKey, g1.One(), privKey)
	return privKey, pubKey
}

// hashToG2 hashes a message to a point on PointG2 using Keccak256
func hashToG2(message []byte) *bls12381.PointG2 {
	h := sha3.NewLegacyKeccak256()
	_, err := h.Write(message)
	if err != nil {
		panic(err)
	}
	sha3.Sum256(message)
	digest := h.Sum(nil)
	domainSeparationTag := []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_NUL_")
	g2 := bls12381.NewG2()
	messagePoint, err := g2.HashToCurve(digest, domainSeparationTag)
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
	negSignature := g2.New()
	g2.Neg(negSignature, sig)

	e.AddPair(g1.One(), negSignature)
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

func AggPubKeys(pubKeys []*bls12381.PointG1) *bls12381.PointG1 {
	if len(pubKeys) == 0 {
		panic("no public keys to aggregate")
	}
	g1 := bls12381.NewG1()
	aggKey := g1.New()
	g1.Add(aggKey, g1.Zero(), pubKeys[0])
	for i := 1; i < len(pubKeys); i++ {
		g1.Add(aggKey, aggKey, pubKeys[i])
	}
	return aggKey
}

func VerifyAggSig(
	publicKeys []*bls12381.PointG1, messages [][]byte, aggSig *bls12381.PointG2) bool {

	if len(publicKeys) != len(messages) {
		panic("number of public keys must match number of messages")
	}
	e := bls12381.NewEngine()
	g2 := bls12381.NewG2()
	// Add pairs of (publicKey_i, H(message_i))
	for i := 0; i < len(publicKeys); i++ {
		e.AddPair(publicKeys[i], hashToG2(messages[i]))
	}
	// Add the negative of the signature with generator point
	g1 := bls12381.NewG1()
	negSig := g2.New()
	g2.Neg(negSig, aggSig)
	e.AddPair(g1.One(), negSig)
	return e.Check()
}
