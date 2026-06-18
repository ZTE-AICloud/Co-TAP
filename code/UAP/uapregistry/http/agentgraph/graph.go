package agentgraph

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"uapregistry/logger"
	"uapregistry/storage/agentgraphstorage"
	"uapregistry/types/agentgraphmodels"

	"github.com/gorilla/mux"
)

type GraphController struct {
}

// GET - HTTP /agentgraph/graph?page=0&limit=1000
func (c *GraphController) Export(w http.ResponseWriter, r *http.Request) {
	pageStr := mux.Vars(r)["page"]
	limitStr := mux.Vars(r)["limit"]
	page, limit, err := parsePaginationParams(pageStr, limitStr)
	if err != nil {
		ResponseCodeBody(w, http.StatusUnprocessableEntity, err.Error())
	}

	graph, err := agentgraphstorage.ExportGraph(page, limit)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
	} else {

		response := agentgraphmodels.GraphResponse{
			Graph: graph,
			Metadata: agentgraphmodels.ExportMetadata{
				ExportedAt:        time.Now().Format(time.RFC3339),
				NodeCount:         len(graph.Nodes),
				RelationshipCount: len(graph.Relationships),
				Format:            "json",
			},
		}
		ResponseCodeBody(w, http.StatusOK, response)
	}
}

// GET - HTTP /agentgraph/graph
func (c *GraphController) Import(w http.ResponseWriter, r *http.Request) {
	graph, err := c.loadGraphData(w, r)

	if err != nil {
		return
	}

	nodes, err := agentgraphstorage.ImportGraph(graph)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
	} else {
		ResponseCodeBody(w, http.StatusOK, nodes)
	}
}

func (c *GraphController) loadGraphData(w http.ResponseWriter, r *http.Request) (graph agentgraphmodels.Graph, err error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.GetLogger().Errorf("failed to ReadAll http body:%v", err)
		ResponseCodeBody(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.Unmarshal(body, &graph); err != nil {
		logger.GetLogger().Errorf("failed to Unmarshal request body, err: %v, body:%s", err, string(body))
		ResponseCodeBody(w, http.StatusUnprocessableEntity, "failed to Unmarshal body, err: "+err.Error()+", body:"+string(body))
		return
	}

	return
}
