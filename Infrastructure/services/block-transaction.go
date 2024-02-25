package services

/* import "crypto/ecdsa" */
import (
	decimal "github.com/shopspring/decimal"
)

type BlockTransactionService interface {
	SignTransaction( /* signingKey ecdsa */ )
	IsValid() bool
	CalculateHash() []byte
}

type BlockTransaction struct {
	FromAddress string
	ToAddress   string
	Amount      decimal.Decimal
	Signature   []byte
}
