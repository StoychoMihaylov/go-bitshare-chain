package services

type BlockService interface {
	CalculateHash() string
	GetTransactionsStringData() string
	HasValidTransactions() bool
	MineBlock(difficulty int64)
	IsBlockSignatureValid() bool
	IsBlockMiner(minerAddress string) bool
	GetHash(hashAlgorithm string, input string) string //hashAlgorithm GO alternative of Cyrptography.HashAlgorithm in C#
}

type service struct{}
