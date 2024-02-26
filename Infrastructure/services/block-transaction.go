package services

import (
	ecdh "crypto/ecdh"
	ecdsa "crypto/ecdsa"
	elliptic "crypto/elliptic"
	rand "crypto/rand"
	sha256 "crypto/sha256"
	hex "encoding/hex"
	errors "errors"
	fmt "fmt"
	big "math/big"

	decimal "github.com/shopspring/decimal"
)

type BlockTransaction struct {
	FromAddress string
	ToAddress   string
	Amount      decimal.Decimal
	Signature   []byte
}

func NewBlockTransaction(fromAddress, toAddress string, amount decimal.Decimal) *BlockTransaction {
	return &BlockTransaction{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Amount:      amount,
	}
}

func (tx *BlockTransaction) SignTransaction(signingKey *ecdsa.PrivateKey) error {
	fromAddressKey, err := ConvertFromHexString(tx.FromAddress)
	if err != nil {
		return err
	}

	ephemeral, err := ecdh.GenerateKey(fromAddressKey.Curve, rand.Reader)
	if err != nil {
		return err
	}

	sharedKey, err := ecdh.GenerateSharedKey(ephemeral, &fromAddressKey, signingKey.PublicKey.X, signingKey.PublicKey.Y)
	if err != nil {
		return err
	}

	transactionHash := tx.CalculateHash()
	signature, err := ecdsa.Sign(rand.Reader, signingKey, transactionHash)
	if err != nil {
		return err
	}

	ephemeralPubKeyBytes := elliptic.Marshal(ephemeral.Curve, ephemeral.PublicKey.X, ephemeral.PublicKey.Y)
	tx.Signature = append(ephemeralPubKeyBytes, signature.R.Bytes()...)
	tx.Signature = append(tx.Signature, signature.S.Bytes()...)

	return nil
}

func (tx *BlockTransaction) IsValid() bool {
	if tx.FromAddress == "" {
		return false
	}

	if tx.Signature == nil || len(tx.Signature) == 0 {
		return false
	}

	fromAddressKey, err := ConvertFromHexString(tx.FromAddress)
	if err != nil {
		return false
	}

	ephemeralPubKeyBytes := tx.Signature[:65]
	signatureBytes := tx.Signature[65:]

	ephemeralPubKey, err := ConvertFromBytes(ephemeralPubKeyBytes)
	if err != nil {
		return false
	}

	sharedKey, err := ecdh.GenerateSharedKey(&fromAddressKey, ephemeralPubKey, fromAddressKey.X, fromAddressKey.Y)
	if err != nil {
		return false
	}

	transactionHash := tx.CalculateHash()
	signature := &ecdsa.Signature{
		R: new(big.Int).SetBytes(signatureBytes[:32]),
		S: new(big.Int).SetBytes(signatureBytes[32:]),
	}

	return ecdsa.Verify(&fromAddressKey, transactionHash, signature)
}

func (tx *BlockTransaction) CalculateHash() []byte {
	sha256Hash := sha256.New()
	data := fmt.Sprintf("%s%s%s", tx.FromAddress, tx.ToAddress, tx.Amount.String())
	sha256Hash.Write([]byte(data))
	return sha256Hash.Sum(nil)
}

func ConvertFromHexString(hexString string) (ecdsa.PublicKey, error) {
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		return ecdsa.PublicKey{}, err
	}

	return ConvertFromBytes(decoded)
}

func ConvertFromBytes(data []byte) (ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(elliptic.P256(), data)
	if x == nil {
		return ecdsa.PublicKey{}, errors.New("invalid byte data")
	}

	return ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, nil
}
