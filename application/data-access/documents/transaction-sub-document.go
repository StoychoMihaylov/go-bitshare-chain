package documents

import (
	ecdsa "crypto/ecdsa"
	elliptic "crypto/elliptic"
	sha256 "crypto/sha256"
	hex "encoding/hex"
	fmt "fmt"
	big "math/big"
	strings "strings"
	time "time"
)

type TransactionSubDocument struct {
	FromAddress string    `bson:"fromAddress,omitempty"`
	ToAddress   string    `bson:"toAddress,omitempty"`
	Amount      float64   `bson:"amount,omitempty"`
	TimeStamp   time.Time `bson:"timeStamp,omitempty"`
	Signature   []byte    `bson:"signature,omitempty"`
}

func (transaction *TransactionSubDocument) SignTransaction(signingKey *ecdsa.PrivateKey) {
	if hex.EncodeToString(signingKey.PublicKey.X.Bytes())+hex.EncodeToString(signingKey.PublicKey.Y.Bytes()) != tx.FromAddress {
		panic("You cannot sign transactions for other wallets!")
	}

	transactionHash := transaction.CalculateHash()
	r, s, err := ecdsa.Sign(strings.NewReader("random"), signingKey, transactionHash)
	if err != nil {
		panic(fmt.Errorf("failed to sign transaction: %v", err))
	}

	transaction.Signature = append(r.Bytes(), s.Bytes()...)
}

func (transaction *TransactionSubDocument) IsValid() bool {
	if transaction.FromAddress == "" {
		return false
	}

	if transaction.Signature == nil || len(transaction.Signature) != 64 {
		panic("Invalid signature in this transaction.")
	}

	randomCoordinate := new(big.Int).SetBytes(transaction.Signature[:32])
	secret := new(big.Int).SetBytes(transaction.Signature[32:])

	fromAddressKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int),
		Y:     new(big.Int),
	}
	fromAddressKey.X, fromAddressKey.Y = fromAddressKey.Curve.ScalarBaseMult(ConvertFromHexString(transaction.FromAddress))

	if fromAddressKey == nil {
		return false
	}

	return ecdsa.Verify(fromAddressKey, transaction.CalculateHash(), randomCoordinate, secret)
}

func (transaction *TransactionSubDocument) CalculateHash() []byte {
	sha256Hash := sha256.New()
	sha256Hash.Write([]byte(transaction.FromAddress + transaction.ToAddress + fmt.Sprintf("%.8f", transaction.Amount)))
	return sha256Hash.Sum(nil)
}

func ConvertFromHexString(hexString string) []byte {
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		panic(err)
	}
	return decoded
}
