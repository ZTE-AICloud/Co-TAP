package neo4j

import (
	"context"
	"fmt"
	"time"
	"uapregistry/logger"
	"uapregistry/types"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Config Neo4j 连接配置
type Config struct {
	URI      string
	Username string
	Password string
	Database string
}

// Client Neo4j 客户端
type Client struct {
	Driver neo4j.DriverWithContext
	Config types.DatabaseConfig
}

// NewClient 创建新的 Neo4j 客户端
func NewClient(config types.DatabaseConfig) (*Client, error) {
	authToken := neo4j.BasicAuth(config.Username, config.Password, "")
	driver, err := neo4j.NewDriverWithContext(config.URI, authToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}

	client := &Client{
		Driver: driver,
		Config: config,
	}

	// 验证连接
	if err := client.VerifyConnectivity(); err != nil {
		return nil, fmt.Errorf("failed to verify connectivity: %w", err)
	}

	return client, nil
}

// VerifyConnectivity 验证与 Neo4j 的连接
func (c *Client) VerifyConnectivity() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.Driver.VerifyConnectivity(ctx)
}

// Close 关闭客户端连接
func (c *Client) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.Driver.Close(ctx)
}

// ImportGraph 创建节点
// label - 节点标签
// properties - 节点属性
func (c *Client) ImportGraph(graph types.Graph) (newGrahp types.Graph, err error) {

	session := c.Driver.NewSession(context.Background(), neo4j.SessionConfig{DatabaseName: c.Config.Database})
	defer session.Close(context.Background())

	_, err = session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		newNodes := make([]neo4j.Node, len(graph.Nodes))
		oldNewNodeIdMap := make(map[string]string, len(graph.Nodes))
		for i, node := range graph.Nodes {
			cypher, params := c.generateCreateNodeParams(node)

			nodeRaw, err := c.excuteInTransaction(tx, cypher, params, "n")
			if err != nil {
				return nil, err
			}
			newNode := nodeRaw.(neo4j.Node)
			oldNewNodeIdMap[node.ElementId] = newNode.ElementId
			newNodes[i] = newNode
		}

		newRelationships := make([]neo4j.Relationship, len(graph.Relationships))
		for i, relationship := range graph.Relationships {
			// 使用新的节点id 替换原来的节点id
			relationship.StartElementId = oldNewNodeIdMap[relationship.StartElementId]
			relationship.EndElementId = oldNewNodeIdMap[relationship.EndElementId]

			cypher, params := c.generateCreateRelationParams(relationship)
			newRelationship, err := c.excuteInTransaction(tx, cypher, params, "r")
			if err != nil {
				return nil, err
			}
			newRelationships[i] = newRelationship.(neo4j.Relationship)
		}

		newGrahp.Nodes = newNodes
		newGrahp.Relationships = newRelationships

		return nil, nil
	})

	return
}
func (c *Client) ExportGraph(page, limit int) (graph types.Graph, err error) {
	session := c.Driver.NewSession(context.Background(), neo4j.SessionConfig{DatabaseName: c.Config.Database})
	defer session.Close(context.Background())

	_, err = session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		nodes, err := c.listNodes(tx, page, limit)
		if err != nil {
			return nil, err
		}
		graph.Nodes = nodes

		relationships, err := c.listRelationships(tx, page, limit)
		if err != nil {
			return nil, err
		}
		graph.Relationships = relationships

		return nil, nil
	})

	return
}

func (c *Client) listNodes(tx neo4j.ManagedTransaction, page, limit int) (newNodes []neo4j.Node, err error) {
	var cypher = "MATCH (n) RETURN elementId(n) AS elementId, labels(n) AS labels, properties(n) AS props"
	params := make(map[string]any, 0)
	if page >= 0 {
		cypher = cypher + " SKIP $skip LIMIT $limit"
		params["skip"] = page * limit
		params["limit"] = limit
	}
	result, err := tx.Run(context.Background(), cypher, params)
	if err != nil {
		logger.GetLogger().Errorf("listNodes query error, %s", err.Error())
		return nil, err
	}

	// 此处不能使用 result.Collect 直接获取全部数据, Collect 方法会一次性打包所有数据，造成内存压力，Next 方法内存稳定，是一段一段处理。
	for result.Next(context.Background()) {
		record := result.Record()

		// 1. 安全解析 labels
		var labelStrings []string
		if labels, ok := record.AsMap()["labels"].([]any); ok {
			for _, l := range labels {
				if str, ok := l.(string); ok {
					labelStrings = append(labelStrings, str)
				}
			}
		}

		newNodes = append(newNodes, neo4j.Node{
			ElementId: record.AsMap()["elementId"].(string),
			Labels:    labelStrings,
			Props:     record.AsMap()["props"].(map[string]any),
		})
	}

	// 有可能由于网络原因退出result.Next 循环， 此处捕获这种可能
	err = result.Err()
	if err != nil {
		logger.GetLogger().Errorf("listNodes get result error, %s", err.Error())
	}

	return
}

func (c *Client) listRelationships(tx neo4j.ManagedTransaction, page, limit int) (relationships []neo4j.Relationship, err error) {
	var cypher = "MATCH ()-[r]->() RETURN elementId(r) AS elementId, elementId(startNode(r)) AS from_id, elementId(endNode(r)) AS to_id, type(r) AS rel_type, properties(r) AS props"
	params := make(map[string]any, 0)
	if page >= 0 {
		cypher = cypher + " SKIP $skip LIMIT $limit"
		params["skip"] = page * limit
		params["limit"] = limit
	}
	result, err := tx.Run(context.Background(), cypher, params)
	if err != nil {
		logger.GetLogger().Errorf("listRelationships query error, %s", err.Error())
		return nil, err
	}

	// 此处不能使用 result.Collect 直接获取全部数据, Collect 方法会一次性打包所有数据，造成内存压力，Next 方法内存稳定，是一段一段处理。
	for result.Next(context.Background()) {
		record := result.Record()

		relationships = append(relationships, neo4j.Relationship{
			ElementId:      record.AsMap()["elementId"].(string),
			StartElementId: record.AsMap()["from_id"].(string),
			EndElementId:   record.AsMap()["to_id"].(string),
			Type:           record.AsMap()["rel_type"].(string),
			Props:          record.AsMap()["props"].(map[string]any),
		})
	}

	// 有可能由于网络原因退出result.Next 循环， 此处捕获这种可能
	err = result.Err()
	if err != nil {
		logger.GetLogger().Errorf("listRelationships get result error, %s", err.Error())
	}

	return
}

func (c *Client) excuteInTransaction(tx neo4j.ManagedTransaction, cypher string, params map[string]any, resultKey string) (res any, err error) {
	result, err := tx.Run(context.Background(), cypher, params)
	if err != nil {
		logger.GetLogger().Errorf("excuteInTransaction err, %s", err.Error())
		return
	}

	record, err := result.Single(context.Background())
	if err != nil { // none or more than one
		logger.GetLogger().Errorf("excuteInTransaction result is not single, %s", err.Error())
		return
	}

	res, found := record.Get(resultKey)
	if !found {
		err = fmt.Errorf("can not find relationship in return")
		logger.GetLogger().Errorf("excuteInTransaction result does not contain key: %s", resultKey)
		return
	}

	return
}
