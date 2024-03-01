package documents

import "go.mongodb.org/mongo-driver/bson/primitive"

type WalletAccountDocument struct {
	ID           primitive.ObjectID       `bson:"_id,omitempty"`
	Address      string                   `bson:"address,omitempty"`
	Transactions []TransactionSubDocument `bson:"transactions,omitempty"`
}
