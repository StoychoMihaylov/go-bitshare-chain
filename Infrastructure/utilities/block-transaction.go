package services

import (
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

func (transaction *BlockTransaction) SignTransaction(signingKey *ecdsa.PrivateKey) error {
	fromAddressKey, err := ConvertFromHexString(transaction.FromAddress)
	if err != nil {
		return err
	}

	ephemeral, err := ecdsa.GenerateKey(fromAddressKey.Curve, rand.Reader)
	if err != nil {
		return err
	}

	transactionHash := transaction.CalculateHash()
	r, s, err := ecdsa.Sign(rand.Reader, signingKey, transactionHash)
	if err != nil {
		return err
	}

	ephemeralPubKeyBytes := elliptic.Marshal(ephemeral.Curve, ephemeral.PublicKey.X, ephemeral.PublicKey.Y)
	transaction.Signature = append(ephemeralPubKeyBytes, r.Bytes()...)
	transaction.Signature = append(transaction.Signature, s.Bytes()...)

	return nil
}

func (transaction *BlockTransaction) IsValid() bool {
	if transaction.FromAddress == "" {
		return false
	}

	if transaction.Signature == nil || len(transaction.Signature) == 0 {
		return false
	}

	fromAddressKey, err := ConvertFromHexString(transaction.FromAddress)
	if err != nil {
		return false
	}

	// Assuming sharedKey is intended to be fromAddressKey
	transactionHash := transaction.CalculateHash()

	// Assuming signature is intended to be a variable of type ecdsa.Signature
	signature := &ecdsa.Signature{
		R: new(big.Int).SetBytes(transaction.Signature[:32]),
		S: new(big.Int).SetBytes(transaction.Signature[32:]),
	}

	return ecdsa.Verify(&fromAddressKey, transactionHash, signature)
}

func (transaction *BlockTransaction) CalculateHash() []byte {
	sha256Hash := sha256.New()
	data := fmt.Sprintf("%s%s%s", transaction.FromAddress, transaction.ToAddress, transaction.Amount.String())
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
