package services

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"
)

type BlockService interface {
	CalculateHash() string
	GetTransactionsStringData() string
	HasValidTransactions() bool
	MineBlock(difficulty int64)
	IsBlockSignatureValid() bool
	IsBlockMiner(minerAddress string) bool
	GetHash(hashAlgorithm string, input string) string //hashAlgorithm GO alternative of Cyrptography.HashAlgorithm in C#
}

type Block struct {
	TimeStamp      time.Time
	Transactions   []BlockTransaction
	Hash           string
	PreviousHash   string
	Nonce          int
	BlockSignature []byte
	BlockMiner     string
	SigningKey     *ecdsa.PrivateKey
}

func NewBlock(timeStamp time.Time, transactions []BlockTransaction, previousHash, miningRewardAddress string, signingKey *ecdsa.PrivateKey) *Block {
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

func (b *Block) CalculateHash() string {
	sha256Hash := sha256.New()
	transactionsStringData := b.GetTransactionsStringData()
	data := fmt.Sprintf("%v%v%v%v", b.TimeStamp, b.PreviousHash, transactionsStringData, b.Nonce)
	return GetHash(sha256Hash, data)
}

func (b *Block) GetTransactionsStringData() string {
	var transactionsStringData strings.Builder
	for _, transaction := range b.Transactions {
		transactionsStringData.WriteString(transaction.Signature)
	}
	return transactionsStringData.String()
}

func (b *Block) HasValidTransactions() bool {
	for _, tx := range b.Transactions {
		if !tx.IsValid() {
			return false
		}
	}
	return true
}

func (b *Block) MineBlock(difficulty int) {
	sw := time.Now()

	for !isValidHash(b.Hash, difficulty) {
		b.Nonce++
		b.Hash = b.CalculateHash()
	}

	b.BlockSignature, _ = ecdsa.Sign(b.SigningKey, []byte(b.Hash))
	result := b.IsBlockSignatureValid()

	elapsed := time.Since(sw)
	fmt.Printf("Mining time: %v for block: %v\n", elapsed, b.Hash)
}

func isValidHash(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

func (b *Block) IsBlockSignatureValid() bool {
	minerKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int),
		Y:     new(big.Int),
	}
	minerKey.X, minerKey.Y = minerKey.Curve.ScalarBaseMult(b.BlockSignature)

	if b.BlockSignature == nil {
		panic("Block doesn't contain block signature.")
	}

	return ecdsa.Verify(minerKey, []byte(b.Hash), b.BlockSignature)
}

func (b *Block) IsBlockMiner(minerAddress string) bool {
	fromAddresskey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int),
		Y:     new(big.Int),
	}
	fromAddresskey.X, fromAddresskey.Y = fromAddresskey.Curve.ScalarBaseMult(ConvertFromHexString(minerAddress))

	if b.BlockSignature == nil {
		panic("Block doesn't contain block signature.")
	}

	if fromAddresskey == nil {
		panic("Invalid miner address.")
	}

	return ecdsa.Verify(fromAddresskey, []byte(b.Hash), b.BlockSignature)
}

func GetHash(hashAlgorithm *sha256.Hash, input string) string {
	hashAlgorithm.Write([]byte(input))
	hash := hashAlgorithm.Sum(nil)
	return hex.EncodeToString(hash)
}

func ConvertFromHexString(hexString string) []byte {
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		panic(err)
	}
	return decoded
}
