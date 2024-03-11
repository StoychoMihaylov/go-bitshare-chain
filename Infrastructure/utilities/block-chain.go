package utilities

import (
	ecdsa "crypto/ecdsa"
	elliptic "crypto/elliptic"
	rand "crypto/rand"
	fmt "fmt"
	time "time"

	decimal "github.com/shopspring/decimal"
)

type Blockchain struct {
	Chain               []Block
	Difficulty          int
	PendingTransactions []BlockTransaction
	MiningReward        decimal.Decimal
}

func NewBlockchain() *Blockchain {
	blockchain := &Blockchain{
		Chain:               createGenesisBlock(),
		Difficulty:          2,
		PendingTransactions: make([]BlockTransaction, 0),
		MiningReward:        decimal.NewFromFloat(100),
	}
	return blockchain
}

// MinePendingTransactions mines pending transactions and adds a new block to the blockchain.
func (blockChain *Blockchain) MinePendingTransactions(miningRewardAddress string, signingKey *ecdsa.PrivateKey) {
	block := NewBlock(time.Now(), blockChain.PendingTransactions, getLatestBlockHash(blockChain.Chain), miningRewardAddress, signingKey)
	block.MineBlock(blockChain.Difficulty)
	blockChain.Chain = append(blockChain.Chain, *block)

	// Reset pending transactions and create a new transaction to send the miner a reward
	blockChain.PendingTransactions = []BlockTransaction{
		{ToAddress: miningRewardAddress, Amount: blockChain.MiningReward},
	}
}

func (blockChain *Blockchain) AddTransaction(transaction BlockTransaction) {
	if transaction.FromAddress == "" || transaction.ToAddress == "" {
		panic("Transaction must include from and to address.")
	}

	if !transaction.IsValid() {
		panic("Cannot add invalid transaction to the chain.")
	}

	blockChain.PendingTransactions = append(blockChain.PendingTransactions, transaction)
}

func (blockChain *Blockchain) GetBalanceOfAddress(address string) decimal.Decimal {
	balance := decimal.NewFromFloat(0)

	for _, block := range blockChain.Chain {
		for _, transaction := range block.Transactions {
			if transaction.FromAddress == address {
				balance = balance.Sub(transaction.Amount)
			}

			if transaction.ToAddress == address {
				balance = balance.Add(transaction.Amount)
			}
		}
	}

	return balance
}

func (blockChain *Blockchain) IsChainValid() bool {
	for index := 1; index < len(blockChain.Chain); index++ {
		currentBlock := &blockChain.Chain[index]
		previousBlock := &blockChain.Chain[index-1]

		if !currentBlock.IsBlockSignatureValid() {
			return false
		}

		if !currentBlock.HasValidTransactions() {
			return false
		}

		if currentBlock.Hash != currentBlock.CalculateHash() {
			return false
		}

		if currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}
	}

	return true
}

func createGenesisBlock() []Block {
	return []Block{
		{
			TimeStamp:    time.Now(),
			Transactions: []BlockTransaction{},
			PreviousHash: "",
			BlockMiner:   "",
			SigningKey:   createSigningKey(),
		},
	}
}

func getLatestBlockHash(chain []Block) string {
	return chain[len(chain)-1].Hash
}

func createSigningKey() *ecdsa.PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(fmt.Errorf("failed to generate key: %v", err))
	}
	return key
}
