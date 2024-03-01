package repositories

import (
	context "context"
	fmt "fmt"

	mongo_context "github.com/bitshare/application/dataaccess/context"
	documents "github.com/bitshare/application/dataaccess/documents"
	bson "go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"
)

type WalletAccountRepository interface {
	CreateWalletAccount(newWalletAccount *documents.WalletAccountDocument) error
	UpdateWalletAccountTransactions(walletTransactions map[string][]documents.TransactionSubDocument) error
}

type walletAccountRepository struct {
	walletAccountCollection *mongo.Collection
}

func NewWalletAccountRepository(mongoContext *mongo_context.MongoContext) WalletAccountRepository {
	return &walletAccountRepository{
		walletAccountCollection: mongoContext.Database.Collection("WalletAccountDocument"),
	}
}

func (r *walletAccountRepository) CreateWalletAccount(newWalletAccount *documents.WalletAccountDocument) error {
	_, err := r.walletAccountCollection.InsertOne(context.Background(), newWalletAccount)
	return err
}

func (r *walletAccountRepository) UpdateWalletAccountTransactions(walletTransactions map[string][]documents.TransactionSubDocument) error {
	for address, transactions := range walletTransactions {
		filter := bson.M{"address": address}
		update := bson.M{"$set": bson.M{"transactions": transactions}}

		_, err := r.walletAccountCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return fmt.Errorf("failed to update wallet account transactions: %v", err)
		}
	}

	return nil
}
