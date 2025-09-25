package utils

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/google/uuid"
	consulapi "github.com/hashicorp/consul/api"

	"uapregistry/logger"
	t "uapregistry/types"
)

func AgentService2Service(agentService *consulapi.AgentService) *t.Service {
	svc := &t.Service{}

	for _, tag := range agentService.Tags {
		if strings.HasPrefix(tag, "fullInfo:") {
			value := strings.TrimPrefix(tag, "fullInfo:")
			err := json.Unmarshal([]byte(value), &svc)
			if err != nil {
				logger.GetLogger().Errorf("Failed to unmarshal base tag:%v", err)
				return nil
			}
			break
		}
	}

	svc.Index = agentService.ModifyIndex
	if svc.HealthStatus == "" {
		svc.HealthStatus = "healthy"
	}

	return svc
}

func AgentServices2Services(agentServices map[string]*consulapi.AgentService) map[string]*t.Service {
	svcs := make(map[string]*t.Service)
	for _, agentService := range agentServices {
		svc := AgentService2Service(agentService)
		if svc != nil {
			svcs[agentService.ID] = svc
		}
	}

	return svcs
}

func CatalogService2Service(catalogService *consulapi.CatalogService) *t.Service {
	svc := &t.Service{}

	for _, tag := range catalogService.ServiceTags {
		if strings.HasPrefix(tag, "fullInfo:") {
			value := strings.TrimPrefix(tag, "fullInfo:")
			err := json.Unmarshal([]byte(value), &svc)
			if err != nil {
				logger.GetLogger().Errorf("Failed to unmarshal fullInfo tag:%v", err)
				return nil
			}
			break
		}
	}

	if svc.HealthStatus == "" {
		svc.HealthStatus = "healthy"
	}

	svc.Index = catalogService.ModifyIndex
	return svc
}

func CatalogServices2Services(css []*consulapi.CatalogService) []*t.Service {
	svcs := make([]*t.Service, 0)
	for _, cs := range css {
		svc := CatalogService2Service(cs)
		if svc != nil {
			svcs = append(svcs, svc)
		}
	}

	return svcs
}

// buildNodeName
func BuildNodeName() string {
	return "default-node"
}

func BuildServiceIDFromServiceUnit(svc *t.Service) string {
	return BuildServiceID(svc.Name, svc.Host, svc.Port)
}

func BuildServiceID(serviceName, host string, port int) string {
	portStr := strconv.Itoa(port)
	key := strings.Join([]string{serviceName, host, portStr}, "#")
	uid := uuid.NewSHA1(uuid.NameSpaceOID, []byte(key))
	return uid.String()
}

func BuildRouteUID(name string) string {
	key := strings.Join([]string{name}, "#")
	uid := uuid.NewSHA1(uuid.NameSpaceOID, []byte(key))
	return uid.String()
}

func buildServiceMeta(su *t.Service) map[string]string {
	var serviceMeta = make(map[string]string, 0)

	if su.Ephemeral {
		serviceMeta["ephemeral"] = "true"
	} else {
		serviceMeta["ephemeral"] = "false"
	}
	serviceMeta["name"] = su.Name
	serviceMeta["protocol"] = su.Protocol
	serviceMeta["agent_protocol"] = su.AgentProtocol

	return serviceMeta
}

// build serviceName in consul
func BuildServiceNameInConsul(serviceName string) (serviceNameInConsul string) {
	return serviceName
}

func GetMaxLastIndex(ccs []*consulapi.CatalogService) uint64 {
	var (
		index uint64
	)
	for _, cs := range ccs {
		if index < cs.ModifyIndex {
			index = cs.ModifyIndex
		}
	}

	return index
}

// build CatalogRegistration
func BuildCatalogRegistration(svc *t.Service) *consulapi.CatalogRegistration {
	return &consulapi.CatalogRegistration{
		Node:    BuildNodeName(),
		Address: "127.0.0.1",
		Service: buildAgentService(svc),
	}
}

// buildAgentService
func buildAgentService(svc *t.Service) *consulapi.AgentService {
	return &consulapi.AgentService{
		ID:      BuildServiceIDFromServiceUnit(svc),
		Service: BuildServiceNameInConsul(svc.Name),
		Tags:    buildServiceTags(svc),
		Address: svc.Host,
		Port:    svc.Port,
		Meta:    buildServiceMeta(svc),
	}
}

// build servicetags
func buildServiceTags(serviceNode *t.Service) []string {
	//build service tags
	baseTag := buildBaseTag(serviceNode)

	return []string{baseTag}
}

// buildBaseTag
func buildBaseTag(svc *t.Service) string {
	j, _ := json.Marshal(svc)
	return `fullInfo:` + string(j)
}
