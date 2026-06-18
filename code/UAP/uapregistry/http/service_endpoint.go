package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"uapregistry/logger"
	"uapregistry/servicemanager"
	"uapregistry/storage/consulagent/cache"
	"uapregistry/types"
	"uapregistry/utils"

	"github.com/gorilla/mux"
)

func checkWaitTimeAndIndex(waitTime, index string) error {
	if waitTime != "" && !utils.IsValidIntStr(waitTime) {
		return errors.New("wait is invalid")
	}
	if index != "" && !utils.IsValidIntStr(index) {
		return errors.New("index is invalid")
	}
	return nil
}

// GET - HTTP /services/{serviceName}
func (s *HTTPServer) GetServiceHandler(w http.ResponseWriter, r *http.Request) {
	serviceName := mux.Vars(r)["serviceName"]
	indexStr := r.URL.Query().Get("index")
	waitTimeStr := r.URL.Query().Get("wait")
	if err := checkWaitTimeAndIndex(waitTimeStr, indexStr); err != nil {
		logger.GetLogger().Errorf("failed to get all services:%v", err)
		Response400(w, err.Error())
	}

	var (
		svcs      []*types.Service
		lastIndex uint64
	)

	index, _ := strconv.Atoi(indexStr)
	waitTime, _ := strconv.Atoi(waitTimeStr)
	// get all service
	if waitTime > 0 {
		svcs, lastIndex = cache.GetServicesCache().BlockingQueryOneServices(serviceName, time.Duration(waitTime)*time.Second, uint64(index))
	} else {
		svcs = cache.GetServicesCache().GetServicesByServiceName(serviceName)
	}

	// marashal serviceUnit list
	j, err := json.Marshal(svcs)
	if err != nil {
		s.log.Errorf("Failed to marshal json data when build response: %v", err)
		Response500(w, "Failed to marshal json data when build response")
		return
	}

	Response200WithCustomHeader(w, j, map[string]string{"X-Uapregistry-Index": strconv.Itoa(int(lastIndex))})
}

// GET - HTTP /services
func (s *HTTPServer) GetAllServiceHandler(w http.ResponseWriter, r *http.Request) {
	s.log.Debugf("Get all services")

	indexStr := r.URL.Query().Get("index")
	waitTimeStr := r.URL.Query().Get("wait")
	if err := checkWaitTimeAndIndex(waitTimeStr, indexStr); err != nil {
		logger.GetLogger().Errorf("failed to get all services:%v", err)
		Response400(w, err.Error())
	}

	var (
		svcs      []*types.Service
		lastIndex uint64
	)

	index, _ := strconv.Atoi(indexStr)
	waitTime, _ := strconv.Atoi(waitTimeStr)
	// get all service
	if waitTime > 0 {
		svcs, lastIndex = cache.GetServicesCache().BlockingQueryAllServices(time.Duration(waitTime)*time.Second, uint64(index))
	} else {
		svcs = cache.GetServicesCache().GetAllService()
	}

	// marashal serviceUnit list
	j, err := json.Marshal(svcs)
	if err != nil {
		s.log.Errorf("Failed to marshal json data when build response: %v", err)
		Response500(w, "Failed to marshal json data when build response")
		return
	}

	Response200WithCustomHeader(w, j, map[string]string{"X-Uapregistry-Index": strconv.Itoa(int(lastIndex))})
}

// POST - HTTP /services/semantic/search
func (s *HTTPServer) SemanticSearchHandler(w http.ResponseWriter, r *http.Request) {
	//get input params
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.GetLogger().Errorf("Failed to ReadAll http body for semantic search:%v", err)
		Response500(w, "Failed to ReadAll the http body for semantic search")
		return
	}

	ssq := &types.SemanticSearchRequest{}
	if err = json.Unmarshal(body, ssq); err != nil {
		logger.GetLogger().Errorf("Semantic search failed, Failed to Unmarshal body for PostService: %v body:%s", err, string(body))
		Response400(w, "Failed to Unmarshal body for PostService")
		return
	}

	logger.GetLogger().Infof("semantic search, agent description: %v", ssq)
	// 3. 参数校验
	if ssq.Query == "" {
		logger.GetLogger().Error("Semantic search agent description cannot be empty")
		Response400(w, "Query cannot be empty")
		return
	}

	svcManager := servicemanager.NewServiceManager()
	svc, err := svcManager.SemanticSearch(ssq)
	if err != nil {
		logger.GetLogger().Error("Failed to Semantic search by agent description, err: %v", err)
		Response500(w, "Failed to Semantic search by agent description")
		return
	}

	// marashal serviceUnit list
	j, err := json.Marshal(svc)
	if err != nil {
		s.log.Errorf("Failed to marshal json data when build response: %v", err)
		Response500(w, "Failed to marshal json data when build response")
		return
	}

	Response201(w, j)
}


// POST - HTTP /services
func (s *HTTPServer) PostServiceHandler(w http.ResponseWriter, r *http.Request) {
	CheckUpdateSvc()
	defer UnlockUpdateSvc()

	//get input params
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.GetLogger().Errorf("Failed to ReadAll http body for PostService:%v", err)
		Response500(w, "Failed to ReadAll the http body for PostService")
		return
	}

	svc := &types.Service{}
	if err = json.Unmarshal(body, svc); err != nil {
		logger.GetLogger().Errorf("Add a service failed, Failed to Unmarshal body for PostService: %v body:%s", err, string(body))
		Response400(w, "Failed to Unmarshal body for PostService")
		return
	}
	utils.FillServiceDefaultValue(svc)
	err = utils.CheckAndModifyService(svc)
	if err != nil {
		logger.GetLogger().Errorf("Add a service failed, Failed to check and modify service: %v", err)
		Response422(w, err.Error())
		return
	}

	s.log.Infow(fmt.Sprintf("Add a service,body:%s", string(body)), "ServiceName", svc.Name)

	if timeout := preDo(svc.ID); timeout {
		s.log.Errorf("(serviceName:%s) Waiting for the completion of last operation is timed out", svc.ID)
		Response408(w)
		return
	}
	defer postDo(svc.ID)

	svcManager := servicemanager.NewServiceManager()
	suNew, statusCode, err := svcManager.PostService(svc, false)
	if err != nil {
		ResponseWithStatusCode(w, statusCode, err.Error())
		return
	}

	j, err := json.Marshal(suNew)
	if err != nil {
		s.log.Errorf("Failed to marshal json data when build response: %v", err)
		Response500(w, "Failed to marshal json data when build response")
		return
	}

	Response201(w, j)
}

// PATCH - HTTP /services/{id}
func (s *HTTPServer) PatchServiceHandler(w http.ResponseWriter, r *http.Request) {
	CheckUpdateSvc()
	defer UnlockUpdateSvc()

	id := mux.Vars(r)["id"]

	//get input params
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.GetLogger().Errorf("Failed to ReadAll http body for UpdateService:%v", err)
		Response500(w, "Failed to ReadAll the http body for UpdateService")
		return
	}

	svc := &types.Service{}
	if err = json.Unmarshal(body, svc); err != nil {
		logger.GetLogger().Errorf("Update a service failed, Failed to Unmarshal body for UpdateService: %v body:%s", err, string(body))
		Response400(w, "Failed to Unmarshal body for UpdateService")
		return
	}

	utils.FillServiceDefaultValue(svc)
	err = utils.CheckAndModifyService(svc)
	if err != nil {
		logger.GetLogger().Errorf("Update a service failed, Failed to check and modify service: %v", err)
		Response422(w, err.Error())
		return
	}

	if timeout := preDo(id); timeout {
		s.log.Errorf("(id:%s) Waiting for the completion of last operation is timed out", id)
		Response408(w)
		return
	}
	defer postDo(id)

	svcManager := servicemanager.NewServiceManager()
	serviceUnit, statusCode, err := svcManager.PatchService(id, svc)
	if err != nil {
		ResponseWithStatusCode(w, statusCode, "Failed to update a service")
		return
	}

	//marshal serviceUnit
	j, err := json.Marshal(serviceUnit)
	if err != nil {
		s.log.Errorf("Failed to marshal json data when build response: %v", err)
		Response500(w, "Failed to marshal json data when build response")
		return
	}

	Response201(w, j)
}

// HTTP DELETE - /services/{id}
func (s *HTTPServer) DeleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	CheckUpdateSvc()
	defer UnlockUpdateSvc()

	id := mux.Vars(r)["id"]

	s.log.Infow("Delete a service", "id", id)
	if timeout := preDo(id); timeout {
		s.log.Errorf("(id:%s) Waiting for the completion of last operation is timed out", id)
		Response408(w)
		return
	}
	defer postDo(id)

	svcManager := servicemanager.NewServiceManager()
	err := svcManager.DeleteServiceByID(id)
	if err != nil {
		s.log.Errorf("(id:%s) Failed to delete a service: %v", id, err)
		ResponseWithStatusCode(w, http.StatusInternalServerError, "Failed to delete a service")
		return
	}

	Response204(w)
}
