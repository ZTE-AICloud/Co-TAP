package neo4j

import (
	"context"
	"fmt"
	"time"
	"uapregistry/logger"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (c *Client) CreateRelationship(relationship neo4j.Relationship) (newRelationship neo4j.Relationship, err error) {

	// 拼装 Cypher
	// 使用 MERGE (from)-[r:%s]->(to) 确保关系的唯一性
	cypher, params := c.generateCreateRelationParams(relationship)

	relations, err := c.executeQueryRelationship(cypher, params)

	if err != nil {
		return
	}

	if len(relations) != 0 {
		newRelationship = relations[0]
	}

	return
}

func (c *Client) CreateRelationships(relationships []neo4j.Relationship) (newRelationships []neo4j.Relationship, err error) {

	session := c.Driver.NewSession(context.Background(), neo4j.SessionConfig{DatabaseName: c.Config.Database})
	defer session.Close(context.Background())

	_, err = session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		for _, relationship := range relationships {
			cypher, params := c.generateCreateRelationParams(relationship)
			newRelationship, err := c.excuteInTransaction(tx, cypher, params, "r")
			if err != nil {
				return nil, err
			}

			newRelationships = append(newRelationships, newRelationship.(neo4j.Relationship))
		}
		return nil, nil
	})

	return
}

func (c *Client) generateCreateRelationParams(relationship neo4j.Relationship) (cypher string, params map[string]any) {
	// 拼装 Cypher
	// 使用 MERGE (from)-[r:%s]->(to) 确保关系的唯一性
	cypher = fmt.Sprintf(`
				MATCH (from) WHERE elementId(from) = $fromId
				MATCH (to) WHERE elementId(to) = $toId
				
				// 建立唯一关系
				MERGE (from)-[r:%s]->(to)
				
				// 如果关系是新创建的，设置初始属性
				ON CREATE SET r.createdAt = $timestamp, r.status = "Initiated"
				
				// 如果关系早就存在，则更新最后活跃时间
				ON MATCH SET r.updatedAt = $timestamp
		
				SET r = $props
				
				RETURN r
			`, relationship.Type)

	params = map[string]any{
		"fromId":    relationship.StartElementId,
		"toId":      relationship.EndElementId,
		"timestamp": time.Now().Format(time.RFC3339),
		"props":     relationship.Props,
	}

	return
}

// UpdateRelationship 只能更行Properties
func (c *Client) UpdateRelationship(elementId string, props map[string]any) (relation neo4j.Relationship, err error) {
	// 动态组装原子覆盖 Cypher 语句
	//
	cypherQuery := `
		MATCH (from)-[r]->(to)
		WHERE elementId(r) = $elementId
		
		SET r = $props
		
		RETURN r
	`

	params := map[string]any{
		"elementId": elementId,
		"props":     props,
	}

	relations, err := c.executeQueryRelationship(cypherQuery, params)
	if err != nil {
		return
	}

	if len(relations) != 1 {
		err = fmt.Errorf("updated relationship len is %d", len(relations))
		logger.GetLogger().Error(err.Error())
		return
	}

	relation = relations[0]

	return
}

func (c *Client) DeleteRelationship(elementId string) (relationships []neo4j.Relationship, err error) {
	cypherQuery := `
		MATCH (from)-[r]->(to)
		WHERE elementId(r) = $elementId
		DELETE r
	`

	params := map[string]any{
		"elementId": elementId,
	}

	return c.executeQueryRelationship(cypherQuery, params)
}

func (c *Client) executeQueryRelationship(cypher string, params map[string]any) (relations []neo4j.Relationship, err error) {
	result, err := c.executeQueryRaw(cypher, params)
	if err != nil {
		return
	}

	for _, record := range result.Records {
		rel, found := record.Get("r")
		if !found {
			continue
		}

		relations = append(relations, rel.(neo4j.Relationship))
	}

	return
}

// QueryRelationships 查询所有节点 page 从 1 开始
func (c *Client) QueryRelationships(typ string, page, limit int) (nodes []neo4j.Relationship, err error) {
	param := map[string]any{}
	var cypher string

	if typ == "" {
		cypher = "MATCH ()-[r]->()  RETURN r"
	} else {
		cypher = fmt.Sprintf("MATCH ()-[r:%s]->() RETURN r", typ)
	}

	if page >= 1 {
		skip := (page - 1) * limit
		param["skip"] = skip
		param["limit"] = limit
		cypher = cypher + " ORDER BY elementId(r) SKIP $skip LIMIT $limit"
	}

	return c.executeQueryRelationship(cypher, param)
}

// QueryRelationship
func (c *Client) QueryRelationship(elementId string) (node neo4j.Relationship, exist bool, err error) {
	param := map[string]any{}
	param["elementId"] = elementId

	cypher := "MATCH ()-[r]->() WHERE elementId(r) = $elementId RETURN r"

	relations, err := c.executeQueryRelationship(cypher, param)
	if err != nil {
		return
	}

	if len(relations) == 0 {
		return
	}

	return relations[0], true, err
}
