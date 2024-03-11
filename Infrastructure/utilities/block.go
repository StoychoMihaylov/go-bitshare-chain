package utilities

import (
	ecdsa "crypto/ecdsa"
	elliptic "crypto/elliptic"
	rand "crypto/rand"
	sha256 "crypto/sha256"
	hex "encoding/hex"
	fmt "fmt"
	hash "hash"
	big "math/big"
	strings "strings"
	time "time"
)

type BlockService interface {
	CalculateHash() string
	GetTransactionsStringData() string
	HasValidTransactions() bool
	MineBlock(difficulty int64)
	IsBlockSignatureValid() bool
	IsBlockMiner(minerAddress string) bool
	GetHash(hashAlgorithm hash.Hash, input string) string
}

type Block struct {
	TimeStamp      time.Time
	Transactions   []BlockTransaction
	PreviousHash   string
	Hash           string
	Nonce          int
	BlockMiner     string
	BlockSignature []byte
	SigningKey     *ecdsa.PrivateKey
}

func NewBlock(timeStamp time.Time, transactions []BlockTransaction, previousHash string, miningRewardAddress string, signingKey *ecdsa.PrivateKey) *Block {
	block := &Block{
		TimeStamp:    timeStamp,
		Transactions: transactions,
		PreviousHash: previousHash,
		Nonce:        0,
		BlockMiner:   miningRewardAddress,
		SigningKey:   signingKey,
	}
	block.Hash = block.CalculateHash()
	return block
}

func (block *Block) CalculateHash() string {
	sha256Hash := sha256.New()
	transactionsStringData := block.GetTransactionsStringData()
	data := fmt.Sprintf("%v%v%v%v", block.TimeStamp, block.PreviousHash, transactionsStringData, block.Nonce)
	return GetHash(sha256Hash, data)
}

func (block *Block) GetTransactionsStringData() string {
	var transactionsStringData strings.Builder
	for _, transaction := range block.Transactions {
		transactionsStringData.WriteString(hex.EncodeToString(transaction.Signature))
	}

	return transactionsStringData.String()
}

func (block *Block) HasValidTransactions() bool {
	for _, tx := range block.Transactions {
		if !tx.IsValid() {
			return false
		}
	}

	return true
}

func (block *Block) MineBlock(difficulty int) {
	stopWatch := time.Now()

	for !isValidHash(block.Hash, difficulty) {
		block.Nonce++
		block.Hash = block.CalculateHash()
	}

	r, s, err := ecdsa.Sign(rand.Reader, block.SigningKey, []byte(block.Hash))
	if err != nil {
		// Handle error
		return
	}

	block.BlockSignature = append(r.Bytes(), s.Bytes()...)
	result := block.IsBlockSignatureValid()

	if result {
		fmt.Printf("Block signature valid: %v\n", block.BlockSignature)
	}

	elapsed := time.Since(stopWatch)
	fmt.Printf("Mining time: %v for block: %v\n", elapsed, block.Hash)
}

func isValidHash(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

func (block *Block) IsBlockSignatureValid() bool {
	minerKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int),
		Y:     new(big.Int),
	}
	minerKey.X, minerKey.Y = minerKey.Curve.ScalarBaseMult(block.BlockSignature)

	if block.BlockSignature == nil || len(block.BlockSignature) != 64 {
		panic("Block doesn't contain a valid block signature.")
	}

	randomCordinate := new(big.Int).SetBytes(block.BlockSignature[:32])
	secret := new(big.Int).SetBytes(block.BlockSignature[32:])

	return ecdsa.Verify(minerKey, []byte(block.Hash), randomCordinate, secret)
}

func (block *Block) IsBlockMiner(minerAddress string) bool {
	minerKeyBytes, err := hex.DecodeString(minerAddress)
	if err != nil {
		panic("Invalid miner address.")
	}

	if block.BlockSignature == nil || len(block.BlockSignature) != 65 {
		panic("Block doesn't contain a valid block signature.")
	}

	// Extract r and s components from the signature
	randomCoordinate := new(big.Int).SetBytes(block.BlockSignature[:32])
	secret := new(big.Int).SetBytes(block.BlockSignature[32:])

	// Decode the minerKeyBytes into an ecdsa.PublicKey
	minerKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int),
		Y:     new(big.Int),
	}
	minerKey.X, minerKey.Y = elliptic.Unmarshal(minerKey.Curve, minerKeyBytes)

	return ecdsa.Verify(minerKey, []byte(block.Hash), randomCoordinate, secret)
}

func GetHash(hashAlgorithm hash.Hash, input string) string {
	hashAlgorithm.Write([]byte(input))
	hash := hashAlgorithm.Sum(nil)
	return hex.EncodeToString(hash)
}
