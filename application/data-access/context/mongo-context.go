package mongo_context

import (
	settings "bitshare-chain/infrastructure/settings"
	context "context"

	mongo "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoContext struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoContext(dbOptions *settings.MongoDbOptions) (*MongoContext, error) {
	clientOptions := options.Client()
	client, err := mongo.Connect(context.Background(), clientOptions.ApplyURI(dbOptions.ConnectionString))
	if err != nil {
		return nil, err
	}

	database := client.Database(dbOptions.DatabaseName)
	return &MongoContext{
		Client:   client,
		Database: database,
	}, nil
}

func (mongoContext *MongoContext) Close() error {
	if mongoContext.Client != nil {
		return mongoContext.Client.Disconnect(context.Background())
	}
	return nil
}

func (mongoContext *MongoContext) CreateIndex(collectionName string, indexKeys []string) error {
	collection := mongoContext.Database.Collection(collectionName)
	indexModel := mongo.IndexModel{
		Keys:    indexKeys,
		Options: options.Index().SetName("customIndexName"), // Customize the index name as needed
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}
