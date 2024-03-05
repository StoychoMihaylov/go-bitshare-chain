package services

import (
	ecdsa "crypto/ecdsa"
	elliptic "crypto/elliptic"
	rand "crypto/rand"
	hex "encoding/hex"
	fmt "fmt"
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

	return publicKey, privateKey
}
