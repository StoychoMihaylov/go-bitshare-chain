package bindingmodels

import (
	ecdsa "crypto/ecdsa"
	sha256 "crypto/sha256"
	errors "errors"
	fmt "fmt"
	big "math/big"
)

type TransactionBindingModel struct {
	FromAddress string  `json:"fromAddress"`
	ToAddress   string  `json:"toAddress"`
	Amount      float64 `json:"amount"`
	Signature   []byte  `json:"signature,omitempty"`
}

func (transaction *TransactionBindingModel) SignTransaction(signingKey *ecdsa.PrivateKey) error {
	if transaction.FromAddress != fmt.Sprintf("%x", signingKey.X) {
		return errors.New("you cannot sign transactions for other wallets")
	}

	transactionHash := transaction.CalculateHash()
	r, s, err := ecdsa.Sign(nil, signingKey, transactionHash)
	if err != nil {
		return err
	}
	signature, err := MarshalECDSASignature(r, s)
	if err != nil {
		return err
	}
	transaction.Signature = signature
	return nil
}

func (transaction *TransactionBindingModel) CalculateHash() []byte {
	hash := sha256.Sum256([]byte(transaction.FromAddress + transaction.ToAddress + fmt.Sprintf("%.2f", transaction.Amount)))
	return hash[:]
}

func MarshalECDSASignature(r, s *big.Int) ([]byte, error) {
	rb := r.Bytes()
	sb := s.Bytes()

	signature := make([]byte, 64)
	copy(signature[32-len(rb):32], rb)
	copy(signature[64-len(sb):], sb)
	return signature, nil
}
