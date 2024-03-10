package documents

type NodeConnectionsSubDocument struct {
	NodeID     string `bson:"nodeId,omitempty"`
	NodeURL    string `bson:"nodeUrl,omitempty"`
	NodeHealth bool   `bson:"nodeHealth,omitempty"`
}
