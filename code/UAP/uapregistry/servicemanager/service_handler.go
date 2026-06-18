package servicemanager

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"uapregistry/logger"
	chroma "uapregistry/storage/chromaagent"
	agent "uapregistry/storage/consulagent"
	"uapregistry/storage/consulagent/cache"
	"uapregistry/types"
	t "uapregistry/types"
	"uapregistry/utils"
)

type ServiceHandler struct {
	chromaAgent *chroma.ChromaAgentManager
	agent       *agent.Agent
	log         logger.Logger
}

func NewServiceHandler(ca *chroma.ChromaAgentManager, agt *agent.Agent) *ServiceHandler {
	return &ServiceHandler{chromaAgent: ca, agent: agt, log: logger.GetLogger()}
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

func (svcHandler *ServiceHandler) SemanticSearch(ssq *types.SemanticSearchRequest) ([]*t.Service, error) {
	logger.GetLogger().Info("query agent id by agent description")
	svcIDs, err := svcHandler.chromaAgent.QueryTopNAgents(ssq)
	if err != nil {
		svcHandler.log.Errorf("(agentDescription:%s)Failed to semantic search: %v", ssq.Query, err)
		return nil, err
	}

	var services []*types.Service
	for _, svcID := range svcIDs{
		services = append(services, cache.GetServicesCache().GetServiceByID(svcID))
	}
	
	return services, nil
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
	err := svcHandler.agent.DeregisterService(svcInCache.ID)
	if err != nil {
		svcHandler.log.Errorf("service(id:%s) delete in consul failed", id)
		return err
	}
	return svcHandler.chromaAgent.DeleteAgent(svcInCache.ID)
}

// get ServiceNodes
func (svcHandler *ServiceHandler) GetServicesByName(serviceName string) ([]*t.Service, error) {
	return svcHandler.agent.GetServicesByName(serviceName)
}

// registerServiceIfCreate
func (svcHandler *ServiceHandler) registerService(svc *t.Service) error {
	err := svcHandler.agent.RegisterService(svc)
	if err != nil {
		return err
	}

	description := buildDescription(svc)
	if description != "" {
		return svcHandler.chromaAgent.RegisterAgent(svc.ID, description, svc.AgentInfo.A2AAgentCard)
	}
	return nil
}

func buildDescription(svc *t.Service) string {
	var buf strings.Builder
	agentCard := svc.AgentInfo.A2AAgentCard

	// 1. 遍历所有技能，按字段带前缀拼接
	for _, skill := range agentCard.Skills {
		// 技能名称：非空才写入
		if skill.Name != "" {
			buf.WriteString("技能名称：")
			buf.WriteString(skill.Name)
			buf.WriteByte('\n')
		}
		// 技能描述：非空才写入
		if skill.Description != "" {
			buf.WriteString("技能描述：")
			buf.WriteString(skill.Description)
			buf.WriteByte('\n')
		}
		// 技能标签：数组非空才写入，用逗号拼接
		if len(skill.Tags) > 0 {
			buf.WriteString("技能标签：")
			buf.WriteString(strings.Join(skill.Tags, ", "))
			buf.WriteByte('\n')
		}
		// 技能示例：数组非空才写入，用逗号拼接
		if len(skill.Examples) > 0 {
			buf.WriteString("技能示例：")
			buf.WriteString(strings.Join(skill.Examples, ", "))
			buf.WriteByte('\n')
		}
		// 不同技能之间空一行，做语义隔离
		buf.WriteByte('\n')
	}

	// 2. Agent 自身字段拼接
	if agentCard.Name != "" {
		buf.WriteString("Agent名称：")
		buf.WriteString(agentCard.Name)
		buf.WriteByte('\n')
	}
	if agentCard.Description != "" {
		buf.WriteString("Agent描述：")
		buf.WriteString(agentCard.Description)
		buf.WriteByte('\n')
	}

	// 去掉末尾多余的换行，得到最终文本
	return strings.TrimRight(buf.String(), "\n")
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
