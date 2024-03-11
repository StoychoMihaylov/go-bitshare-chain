package documents

import primitive "go.mongodb.org/mongo-driver/bson/primitive"

type NodeMetadataDocument struct {
	Id              primitive.ObjectID           `bson:"_id,omitempty"`
	NodeId          string                       `bson:"nodeId,omitempty"`
	RewardAddress   string                       `bson:"rewardAddress,omitempty"`
	NodeConnections []NodeConnectionsSubDocument `bson:"nodeConnections,omitempty"`
}
