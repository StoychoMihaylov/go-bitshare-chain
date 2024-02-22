package services

//import "google.golang.org/genproto/googleapis/type/decimal"

//import "crypto/ecdsa"

type ChainService interface {
	MinePendingTransactions(miningRewardAddress string /* signingKey ecdsa */)
	AddTransaction( /* transaction BlockTransaction */ )
	GetBalanceOfAddress(address string) /* decimal */
	IsChainValid() bool
	CreateGenesisBlock() /* []Block */
	GetLatesBlock()      /* Block */
}

type service struct{}
