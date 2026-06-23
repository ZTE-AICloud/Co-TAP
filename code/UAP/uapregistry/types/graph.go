package types

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	// 节点属性列表
	NodeProperty_AgentCardName = "agentCardName" // name in  agent card
	NodeProperty_ServiceName   = "serviceName"   // name in register center
	NodeProperty_Namespace     = "namespace"     //
	NodeProperty_Cluster       = "cluster"
)

type DatabaseConfig struct {
	URI      string
	Username string
	Password string
	Database string
}

type RelatedNode struct {
	Relationship neo4j.Relationship `json:"relationship"`
	Node         neo4j.Node         `json:"node"`
}

type Graph struct {
	Relationships []neo4j.Relationship `json:"relationships"`
	Nodes         []neo4j.Node         `json:"nodes"`
}

type ExportMetadata struct {
	ExportedAt        string `json:"exportedAt"`
	NodeCount         int    `json:"nodeCount"`
	RelationshipCount int    `json:"relationshipCount"`
	Format            string `json:"format"`
}

type GraphResponse struct {
	Graph    Graph          `json:"data"`
	Metadata ExportMetadata `json:"metadata"`
}
