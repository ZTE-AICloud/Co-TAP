package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
	"uapregistry/logger"
	"uapregistry/servicemanager"
	"uapregistry/storage/consulagent/cache"
	t "uapregistry/types"
	"uapregistry/utils"

	"github.com/gorilla/mux"
)

// GET - HTTP /routes/{routeName}
func (s *HTTPServer) GetRouteHandler(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["routeName"]

	s.log.Debugf("Get route %s", name)

	route := cache.GetRoutesCache().GetRouteByName(name)
	if route == nil {
		Response404(w, fmt.Sprintf("route %s not found", name))
	}

	// marashal serviceUnit list
	j, err := json.Marshal(route)
	if err != nil {
		s.log.Errorf("Failed to marshal json data when build response: %v", err)
		Response500(w, "Failed to marshal json data when build response")
		return
	}

	Response200(w, j)
}

// GET - HTTP /routes
func (s *HTTPServer) GetAllRoutesHandler(w http.ResponseWriter, r *http.Request) {
	s.log.Debugf("Get all routes")

	indexStr := r.URL.Query().Get("index")
	waitTimeStr := r.URL.Query().Get("wait")
	if err := checkWaitTimeAndIndex(waitTimeStr, indexStr); err != nil {
		logger.GetLogger().Errorf("failed to get all routes:%v", err)
		Response400(w, err.Error())
	}

	var (
		routes    []*t.Route
		lastIndex uint64
	)

	index, _ := strconv.Atoi(indexStr)
	waitTime, _ := strconv.Atoi(waitTimeStr)
	// get all service
	if waitTime > 0 {
		routes, lastIndex = cache.GetRoutesCache().BlockingQueryAllRoutes(time.Duration(waitTime)*time.Second, uint64(index))
	} else {
		routes = cache.GetRoutesCache().GetAllRoutes()
	}

	// marashal serviceUnit list
	j, err := json.Marshal(routes)
	if err != nil {
		s.log.Errorf("Failed to marshal json data when build response: %v", err)
		Response500(w, "Failed to marshal json data when build response")
		return
	}

	Response200WithCustomHeader(w, j, map[string]string{"X-Uapregistry-Index": strconv.Itoa(int(lastIndex))})
}

// POST - HTTP /routes
func (s *HTTPServer) PostRouteHandler(w http.ResponseWriter, r *http.Request) {
	CheckUpdateSvc()
	defer UnlockUpdateSvc()

	//get input params
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.GetLogger().Errorf("Failed to ReadAll http body for PostRoute:%v", err)
		Response500(w, "Failed to ReadAll the http body for PostRoute")
		return
	}

	route := &t.Route{}
	if err = json.Unmarshal(body, route); err != nil {
		logger.GetLogger().Errorf("Add a route failed, Failed to Unmarshal body for PostRoute: %v body:%s", err, string(body))
		Response400(w, "Failed to Unmarshal body for PostRoute")
		return
	}
	utils.FillRouteDefaultValue(route)
	err = utils.CheckAndModifyRoute(route)
	if err != nil {
		logger.GetLogger().Errorf("Add a route failed, Failed to check and modify route: %v", err)
		Response422(w, err.Error())
		return
	}

	s.log.Infow(fmt.Sprintf("Add a route,body:%s", string(body)), "routeName", route.Name)

	if timeout := preDo(route.ID); timeout {
		s.log.Errorf("(routeName:%s) Waiting for the completion of last operation is timed out", route.ID)
		Response408(w)
		return
	}
	defer postDo(route.ID)

	svcManager := servicemanager.NewServiceManager()
	route, statusCode, err := svcManager.PostRoute(route)
	if err != nil {
		ResponseWithStatusCode(w, statusCode, err.Error())
		return
	}

	j, err := json.Marshal(route)
	if err != nil {
		s.log.Errorf("Failed to marshal json data when build response: %v", err)
		Response500(w, "Failed to marshal json data when build response")
		return
	}

	Response201(w, j)
}

// PUT - HTTP /routes/{routeName}
func (s *HTTPServer) UpdateRouteHandler(w http.ResponseWriter, r *http.Request) {
	CheckUpdateSvc()
	defer UnlockUpdateSvc()

	//get input params
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.GetLogger().Errorf("Failed to ReadAll http body for PostRoute:%v", err)
		Response500(w, "Failed to ReadAll the http body for PostRoute")
		return
	}

	route := &t.Route{}
	if err = json.Unmarshal(body, route); err != nil {
		logger.GetLogger().Errorf("Add a route failed, Failed to Unmarshal body for PostRoute: %v body:%s", err, string(body))
		Response400(w, "Failed to Unmarshal body for PostRoute")
		return
	}
	utils.FillRouteDefaultValue(route)
	err = utils.CheckAndModifyRoute(route)
	if err != nil {
		logger.GetLogger().Errorf("Add a route failed, Failed to check and modify route: %v", err)
		Response422(w, err.Error())
	}

	s.log.Infow(fmt.Sprintf("Add a route,body:%s", string(body)), "routeName", route.Name)

	if timeout := preDo(route.ID); timeout {
		s.log.Errorf("(routeName:%s) Waiting for the completion of last operation is timed out", route.ID)
		Response408(w)
		return
	}
	defer postDo(route.ID)

	svcManager := servicemanager.NewServiceManager()
	route, statusCode, err := svcManager.UpdateRoute(route)
	if err != nil {
		ResponseWithStatusCode(w, statusCode, err.Error())
		return
	}

	j, err := json.Marshal(route)
	if err != nil {
		s.log.Errorf("Failed to marshal json data when build response: %v", err)
		Response500(w, "Failed to marshal json data when build response")
		return
	}

	Response201(w, j)
}

// HTTP DELETE - /routes/{routeName}
func (s *HTTPServer) DeleteRouteHandler(w http.ResponseWriter, r *http.Request) {
	CheckUpdateSvc()
	defer UnlockUpdateSvc()

	routeName := mux.Vars(r)["routeName"]

	s.log.Infow("Delete a route", "routeName", routeName)
	if timeout := preDo(routeName); timeout {
		s.log.Errorf("(routeName:%s) Waiting for the completion of last operation is timed out", routeName)
		Response408(w)
		return
	}
	defer postDo(routeName)

	svcManager := servicemanager.NewServiceManager()
	err := svcManager.DeleteRoute(routeName)
	if err != nil {
		s.log.Errorf("(routeName:%s) Failed to delete a route: %v", routeName, err)
		ResponseWithStatusCode(w, http.StatusInternalServerError, "Failed to delete a route")
		return
	}

	Response204(w)
}
