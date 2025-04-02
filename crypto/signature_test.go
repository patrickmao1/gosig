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
	msg1 := []byte("First message")
	msg2 := []byte("Second message")

	sig1 := Sign(privKey1, msg1)
	sig2 := Sign(privKey2, msg2)

	valid1 := VerifySig(pubKey1, msg1, sig1)
	require.True(t, valid1)
	valid2 := VerifySig(pubKey2, msg2, sig2)
	require.True(t, valid2)

	sigs := []*bls12381.PointG2{sig1, sig2}
	aggSig := AggregateSignatures(sigs)

	pubKeys := []*bls12381.PointG1{pubKey1, pubKey2}
	msgs := [][]byte{msg1, msg2}

	validAgg := VerifyAggregateSignature(pubKeys, msgs, aggSig)
	require.True(t, validAgg)

	pubKeysReversed := []*bls12381.PointG1{pubKey2, pubKey1}
	msgsReversed := [][]byte{msg2, msg1}

	validAggReversed := VerifyAggregateSignature(pubKeysReversed, msgsReversed, aggSig)
	require.True(t, validAggReversed)
}
