package graph

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"uapregistry/logger"
	graphstorage "uapregistry/storage/graph"
	"uapregistry/types"

	"github.com/gorilla/mux"
)

type GraphController struct {
}

// Export 导出图数据
// @Summary      获取完整的知识图谱信息
// @Description  获取完整知识图谱信息，可以支持分页查询
// @Tags         graph
// @Accept       json
// @Produce      json
// @Param        page   query      int  false  "分页参数,当前页数，从1开始,默认不分页"
// @Param        limit   query      int  false  "分页参数，每页最大条数"
// @Success      200  {object}  types.GraphResponse
// @Failure      422  {string}  string
// @Failure      500  {string}  string
// @Router       /knowledgegraph/graph [get]
func (c *GraphController) Export(w http.ResponseWriter, r *http.Request) {
	pageStr := mux.Vars(r)["page"]
	limitStr := mux.Vars(r)["limit"]
	page, limit, err := parsePaginationParams(pageStr, limitStr)
	if err != nil {
		ResponseCodeBody(w, http.StatusUnprocessableEntity, err.Error())
	}

	graph, err := graphstorage.ExportGraph(page, limit)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
	} else {

		response := types.GraphResponse{
			Graph: graph,
			Metadata: types.ExportMetadata{
				ExportedAt:        time.Now().Format(time.RFC3339),
				NodeCount:         len(graph.Nodes),
				RelationshipCount: len(graph.Relationships),
				Format:            "json",
			},
		}
		ResponseCodeBody(w, http.StatusOK, response)
	}
}

// Import 导入图数据
// @Summary      导入完整的图数据
// @Description  导入完整的图数据信息，包括节点和关系
// @Tags         graph
// @Accept       json
// @Produce      json
// @Param        graph   body   types.Graph  true  "图中的节点与关系信息"
// @Success      200  {object}  types.Graph "导入后的图信息"
// @Failure      422  {string}  string "参数错误"
// @Failure      500  {string}  string "内部请求错误"
// @Router       /knowledgegraph/graph [POST]
// POST - HTTP /knowledgegraph/graph
func (c *GraphController) Import(w http.ResponseWriter, r *http.Request) {
	graph, err := c.loadGraphData(w, r)

	if err != nil {
		return
	}

	nodes, err := graphstorage.ImportGraph(graph)
	if err != nil {
		ResponseCodeBody(w, http.StatusInternalServerError, err.Error())
	} else {
		ResponseCodeBody(w, http.StatusOK, nodes)
	}
}

func (c *GraphController) loadGraphData(w http.ResponseWriter, r *http.Request) (graph types.Graph, err error) {
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
