package graph

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"uapregistry/logger"
	graphstorage "uapregistry/storage/graph"

	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type NodeController struct {
}

// CreateBulk 批量创建节点
// @Summary      批量创建节点
// @Description  导入节点数组批量创建节点
// @Tags         node
// @Accept       json
// @Produce      json
// @Param        nodes   body      []neo4j.Node  true  "节点列表"
// @Success      201  {array}  neo4j.Node
// @Failure      422  {string}  string
// @Failure      500  {string}  string
// @Router       /knowledgegraph/nodes/bulk [POST]
func (c *NodeController) CreateBulk(w http.ResponseWriter, r *http.Request) {
	nodes, err := c.loadNodes(w, r)

	if err != nil {
		return
	}

	newNodes, err := graphstorage.CreateNodes(nodes)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
	} else {
		ResponseCodeBody(w, http.StatusCreated, newNodes)
	}
}

// Create 创建单个节点
// @Summary      创建单个节点
// @Description  创建单个节点
// @Tags         node
// @Accept       json
// @Produce      json
// @Param        node   body      neo4j.Node  true  "节点信息"
// @Success      201  {object}  neo4j.Node
// @Failure      422  {string}  string
// @Failure      500  {string}  string
// @Router       /knowledgegraph/nodes [POST]
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

// Put 更新单个节点
// @Summary      更新单个节点
// @Description  更新单个节点
// @Tags         node
// @Accept       json
// @Produce      json
// @Param        elementId   path      string  true  "旧节点elementId"
// @Param        node   body      neo4j.Node  true  "节点信息"
// @Success      201  {object}  neo4j.Node
// @Failure      422  {string}  string
// @Failure      500  {string}  string
// @Router       /knowledgegraph/nodes/{elementId} [PUT]
func (c *NodeController) Put(w http.ResponseWriter, r *http.Request) {
	elementId := mux.Vars(r)["elementId"]

	newNode, err := c.loadNode(w, r)
	if err != nil {
		return
	}

	node, err := graphstorage.UpdateNode(elementId, newNode)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
	} else {
		ResponseCodeBody(w, http.StatusCreated, node)
	}
}

// Delete 删除单个节点
// @Summary      删除单个节点
// @Description  删除单个节点
// @Tags         node
// @Accept       json
// @Produce      json
// @Param        elementId   path      string  true  "旧节点elementId"
// @Param        force   query      bool  false  "是否强制删除，默认false，当存在相关的关系时不删除"
// @Success      204  "删除成功"
// @Failure      409  {string}  string "无法删除有关联的节点（force=false）"
// @Failure      422  {string}  string "请求参数错误"
// @Failure      500  {string}  string
// @Router       /knowledgegraph/nodes/{elementId} [Delete]
func (c *NodeController) Delete(w http.ResponseWriter, r *http.Request) {
	elementId := mux.Vars(r)["elementId"]
	forceStr := r.URL.Query().Get("force")

	if forceStr == "" {
		forceStr = "false"
	}
	force, err := strconv.ParseBool(forceStr)
	if err != nil {
		ResponseCodeBody(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	canDelete, err := graphstorage.DeleteNode(elementId, force)
	if err != nil {
		if canDelete {
			ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
		} else {
			ResponseCodeBody(w, http.StatusConflict, err.Error())
		}
	} else {
		ResponseCodeBody(w, http.StatusNoContent, "")
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

// GetNodes 查询节点列表
// @Summary      查询节点列表
// @Description  查询节点列表
// @Tags         node
// @Accept       json
// @Produce      json
// @Param        page   query      int      false  "分页参数,当前页数，从1开始,默认不分页"
// @Param        limit  query      int      false  "分页参数，每页最大条数"
// @Param        label  query      string   true   "节点标签, 默认为空，不区分标签"
// @Success      200  {array}   neo4j.Node "查询成功"
// @Failure      422  {string}  string
// @Failure      500  {string}  string
// @Router       /knowledgegraph/nodes [GET]
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

// GetNodeByID 查询指定节点
// @Summary      查询指定节点
// @Description  查询指定节点
// @Tags         node
// @Accept       json
// @Produce      json
// @Param        elementId   path      string      true  "节点elementId"
// @Success      200  {object}  neo4j.Node "查询成功"
// @Failure      404  {string}  string "资源不存在"
// @Failure      500  {string}  string
// @Router       /knowledgegraph/nodes [GET]
// GET - HTTP /knowledgegraph/nodes/{elementId}
func (c *NodeController) GetNodeByID(w http.ResponseWriter, r *http.Request) {
	elementId := mux.Vars(r)["elementId"]

	node, exist, err := graphstorage.QueryNodeByID(elementId)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
		return
	}

	if exist {
		ResponseCodeBody(w, http.StatusOK, node)
	} else {
		ResponseCodeBody(w, http.StatusNotFound, "resource not found")
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
