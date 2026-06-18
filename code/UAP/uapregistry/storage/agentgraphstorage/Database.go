package agentgraphstorage

import (
	"uapregistry/logger"
	neo4jclient "uapregistry/storage/agentgraphstorage/neo4j"
	"uapregistry/types/agentgraphmodels"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Database interface {
	ImportGraph(graph agentgraphmodels.Graph) (newGraph agentgraphmodels.Graph, err error)
	ExportGraph(page, limit int) (graph agentgraphmodels.Graph, err error)

	CreateNode(node neo4j.Node) (newNode neo4j.Node, err error)
	CreateNodes(nodes []neo4j.Node) (newNodes []neo4j.Node, err error)
	DeleteNode(elementId string) (nodes []neo4j.Node, err error)
	UpdateNode(elementId string, newNode neo4j.Node) (node []neo4j.Node, err error)

	CreateRelationship(relationship neo4j.Relationship) (newRelationship neo4j.Relationship, err error)
	CreateRelationships(relationships []neo4j.Relationship) (newRelationships []neo4j.Relationship, err error)
	UpdateRelationship(elementId string, propMap map[string]any) (newRelationship []neo4j.Relationship, err error)
	DeleteRelationship(elementId string) (newRelationship []neo4j.Relationship, err error)

	QueryNodeByID(id string) (node neo4j.Node, err error)
	QueryNode(name, namespace, cluster string) (nodes []neo4j.Node, err error)
	QueryNodes(label string, page, limit int) (nodes []neo4j.Node, err error)
	QueryRelationshipsByNode(nodeId string) (relations []agentgraphmodels.RelatedNode, err error)

	QueryRelationship(elementId string) (relatinships []neo4j.Relationship, err error)
	QueryRelationships(typ string, page, limit int) (relatinships []neo4j.Relationship, err error)
}

var client Database
var log = logger.GetLogger()

func InitDatabase(config agentgraphmodels.DatabaseConfig) error {
	if config.URI == "" || config.Username == "" || config.Password == "" || config.Database == "" {
		log.Info("agent graph database config is empty, skip init agent graph database")
		return nil
	}

	var err error
	client, err = neo4jclient.NewClient(config)
	if err != nil {
		log.Errorf("Failed to create Neo4j client: %v", err)
	} else {
		log.Info("succeed to init storage")
	}

	return err
}

func ImportGraph(graph agentgraphmodels.Graph) (newGraph agentgraphmodels.Graph, err error) {
	return client.ImportGraph(graph)
}

func ExportGraph(page, limit int) (newGraph agentgraphmodels.Graph, err error) {
	return client.ExportGraph(page, limit)
}

func CreateNode(node neo4j.Node) (newNode neo4j.Node, err error) {
	return client.CreateNode(node)
}
func CreateNodes(nodes []neo4j.Node) (newNodes []neo4j.Node, err error) {
	return client.CreateNodes(nodes)
}

func DeleteNode(elementId string) (nodes []neo4j.Node, err error) {
	return client.DeleteNode(elementId)
}

func UpdateNode(elementId string, newNode neo4j.Node) (node []neo4j.Node, err error) {
	return client.UpdateNode(elementId, newNode)
}

func CreateRelationship(relationship neo4j.Relationship) (newRelationship neo4j.Relationship, err error) {
	return client.CreateRelationship(relationship)
}
func CreateRelationships(relationships []neo4j.Relationship) (newRelationships []neo4j.Relationship, err error) {
	return client.CreateRelationships(relationships)
}

func DeleteRelationship(elementId string) (nodes []neo4j.Relationship, err error) {
	return client.DeleteRelationship(elementId)
}

func UpdateRelationship(elementId string, props map[string]any) (updatedRelation []neo4j.Relationship, err error) {
	return client.UpdateRelationship(elementId, props)
}

func QueryNode(agentCardName, namespace, cluster string) (node []neo4j.Node, err error) {
	return client.QueryNode(agentCardName, namespace, getDefaultCluster(cluster))
}

func QueryNodeByID(elementId string) (node neo4j.Node, err error) {
	return client.QueryNodeByID(elementId)
}

func QueryNodes(label string, page, limit int) (node []neo4j.Node, err error) {
	return client.QueryNodes(label, page, limit)
}

func QueryRelatedNodes(elementId string) (relatedNodes []agentgraphmodels.RelatedNode, err error) {
	return client.QueryRelationshipsByNode(elementId)
}

func QueryRelationships(typ string, page, limit int) (node []neo4j.Relationship, err error) {
	return client.QueryRelationships(typ, page, limit)
}

func QueryRelationship(elementId string) (node []neo4j.Relationship, err error) {
	return client.QueryRelationship(elementId)
}

func getDefaultCluster(cluster string) string {
	if cluster == "" {
		cluster = "local"
	}
	return cluster
}
