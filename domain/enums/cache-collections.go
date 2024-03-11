package enums

type CacheCollections int

const (
	MinerPrivateKey CacheCollections = iota
	NodeMetadata
	PendingTransactions
)
