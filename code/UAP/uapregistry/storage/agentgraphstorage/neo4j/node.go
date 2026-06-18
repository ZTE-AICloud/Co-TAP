package neo4j

import (
	"context"
	"fmt"
	"strings"
	"time"
	"uapregistry/types/agentgraphmodels"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// DeleteNode 删除节点
func (c *Client) executeQueryRaw(cypher string, params map[string]any) (res *neo4j.EagerResult, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return neo4j.ExecuteQuery(ctx, c.Driver, cypher, params, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(c.Config.Database))
}

// executeQueryNode 查询节点信息
// cypher - Cypher 查询语句
// params - 参数
func (c *Client) executeQueryNode(cypher string, params map[string]any) (nodes []neo4j.Node, err error) {
	result, err := c.executeQueryRaw(cypher, params)
	if err != nil {
		return
	}

	for _, record := range result.Records {
		nodeRaw, found := record.Get("n")
		if !found {
			continue
		}

		node := nodeRaw.(neo4j.Node)
		nodes = append(nodes, node)
	}
	return
}

func (c *Client) generateCreateNodeParams(node neo4j.Node) (cypher string, params map[string]any) {
	params = make(map[string]any)

	var cypherBuild strings.Builder

	// 处理属性列表
	if len(node.Props) != 0 {
		params["props"] = node.Props
		cypherBuild.WriteString("CREATE (n $props)")
	} else {
		cypherBuild.WriteString("CREATE (n)")
	}

	// 处理标签列表
	if len(node.Labels) > 0 {
		cypherBuild.WriteString(" SET n:" + strings.Join(node.Labels, ":"))
	}

	// 添加返回值
	cypherBuild.WriteString(" RETURN n")
	cypher = cypherBuild.String()

	return
}

// CreateNode 创建节点
// label - 节点标签
// properties - 节点属性
func (c *Client) CreateNode(node neo4j.Node) (newNode neo4j.Node, err error) {
	cypher, params := c.generateCreateNodeParams(node)

	nodes, err := c.executeQueryNode(cypher, params)
	if err != nil {
		return
	}

	if len(nodes) == 0 {
		err = fmt.Errorf("failed to create node")
		return
	}

	return nodes[0], nil
}

// CreateNodes 创建节点
// label - 节点标签
// properties - 节点属性
func (c *Client) CreateNodes(nodes []neo4j.Node) (newNodes []neo4j.Node, err error) {

	session := c.Driver.NewSession(context.Background(), neo4j.SessionConfig{DatabaseName: c.Config.Database})
	defer session.Close(context.Background())

	_, err = session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		for _, node := range nodes {
			cypher, params := c.generateCreateNodeParams(node)
			nodeRaw, err := c.excuteInTransaction(tx, cypher, params, "n")

			if err != nil {
				return nil, err
			}
			newNodes = append(newNodes, nodeRaw.(neo4j.Node))
		}
		return nil, nil
	})

	return
}

func (c *Client) UpdateNode(elementId string, node neo4j.Node) (newNode []neo4j.Node, err error) {
	params := make(map[string]any)
	params["elementId"] = elementId
	params["props"] = node.Props

	var setNewLabel string
	// 处理标签
	if len(node.Labels) != 0 {
		setNewLabel = `SET n:` + strings.Join(node.Labels, ":")
	}

	cypher := fmt.Sprintf(`
		MATCH (n)
		WHERE elementId(n) = $elementId

		// 设置新标签
		%s
		SET n = $props

		RETURN n

	`, setNewLabel)

	return c.executeQueryNode(cypher, params)
}

// DeleteNode 删除节点
func (c *Client) DeleteNode(elementId string) (nodes []neo4j.Node, err error) {
	cypher := `
	MATCH (n)
	WHERE elementId(n) = $elementId
	// 打包完整的节点快照数据
	WITH n, {
		elementId: elementId(n),
		labels: labels(n),
		props: properties(n)
	} AS nodeSnapshot
	DETACH DELETE n 
	RETURN nodeSnapshot`

	params := map[string]any{
		"elementId": elementId,
	}

	return c.executeQueryNode(cypher, params)
}

// QueryNodes 查询所有节点 page 从 0 开始
func (c *Client) QueryNodes(label string, page, limit int) (nodes []neo4j.Node, err error) {
	param := map[string]any{}
	var cypher string

	if label == "" {
		cypher = "MATCH (n) RETURN n"
	} else {
		cypher = fmt.Sprintf("MATCH (n:%s) RETURN n", label)
	}

	if page >= 0 {
		skip := page * limit
		param["skip"] = skip
		param["limit"] = limit
		cypher = cypher + " ORDER BY elementId(n) SKIP $skip LIMIT $limit"
	}

	return c.executeQueryNode(cypher, param)
}

// QueryNodeByID 根据 ID 查询节点
func (c *Client) QueryNodeByID(id string) (node neo4j.Node, err error) {
	cypher := "MATCH (n) WHERE elementId(n) = $id RETURN n"
	params := map[string]any{"id": id}

	nodes, err := c.executeQueryNode(cypher, params)
	if err != nil {
		return
	}

	if len(nodes) == 0 {
		err = fmt.Errorf("node with id %s not found", id)
		return
	}

	return nodes[0], nil
}

// QueryNode
func (c *Client) QueryNode(agentName, namespace, cluster string) (nodes []neo4j.Node, err error) {
	cypher := "MATCH (n {agentName: $agentName, namespace: $namespace, cluster: $cluster } ) RETURN n"
	params := map[string]any{
		"agentName": agentName,
		"namespace": namespace,
		"cluster":   cluster,
	}

	return c.executeQueryNode(cypher, params)
}

func (c *Client) QueryRelationshipsByNode(elementId string) (relatedNodes []agentgraphmodels.RelatedNode, err error) {
	cypher := `
		MATCH (n)-[r]-(neighbor)
		WHERE elementId(n) = $elementId
		RETURN r, neighbor
	`

	params := map[string]any{
		"elementId": elementId,
	}

	result, err := c.executeQueryRaw(cypher, params)
	if err != nil {
		return
	}

	for _, record := range result.Records {
		nodeRaw, found := record.Get("neighbor")
		if !found {
			continue
		}

		node := nodeRaw.(neo4j.Node)

		relationshipRaw, found := record.Get("r")
		if !found {
			continue
		}

		relationship := relationshipRaw.(neo4j.Relationship)

		relatedNodes = append(relatedNodes, agentgraphmodels.RelatedNode{
			Node:         node,
			Relationship: relationship,
		})
	}

	return
}
