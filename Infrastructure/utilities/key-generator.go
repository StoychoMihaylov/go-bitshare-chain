package utilities

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	big "math/big"
)

type IKeyGenerator interface {
	GeneratePublicAndPrivateKey() (publicKey string, privateKey string)
}

type KeyGenerator struct{}

func (keyGenerator *KeyGenerator) GeneratePublicAndPrivateKey() (publicKey string, privateKey string) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(fmt.Errorf("failed to generate key: %v", err))
	}

	publicKeyBytes := elliptic.Marshal(key.PublicKey.Curve, key.PublicKey.X, key.PublicKey.Y)
	privateKeyBytes := key.D.Bytes()

	publicKey = hex.EncodeToString(publicKeyBytes)
	privateKey = hex.EncodeToString(privateKeyBytes)

	if !verifyKeys(key, publicKeyBytes, privateKeyBytes) {
		panic("generated keys do not correspond to each other")
	}

	return publicKey, privateKey
}

func verifyKeys(key *ecdsa.PrivateKey, publicKeyBytes, privateKeyBytes []byte) bool {
	curve := elliptic.P256()
	parsedPrivateKey := new(ecdsa.PrivateKey)
	parsedPrivateKey.Curve = curve
	parsedPrivateKey.D = new(big.Int).SetBytes(privateKeyBytes)
	parsedPrivateKey.PublicKey.Curve = curve
	parsedPrivateKey.PublicKey.X, parsedPrivateKey.PublicKey.Y = curve.ScalarBaseMult(privateKeyBytes)

	publicKeyFromPrivate := &parsedPrivateKey.PublicKey
	derivedPublicKeyBytes := elliptic.Marshal(key.PublicKey.Curve, publicKeyFromPrivate.X, publicKeyFromPrivate.Y)

	return hex.EncodeToString(publicKeyBytes) == hex.EncodeToString(derivedPublicKeyBytes)
}
