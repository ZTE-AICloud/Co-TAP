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

// POST - HTTP /knowledgegraph/relationships/bulk
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

// POST - HTTP /knowledgegraph/relationships
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

// PUT - HTTP /knowledgegraph/relationships/{elementId}
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

	relationships, err := graphstorage.UpdateRelationship(elementId, props)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(relationships) != 1 {
		err = fmt.Errorf("relationship[%s] does not exist, failed to update", elementId)
		logger.GetLogger().Errorf(err.Error())
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseCodeBody(w, http.StatusCreated, relationships[0])
}

// Delete - HTTP /knowledgegraph/relationships/{elementId}
func (c *RelationshipController) Delete(w http.ResponseWriter, r *http.Request) {
	elementId := mux.Vars(r)["elementId"]

	_, err := graphstorage.DeleteRelationship(elementId)

	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseCodeBody(w, http.StatusNoContent, "")
}

// GET - HTTP /knowledgegraph/relationships/{elementId}
func (c *RelationshipController) GetRelationship(w http.ResponseWriter, r *http.Request) {
	elementId := mux.Vars(r)["elementId"]

	relations, err := graphstorage.QueryRelationship(elementId)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
	} else {
		ResponseCodeBody(w, http.StatusOK, relations)
	}
}

// GET - HTTP /knowledgegraph/relationships?page=1&limit=100&type=depend
func (c *RelationshipController) GetNodes(w http.ResponseWriter, r *http.Request) {

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
