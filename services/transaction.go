package services

/* import "crypto/ecdsa" */

type TransactionService interface {
	SignTransaction( /* signingKey ecdsa */ )
	IsValid() bool
	CalculateHash() []byte
}

type service struct{}
