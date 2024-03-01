package mongo_context

import (
	settings_options "bitshare-chain/infrastructure/options"
	context "context"

	mongo "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

// MongoContext represents the MongoDB context.
type MongoContext struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewMongoContext creates and initializes a new MongoContext.
func NewMongoContext(dbOptions *settings_options.MongoDbOptions) (*MongoContext, error) {
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

// Close closes the MongoDB client.
func (m *MongoContext) Close() error {
	if m.Client != nil {
		return m.Client.Disconnect(context.Background())
	}
	return nil
}

// CreateIndex creates an index on the specified collection.
func (m *MongoContext) CreateIndex(collectionName string, indexKeys []string) error {
	collection := m.Database.Collection(collectionName)
	indexModel := mongo.IndexModel{
		Keys:    indexKeys,
		Options: options.Index().SetName("customIndexName"), // Customize the index name as needed
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}
