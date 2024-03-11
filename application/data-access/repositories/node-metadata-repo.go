package repositories

import (
	mongo_context "bitshare-chain/application/data-access/context"
	documents "bitshare-chain/application/data-access/documents"
	context "context"

	uuid "github.com/google/uuid"
	bson "go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"
)

type INodeMetadataRepository interface {
	InsertMetadataPublicKey(publicKey string) (string, error)
	GetRewardAddress(ctx context.Context) (*documents.NodeMetadataDocument, error)
	GetNodeConnections(ctx context.Context) ([]documents.NodeConnectionsSubDocument, error)
	UpdateNodeMetadataConnections(ctx context.Context, nodeMetadataSubDocuments []documents.NodeConnectionsSubDocument) error
}

type NodeMetadataRepository struct {
	nodeMetadata *mongo.Collection
}

func NewNodeMetadataRepository(mongoContext *mongo_context.MongoContext) *NodeMetadataRepository {
	nodeMetadata := mongoContext.Database.Collection("NodeMetadataDocument")
	return &NodeMetadataRepository{
		nodeMetadata: nodeMetadata,
	}
}

func (repo *NodeMetadataRepository) InsertMetadataPublicKey(publicKey string) (string, error) {
	filter := bson.M{}
	nodeMetadata := &documents.NodeMetadataDocument{}
	err := repo.nodeMetadata.FindOne(context.Background(), filter).Decode(nodeMetadata)

	if err == mongo.ErrNoDocuments {
		newNodeMetadata := &documents.NodeMetadataDocument{
			NodeId:        uuid.NewString(),
			RewardAddress: publicKey,
		}

		_, err := repo.nodeMetadata.InsertOne(context.Background(), newNodeMetadata)
		if err != nil {
			return "", err
		}
		return newNodeMetadata.NodeId, nil
	} else if err == nil {
		update := bson.M{"$set": bson.M{"RewardAddress": publicKey}}
		_, err := repo.nodeMetadata.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return "", err
		}
		return nodeMetadata.NodeId, nil
	}

	return "", err
}

func (repo *NodeMetadataRepository) GetRewardAddress(ctx context.Context) (*documents.NodeMetadataDocument, error) {
	nodeMetadata := &documents.NodeMetadataDocument{}
	err := repo.nodeMetadata.FindOne(ctx, bson.M{}).Decode(nodeMetadata)
	if err != nil {
		return nil, err
	}
	return nodeMetadata, nil
}

func (repo *NodeMetadataRepository) GetNodeConnections(ctx context.Context) ([]documents.NodeConnectionsSubDocument, error) {
	nodeMetadata, err := repo.GetRewardAddress(ctx)
	if err != nil {
		return nil, err
	}

	if nodeMetadata.NodeConnections == nil {
		return []documents.NodeConnectionsSubDocument{}, nil
	}

	return nodeMetadata.NodeConnections, nil
}

func (repo *NodeMetadataRepository) UpdateNodeMetadataConnections(ctx context.Context, nodeMetadataSubDocuments []documents.NodeConnectionsSubDocument) error {
	nodeMetadata, err := repo.GetRewardAddress(ctx)
	if err != nil {
		return err
	}

	update := bson.M{"$push": bson.M{"NodeConnections": bson.M{"$each": nodeMetadataSubDocuments}}}
	_, err = repo.nodeMetadata.UpdateOne(ctx, bson.M{"_id": nodeMetadata.Id}, update)
	return err
}
