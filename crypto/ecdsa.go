package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/blake2b"
	"math/big"
)

func GenKeyECDSA() *ecdsa.PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return key
}

func SignECDSA(key *ecdsa.PrivateKey, msg []byte) []byte {
	hash := blake2b.Sum256(msg)
	sig, err := ecdsa.SignASN1(rand.Reader, key, hash[:])
	if err != nil {
		panic(err)
	}
	return sig
}

var sigCache = make(map[string]bool)

func VerifyECDSABytes(key, msg, sig []byte) bool {
	mapKey := string(sig) + string(msg)
	if sigCache[mapKey] {
		return true
	}
	pk := UnmarshalECDSAPublic(key)
	pass := VerifyECDSA(pk, msg, sig)
	sigCache[mapKey] = pass
	return pass
}

func VerifyECDSA(key *ecdsa.PublicKey, msg []byte, sig []byte) bool {
	hash := blake2b.Sum256(msg)
	return ecdsa.VerifyASN1(key, hash[:], sig)
}

func MarshalECDSAPublic(key *ecdsa.PublicKey) []byte {
	return elliptic.MarshalCompressed(key.Curve, key.X, key.Y)
}

var keyCache = map[string]*ecdsa.PublicKey{}

func UnmarshalECDSAPublic(pubkey []byte) (key *ecdsa.PublicKey) {
	if key, ok := keyCache[string(pubkey)]; ok {
		return key
	}
	key = new(ecdsa.PublicKey)
	key.Curve = elliptic.P256()
	key.X, key.Y = elliptic.UnmarshalCompressed(key.Curve, pubkey)
	keyCache[string(pubkey)] = key
	return key
}

func MarshalECDSA(key *ecdsa.PrivateKey) []byte {
	return key.D.Bytes()
}

func UnmarshalHexECDSA(hx string) (key *ecdsa.PrivateKey) {
	bs, err := hex.DecodeString(hx)
	if err != nil {
		panic(err)
	}
	return UnmarshalECDSA(bs)
}

func UnmarshalECDSA(bs []byte) (key *ecdsa.PrivateKey) {
	if len(bs) != 32 {
		panic("bad key length")
	}
	key = new(ecdsa.PrivateKey)
	key.Curve = elliptic.P256()
	key.D = new(big.Int).SetBytes(bs)
	key.PublicKey.X, key.PublicKey.Y = key.Curve.ScalarBaseMult(bs)
	return key
}
