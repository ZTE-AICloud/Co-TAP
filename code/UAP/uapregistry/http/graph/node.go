package graph

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"uapregistry/logger"
	graphstorage "uapregistry/storage/graph"

	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type NodeController struct {
}

// POST - HTTP /knowledgegraph/nodes/bulk
func (c *NodeController) CreateBulk(w http.ResponseWriter, r *http.Request) {
	nodes, err := c.loadNodes(w, r)

	if err != nil {
		return
	}

	newNode, err := graphstorage.CreateNodes(nodes)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
	} else {
		ResponseCodeBody(w, http.StatusCreated, newNode)
	}
}

// POST - HTTP /knowledgegraph/nodes
func (c *NodeController) Create(w http.ResponseWriter, r *http.Request) {
	node, err := c.loadNode(w, r)

	if err != nil {
		return
	}

	newNode, err := graphstorage.CreateNode(node)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
	} else {
		ResponseCodeBody(w, http.StatusCreated, newNode)
	}
}

// PUT - HTTP /knowledgegraph/nodes/{elementId}
func (c *NodeController) Put(w http.ResponseWriter, r *http.Request) {
	elementId := mux.Vars(r)["elementId"]

	newNode, err := c.loadNode(w, r)
	if err != nil {
		return
	}

	nodes, err := graphstorage.UpdateNode(elementId, newNode)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
	} else {
		ResponseCodeBody(w, http.StatusCreated, nodes)
	}
}

// Delete - HTTP /knowledgegraph/nodes/{elementId}
func (c *NodeController) Delete(w http.ResponseWriter, r *http.Request) {
	elementId := mux.Vars(r)["elementId"]

	_, err := graphstorage.DeleteNode(elementId)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
	} else {
		ResponseCodeBody(w, http.StatusNoContent, nil)
	}
}

func (c *NodeController) loadNode(w http.ResponseWriter, r *http.Request) (node neo4j.Node, err error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.GetLogger().Errorf("failed to ReadAll http body:%v", err)
		ResponseCodeBody(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.Unmarshal(body, &node); err != nil {
		logger.GetLogger().Errorf("failed to Unmarshal body, err: %v, body:%s", err, string(body))
		ResponseCodeBody(w, http.StatusUnprocessableEntity, "failed to Unmarshal body, err: "+err.Error()+", body:"+string(body))
		return
	}

	return
}
func (c *NodeController) loadNodes(w http.ResponseWriter, r *http.Request) (nodes []neo4j.Node, err error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.GetLogger().Errorf("failed to ReadAll http body:%v", err)
		ResponseCodeBody(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.Unmarshal(body, &nodes); err != nil {
		logger.GetLogger().Errorf("failed to Unmarshal body, err: %v, body:%s", err, string(body))
		ResponseCodeBody(w, http.StatusUnprocessableEntity, "failed to Unmarshal body, err: "+err.Error()+", body:"+string(body))
		return
	}

	return
}

// GET - HTTP /knowledgegraph/nodes/{agentCardName}/namespace/{ns1}?cluster=c1
func (c *NodeController) GetNodesByName(w http.ResponseWriter, r *http.Request) {
	agentCardName := mux.Vars(r)["agentCardName"]
	namespace := mux.Vars(r)["namespace"]
	cluster := r.URL.Query().Get("cluster")

	nodes, err := graphstorage.QueryNode(agentCardName, namespace, cluster)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
	} else {
		ResponseCodeBody(w, http.StatusOK, nodes)
	}
}

// GET - HTTP /knowledgegraph/nodes?page=1&limit=100&label=Person
func (c *NodeController) GetNodes(w http.ResponseWriter, r *http.Request) {
	label := r.URL.Query().Get("label")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, limit, err := parsePaginationParams(pageStr, limitStr)
	if err != nil {
		errInfo := fmt.Sprintf("page:%s, limit:%s, err:%v", pageStr, limitStr, err.Error())
		logger.GetLogger().Errorf(errInfo)
		ResponseCodeBody(w, http.StatusUnprocessableEntity, errInfo)
		return
	}

	nodes, err := graphstorage.QueryNodes(label, page, limit)

	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
	} else {
		ResponseCodeBody(w, http.StatusOK, nodes)
	}
}

// GET - HTTP /knowledgegraph/nodes/{elementId}
func (c *NodeController) GetNodeByID(w http.ResponseWriter, r *http.Request) {
	elementId := mux.Vars(r)["elementId"]

	node, err := graphstorage.QueryNodeByID(elementId)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
	} else {
		ResponseCodeBody(w, http.StatusOK, node)
	}
}

// GET - HTTP /knowledgegraph/nodes/{elementId}/relations
func (c *NodeController) GetRelatedNodes(w http.ResponseWriter, r *http.Request) {
	elementId := mux.Vars(r)["elementId"]
	nodes, err := graphstorage.QueryRelatedNodes(elementId)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
	} else {
		ResponseCodeBody(w, http.StatusOK, nodes)
	}
}
