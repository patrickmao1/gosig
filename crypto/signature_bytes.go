package crypto

import (
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

func SignBytes(privKey []byte, message []byte) []byte {
	sig := Sign(bytes2Fr(privKey), message)
	return bls12381.NewG2().ToCompressed(sig)
}

func VerifySigBytes(pubKey, msg, sig []byte) bool {
	return VerifySig(mustBytesToG1(pubKey), msg, mustBytesToG2(sig))
}

func AggSigBytes(sigs [][]byte) []byte {
	var sigG2s []*bls12381.PointG2
	g2 := bls12381.NewG2()
	for _, sig := range sigs {
		sigG2s = append(sigG2s, mustBytesToG2(sig))
	}
	aggSig := AggSigs(sigG2s)
	return g2.ToCompressed(aggSig)
}

func AggPubKeysBytes(pubKeys [][]byte) []byte {
	var keys []*bls12381.PointG1
	g1 := bls12381.NewG1()
	for _, pubKey := range pubKeys {
		keys = append(keys, mustBytesToG1(pubKey))
	}
	return g1.ToCompressed(AggPubKeys(keys))
}
