package crypto

import (
	bls12381 "github.com/kilic/bls12-381"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	privKey, pubKey := GenKeyPair()
	privKey2, pubKey2 := GenKeyPair()
	require.NotNil(t, privKey)
	require.NotNil(t, privKey2)
	require.NotNil(t, pubKey)
	require.NotNil(t, pubKey2)
	if privKey.Equal(privKey2) {
		t.Error("equal privKeys")
	}
}

func TestSignAndVerify(t *testing.T) {
	privKey, pubKey := GenKeyPair()
	message := []byte("hello")
	sig := Sign(privKey, message)
	valid := VerifySig(pubKey, message, sig)
	require.True(t, valid)
	wrongMessage := []byte("Wrong message")
	valid = VerifySig(pubKey, wrongMessage, sig)
	require.False(t, valid)
	_, wrongPubKey := GenKeyPair()
	valid = VerifySig(wrongPubKey, message, sig)
	require.False(t, valid)
}

func TestSignAndVerifyEmptyMessage(t *testing.T) {
	privKey, pubKey := GenKeyPair()
	emptyMessage := []byte{}
	sig := Sign(privKey, emptyMessage)
	valid := VerifySig(pubKey, emptyMessage, sig)
	require.True(t, valid)
}

func TestSignatureAggregation(t *testing.T) {
	// Generate key pairs
	privKey1, pubKey1 := GenKeyPair()
	privKey2, pubKey2 := GenKeyPair()
	msg := []byte("msg")

	sig1 := Sign(privKey1, msg)
	sig2 := Sign(privKey2, msg)

	valid1 := VerifySig(pubKey1, msg, sig1)
	require.True(t, valid1)
	valid2 := VerifySig(pubKey2, msg, sig2)
	require.True(t, valid2)

	sigs := []*bls12381.PointG2{sig1, sig2}
	aggSig := AggSigs(sigs)

	pubKeys := []*bls12381.PointG1{pubKey1, pubKey2}
	mul := []uint32{1, 1}

	// normal case
	validAgg := VerifyAggSig(pubKeys, mul, msg, aggSig)
	require.True(t, validAgg)

	// reversing signature order should also work
	pubKeysReversed := []*bls12381.PointG1{pubKey2, pubKey1}
	validAggReversed := VerifyAggSig(pubKeysReversed, mul, msg, aggSig)
	require.True(t, validAggReversed)

	// should catch incorrect multiplicities
	invalidMul := VerifyAggSig(pubKeys, []uint32{1, 0}, msg, aggSig)
	require.False(t, invalidMul)

	// adding sig1 to aggSig multiple times
	aggSig = AggSigs([]*bls12381.PointG2{aggSig, sig1})
	validMul := VerifyAggSig(pubKeys, []uint32{2, 1}, msg, aggSig)
	require.True(t, validMul)
}
