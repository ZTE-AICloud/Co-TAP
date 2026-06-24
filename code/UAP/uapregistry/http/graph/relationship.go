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

type RelationshipController struct {
}

// CreateBulk 批量创建关系
// @Summary      批量创建关系
// @Description  导入节点数组批量创建关系
// @Tags         relationship
// @Accept       json
// @Produce      json
// @Param        relationships   body      []neo4j.Relationship  true  "关系列表"
// @Success      201  {array}  neo4j.Relationship
// @Failure      422  {string}  string
// @Failure      500  {string}  string
// @Router       /knowledgegraph/relationships/bulk [POST]
func (c *RelationshipController) CreateBulk(w http.ResponseWriter, r *http.Request) {
	var relationships []neo4j.Relationship
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.GetLogger().Errorf("Failed to ReadAll http body:%v", err)
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
		return
	}

	// invalid params
	if err = json.Unmarshal(body, &relationships); err != nil {
		err = fmt.Errorf("Failed to Unmarshal body, err: %v, body:%s", err, string(body))
		logger.GetLogger().Errorf(err.Error())
		ResponseCodeBody(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	relationships, err = graphstorage.CreateRelationships(relationships)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
	} else {
		ResponseCodeBody(w, http.StatusCreated, relationships)
	}
}

// Create 创建关系
// @Summary      创建关系
// @Description  创建关系
// @Tags         relationship
// @Accept       json
// @Produce      json
// @Param        relationship   body      neo4j.Relationship  true  "关系"
// @Success      201  {object}  neo4j.Relationship
// @Failure      422  {string}  string
// @Failure      500  {string}  string
// @Router       /knowledgegraph/relationships [POST]
func (c *RelationshipController) Create(w http.ResponseWriter, r *http.Request) {
	var relationship neo4j.Relationship
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.GetLogger().Errorf("Failed to ReadAll http body:%v", err)
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
		return
	}

	// invalid params
	if err = json.Unmarshal(body, &relationship); err != nil {
		err = fmt.Errorf("Failed to Unmarshal body, err: %v, body:%s", err, string(body))
		logger.GetLogger().Errorf(err.Error())
		ResponseCodeBody(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	relationship, err = graphstorage.CreateRelationship(relationship)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
	} else {
		ResponseCodeBody(w, http.StatusCreated, relationship)
	}
}

// Put 更新关系
// @Summary      更新关系
// @Description  更新关系
// @Tags         relationship
// @Accept       json
// @Produce      json
// @Param        elementId   path      string  true  "关系elementId"
// @Param        props   body      map[string]interface{}  true  "关系属性"
// @Success      201  {object}  neo4j.Relationship
// @Failure      422  {string}  string
// @Failure      500  {string}  string
// @Router       /knowledgegraph/relationships/{elementId} [PUT]
func (c *RelationshipController) Put(w http.ResponseWriter, r *http.Request) {
	elementId := mux.Vars(r)["elementId"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.GetLogger().Errorf("Failed to ReadAll http body:%v", err)
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
		return
	}

	var props map[string]any

	if err = json.Unmarshal(body, &props); err != nil {
		errInfo := fmt.Sprintf("failed to Unmarshal body, err: %v, body:%s", err, string(body))
		logger.GetLogger().Errorf(errInfo)
		ResponseCodeBody(w, http.StatusUnprocessableEntity, errInfo)
		return
	}

	relationship, err := graphstorage.UpdateRelationship(elementId, props)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseCodeBody(w, http.StatusCreated, relationship)
}

// Delete 删除关系
// @Summary      删除关系
// @Description  删除关系
// @Tags         relationship
// @Accept       json
// @Produce      json
// @Param        elementId   path      string  true  "关系elementId"
// @Success      204  "删除成功"
// @Failure      500  {string}  string
// @Router       /knowledgegraph/relationships/{elementId} [DELETE]
func (c *RelationshipController) Delete(w http.ResponseWriter, r *http.Request) {
	elementId := mux.Vars(r)["elementId"]

	_, err := graphstorage.DeleteRelationship(elementId)

	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseCodeBody(w, http.StatusNoContent, "")
}

// GetRelationship  查询指定关系
// @Summary      查询指定关系
// @Description  查询指定关系
// @Tags         relationship
// @Accept       json
// @Produce      json
// @Param        elementId   path      string  true  "关系elementId"
// @Success      200  {object} neo4j.Relationship  "删除成功"
// @Failure      404  {string}  string "资源不存在"
// @Failure      500  {string}  string
// @Router       /knowledgegraph/relationships/{elementId} [GET]
func (c *RelationshipController) GetRelationship(w http.ResponseWriter, r *http.Request) {
	elementId := mux.Vars(r)["elementId"]

	relationship, exist, err := graphstorage.QueryRelationship(elementId)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
		return
	}

	if exist {
		ResponseCodeBody(w, http.StatusOK, relationship)
	} else {
		ResponseCodeBody(w, http.StatusNotFound, "resource not found")
	}
}

// GetRelationships 查询节点列表
// @Summary      查询节点列表
// @Description  查询节点列表
// @Tags         node
// @Accept       json
// @Produce      json
// @Param        page   query      int      false  "分页参数,当前页数，从1开始,默认不分页"
// @Param        limit  query      int      false  "分页参数，每页最大条数"
// @Param        type  query       string   true   "关系分类 默认为空，不区分分类"
// @Success      200  {array}   neo4j.Relationship "查询成功"
// @Failure      422  {string}  string
// @Failure      500  {string}  string
// @Router       /knowledgegraph/relationships [GET]
func (c *RelationshipController) GetRelationships(w http.ResponseWriter, r *http.Request) {

	typ := r.URL.Query().Get("type")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, limit, err := parsePaginationParams(pageStr, limitStr)
	if err != nil {
		errInfo := fmt.Sprintf("page:%s, limit:%s, err:%v", pageStr, limitStr, err.Error())
		logger.GetLogger().Errorf(errInfo)
		ResponseCodeBody(w, http.StatusUnprocessableEntity, errInfo)
		return
	}

	relations, err := graphstorage.QueryRelationships(typ, page, limit)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
	} else {
		ResponseCodeBody(w, http.StatusOK, relations)
	}
}
