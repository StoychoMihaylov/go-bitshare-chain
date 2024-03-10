package documents

import primitive "go.mongodb.org/mongo-driver/bson/primitive"

type NodeMetadataDocument struct {
	ID              primitive.ObjectID           `bson:"_id,omitempty"`
	NodeID          string                       `bson:"nodeId,omitempty"`
	RewardAddress   string                       `bson:"rewardAddress,omitempty"`
	NodeConnections []NodeConnectionsSubDocument `bson:"nodeConnections,omitempty"`
}
