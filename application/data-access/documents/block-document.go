package documents

import (
	ecdsa "crypto/ecdsa"
	elliptic "crypto/elliptic"
	rand "crypto/rand"
	sha256 "crypto/sha256"
	hex "encoding/hex"
	errors "errors"
	fmt "fmt"
	hash "hash"
	big "math/big"
	strings "strings"
	time "time"
)

type BlockDocument struct {
	ID                string                   `bson:"_id,omitempty"`
	Index             int64                    `bson:"index,omitempty"`
	TimeStamp         time.Time                `bson:"timeStamp,omitempty"`
	Transactions      []TransactionSubDocument `bson:"transactions,omitempty"`
	Hash              string                   `bson:"hash,omitempty"`
	PreviousHash      string                   `bson:"previousHash,omitempty"`
	Nonce             int                      `bson:"nonce,omitempty"`
	BlockSignature    []byte                   `bson:"blockSignature,omitempty"`
	BlockMinerAddress string                   `bson:"blockMinerAddress,omitempty"`
}

func NewBlockDocument() *BlockDocument {
	return &BlockDocument{
		Transactions: []TransactionSubDocument{},
	}
}

func (block *BlockDocument) AddTransaction(transaction TransactionSubDocument) error {
	if transaction.FromAddress == "" || transaction.ToAddress == "" {
		return errors.New("Transaction must include from and to address")
	}

	if !transaction.IsValid() {
		return errors.New("Cannot add invalid transaction to the chain")
	}

	block.Transactions = append(block.Transactions, transaction)
	return nil
}

func (block *BlockDocument) MineBlock(difficulty int, signingKey *ecdsa.PrivateKey) {
	sw := time.Now()

	block.Hash = block.CalculateHash()
	for block.Hash[:difficulty] != strings.Join([]string{"0"}, strings.Repeat("0", difficulty)) {
		block.Nonce++
		block.Hash = block.CalculateHash()
	}

	hashBytes, _ := hex.DecodeString(block.Hash)
	block.BlockSignature, _ = signingKey.Sign(rand.Reader, hashBytes, nil)

	elapsed := time.Since(sw)
	fmt.Printf("Mining time: %s for block: %s\n", elapsed, block.Hash)
}

func (block *BlockDocument) CalculateHash() string {
	transactionsStringData := block.GetTransactionsStringData()
	hashInput := fmt.Sprintf("%v%v%v%v", block.TimeStamp, block.PreviousHash, transactionsStringData, block.Nonce)
	return GetHash(sha256.New(), hashInput)
}

func (block *BlockDocument) GetTransactionsStringData() string {
	var transactionsStringData strings.Builder
	for _, transaction := range block.Transactions {
		signatureString := hex.EncodeToString(transaction.Signature)
		transactionsStringData.WriteString(signatureString)
	}
	return transactionsStringData.String()
}

func (block *BlockDocument) IsSignatureValid(blockMinerAddress string) bool {
	minerAddressBytes, err := hex.DecodeString(blockMinerAddress)
	if err != nil {
		fmt.Println("Error decoding hex string:", err)
		return false
	}

	// Unmarshal the public key
	var minerKey ecdsa.PublicKey
	minerKey.X, minerKey.Y = elliptic.Unmarshal(elliptic.P256(), minerAddressBytes)
	minerKey.Curve = elliptic.P256()

	if block.BlockSignature == nil {
		fmt.Println("Block doesn't contain block signature.")
		return false
	}

	r := new(big.Int).SetBytes(block.BlockSignature[:32])
	s := new(big.Int).SetBytes(block.BlockSignature[32:])

	return ecdsa.Verify(&minerKey, []byte(block.Hash), r, s)
}

func (block *BlockDocument) HasValidTransactions() bool {
	for _, tx := range block.Transactions {
		if !tx.IsValid() {
			return false
		}
	}
	return true
}

func GetHash(hashAlgorithm hash.Hash, input string) string {
	hash := hashAlgorithm.Sum([]byte(input))
	return hex.EncodeToString(hash)
}
