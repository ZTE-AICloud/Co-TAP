package servicemanager

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"uapregistry/logger"
	agent "uapregistry/storage/consulagent"
	"uapregistry/storage/consulagent/cache"
	t "uapregistry/types"
	"uapregistry/utils"
)

type ServiceHandler struct {
	agent *agent.Agent
	log   logger.Logger
}

func NewServiceHandler(agt *agent.Agent) *ServiceHandler {
	return &ServiceHandler{agent: agt, log: logger.GetLogger()}
}

func (svcHandler *ServiceHandler) PostService(svc *t.Service, overWrite bool) (*t.Service, int, error) {
	//get serviceUnit from consul
	if !overWrite {
		svcs, err := svcHandler.GetServicesByName(svc.Name)
		if err != nil {
			svcHandler.log.Errorf("(serviceName:%s) Failed to GetServicesByName after register:%v", svc.Name, err)
			return nil, http.StatusInternalServerError, fmt.Errorf("failed to get service info after register:%v", err)
		}

		if isServiceIDExist(svcs, svc.ID) {
			svcHandler.log.Errorf("(serviceName:%s)Failed to create a service:service has existed", svc.Name)
			return nil, http.StatusConflict, errors.New("failed to create a service:service has existed")
		}
	}

	err := svcHandler.registerService(svc)
	if err != nil {
		svcHandler.log.Errorf("(serviceName:%s)Failed to register service: %v", svc.Name, err)
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to register service: %v", err)
	}

	return svc, http.StatusCreated, nil
}

func isServiceIDExist(svcs []*t.Service, id string) bool {
	for _, s := range svcs {
		if s.ID == id {
			return true
		}
	}

	return false
}

func buildPatchService(patchSvc, svcCurrent *t.Service) *t.Service {
	svcCurrent.UpdatedAt = t.Timestamp(time.Now())
	svcCurrent.EphemeralCheck = patchSvc.EphemeralCheck
	svcCurrent.PersistentCheck = patchSvc.PersistentCheck
	svcCurrent.Host = patchSvc.Host
	svcCurrent.Port = patchSvc.Port
	svcCurrent.Path = patchSvc.Path
	svcCurrent.Tags = patchSvc.Tags
	svcCurrent.ConnectTimeout = patchSvc.ConnectTimeout
	svcCurrent.WriteTimeout = patchSvc.WriteTimeout
	svcCurrent.ReadTimeout = patchSvc.ReadTimeout
	svcCurrent.Retries = patchSvc.Retries
	svcCurrent.AgentInfo = patchSvc.AgentInfo
	svcCurrent.AgentInfoUrl = patchSvc.AgentInfoUrl

	return patchSvc
}

func (svcHandler *ServiceHandler) PatchService(id string, svc *t.Service) (*t.Service, int, error) {
	//check service exit or not
	svcInCache := cache.GetServicesCache().GetServiceByID(id)
	if svcInCache == nil {
		return nil, http.StatusNotFound, errors.New("service not found")
	}

	svc = buildPatchService(svc, svcInCache)
	err := svcHandler.registerService(svc)
	if err != nil {
		svcHandler.log.Errorf("(serviceName:%s)Failed to register service: %v", svc.Name, err)
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to register service: %v", err)
	}

	return svc, http.StatusCreated, nil
}

func (svcHandler *ServiceHandler) DeleteServiceByID(id string) error {
	// check service exit or not
	svcInCache := cache.GetServicesCache().GetServiceByID(id)
	if svcInCache == nil {
		svcHandler.log.Warnf("service(id:%s) not found", id)
		return nil
	}

	return svcHandler.agent.DeregisterService(svcInCache.ID)
}

// get ServiceNodes
func (svcHandler *ServiceHandler) GetServicesByName(serviceName string) ([]*t.Service, error) {
	return svcHandler.agent.GetServicesByName(serviceName)
}

// registerServiceIfCreate
func (svcHandler *ServiceHandler) registerService(svc *t.Service) error {
	return svcHandler.agent.RegisterService(svc)
}

func (svcHandler *ServiceHandler) UpdateServiceStatus(id, status string) error {
	svcInCache := cache.GetServicesCache().GetServiceByID(id)
	if svcInCache == nil {
		return errors.New("service not found")
	}

	svcInCache.HealthStatus = status
	svcInCache.UpdatedAt = t.Timestamp(time.Now())
	err := svcHandler.registerService(svcInCache)
	if err != nil {
		svcHandler.log.Errorf("(id:%s)Failed to register service: %v", svcInCache.ID, err)
		return err
	}

	return nil
}

func (svcHandler *ServiceHandler) PostRoute(route *t.Route) (*t.Route, int, error) {
	//get serviceUnit from consul
	routeCur, err := svcHandler.GetRouteByName(route.Name)
	if err != nil {
		svcHandler.log.Errorf("(routeName:%s) Failed to GetRouteByName after register:%v", route.Name, err)
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to get route info after register:%v", err)
	}

	if routeCur != nil {
		return nil, http.StatusConflict, errors.New("failed to create a route:route has existed")
	}

	err = svcHandler.agent.UpdateRouteByName(route)
	if err != nil {
		svcHandler.log.Errorf("(serviceName:%s)Failed to register route: %v", route.Name, err)
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to register route: %v", err)
	}

	return route, http.StatusCreated, nil
}

func (svcHandler *ServiceHandler) UpdateRoute(route *t.Route) (*t.Route, int, error) {
	routeCur, err := svcHandler.GetRouteByName(route.Name)
	if err != nil {
		svcHandler.log.Errorf("(routeName:%s) Failed to GetRouteByName after register:%v", route.Name, err)
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to get route info after register:%v", err)
	}

	if routeCur == nil {
		return nil, http.StatusNotFound, errors.New("route not found")
	}

	route = buildUpdateRoute(route, routeCur)
	err = svcHandler.agent.UpdateRouteByName(route)
	if err != nil {
		svcHandler.log.Errorf("(serviceName:%s)Failed to register route: %v", route.Name, err)
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to register route: %v", err)
	}

	return route, http.StatusCreated, nil
}

func buildUpdateRoute(routeNew, routeCur *t.Route) *t.Route {
	routeNew.ID = routeCur.ID
	routeNew.Name = routeCur.Name
	routeNew.CreatedAt = routeCur.CreatedAt

	return routeNew
}

func (svcHandler *ServiceHandler) GetRouteByName(name string) (*t.Route, error) {
	return svcHandler.agent.GetRouteByName(name)
}

func (svcHandler *ServiceHandler) DeleteRoute(routeName string) error {
	k := utils.BuildRouteKey(routeName)
	return svcHandler.agent.DeleteKVPair(k)
}
