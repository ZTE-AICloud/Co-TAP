package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	consulapi "github.com/hashicorp/consul/api"

	"uapregistry/logger"
	t "uapregistry/types"
	"uapregistry/utils"
)

const (
	// Maximum number of attempts to connect to consul server.
	maxConnectAttempts = 12

	// Maximum time await the response from agent to log the event
	MaxRespondTime = 500
)

type Agent struct {
	sync.RWMutex
	logThreshold time.Duration
	logWriter    logger.Logger
	client       *consulapi.Client
	catalogEP    *consulapi.Catalog
	healthEP     *consulapi.Health
	agentEP      *consulapi.Agent
	statusEP     *consulapi.Status
	kvEp         *consulapi.KV
}

func Create(cfg *Config) (*Agent, error) {
	var err error
	var consulClient *consulapi.Client
	if cfg == nil {
		return nil, errors.New("configuration object is nil")
	}

	lw := logger.GetLogger()
	if lw == nil {
		return nil, errors.New("failed to create logger object")
	}

	if consulClient, err = NewConsulClient(cfg.ConsulAgent, maxConnectAttempts); err != nil {
		return nil, errors.New("failed to create consul api client")
	}

	agent, err := BuildAgentByConsulClient(consulClient)
	if err != nil {
		logger.GetLogger().Errorf("failed to BuildAgentByConsulClient:%v", err)
		return nil, err
	}

	return agent, nil
}

func BuildAgentByConsulConfig(consulAgentConfig string) (*Agent, error) {
	consulClient, err := NewConsulClient(consulAgentConfig, 1)
	if err != nil {
		return nil, err
	}

	agent, err := BuildAgentByConsulClient(consulClient)
	if err != nil {
		logger.GetLogger().Errorf("failed to BuildAgentByConsulClient:%v", err)
		return nil, err
	}
	return agent, nil
}

func BuildAgentByConsulClient(consulClient *consulapi.Client) (*Agent, error) {
	if consulClient == nil {
		return nil, errors.New("consul-client is nil")
	}
	if consulClient.Catalog() == nil {
		return nil, errors.New("consul-client CatalogEP is nil")
	}
	if consulClient.Health() == nil {
		return nil, errors.New("consul-client HealthEP is nil")
	}
	if consulClient.Agent() == nil {
		return nil, errors.New("consul-client AgentEP is nil")
	}
	if consulClient.Status() == nil {
		return nil, errors.New("consul-client StatusEP is nil")
	}
	if consulClient.KV() == nil {
		return nil, errors.New("consul-client KVEP is nil")
	}

	agent := &Agent{
		logThreshold: MaxRespondTime,
		logWriter:    logger.GetLogger(),
		client:       consulClient,
		catalogEP:    consulClient.Catalog(),
		healthEP:     consulClient.Health(),
		agentEP:      consulClient.Agent(),
		statusEP:     consulClient.Status(),
		kvEp:         consulClient.KV(),
	}

	return agent, nil
}

func (agent *Agent) LogThreshold() time.Duration {
	return agent.logThreshold
}

func (agent *Agent) GetConsulClient() *consulapi.Client {
	return agent.client
}

func (agent *Agent) GetServicesByName(serviceName string) ([]*t.Service, error) {
	serviceNameInConsul := utils.BuildServiceNameInConsul(serviceName)
	css, _, err := agent.catalogService(serviceNameInConsul, "", nil)
	if err != nil {
		return nil, err
	}
	svcs := utils.CatalogServices2Services(css)
	return svcs, nil
}

func (agent *Agent) RegisterService(svc *t.Service) error {
	return agent.registerServiceWithCatalogMode(svc)
}

// deregister service
func (agent *Agent) DeregisterService(id string) error {
	catalogDreg := &consulapi.CatalogDeregistration{
		Node:      utils.BuildNodeName(),
		ServiceID: id,
	}
	startTime := time.Now()
	_, err := agent.catalogDeregister(catalogDreg, nil)
	if err != nil {
		agent.logWriter.Errorf("(Node:%s,ServiceID:%s) Failed to catalog Deregister:%v", catalogDreg.Node, catalogDreg.ServiceID, err)
		return err
	}
	logIfLongResponseTime(startTime, fmt.Sprintf("agent.catalogEP.Deregister(Node:%s,ServiceID:%s) in deregisterServiceWithCatalogMode", catalogDreg.Node, catalogDreg.ServiceID))

	return nil
}

// updateKVPair
func (agent *Agent) UpdateKVPair(k, v string) error {
	startTime := time.Now()
	kvPair := &consulapi.KVPair{Key: k, Value: []byte(v)}
	_, err := agent.kvPut(kvPair, nil)
	if err != nil {
		agent.logWriter.Errorf("(k:%s,v:%s) Put a KVPair failed:%v", k, v, err)
	}
	logIfLongResponseTime(startTime, fmt.Sprintf("agent.kvEp.Put(key:%s) in UpdateKVPair", k))

	return err
}

// getKVPair
func (agent *Agent) GetKVPair(k string) (string, error) {
	value, err := agent.GetValueFromConsulKV(k)
	if err != nil {
		agent.logWriter.Errorf("Failed to GetValueFromConsulKV(key:%s)", k)
		return "", err
	}
	if value == nil {
		return "", nil
	}

	return string(value), nil
}

func (agent *Agent) GetValueFromConsulKV(k string) ([]byte, error) {
	startTime := time.Now()
	kvPair, _, err := agent.kvGet(k, &consulapi.QueryOptions{RequireConsistent: true})
	if err != nil {
		agent.logWriter.Errorf("(key:%s) Failed to get kv from consul:%v", k, err)
		return nil, err
	}
	logIfLongResponseTime(startTime, fmt.Sprintf("agent.kvEp.Get(key:%s) in GetValueFromConsulKV", k))

	if kvPair == nil {
		//		agent.logWriter.Debugf("(key:%s) Can not find kv in consul kvStore", k)
		return nil, nil
	}

	return kvPair.Value, nil
}

func (agent *Agent) GetConsulKV(k string) (*consulapi.KVPair, error) {
	startTime := time.Now()
	kvPair, _, err := agent.kvGet(k, &consulapi.QueryOptions{RequireConsistent: true})
	if err != nil {
		agent.logWriter.Errorf("(key:%s) Failed to get kv from consul:%v", k, err)
		return nil, err
	}
	logIfLongResponseTime(startTime, fmt.Sprintf("agent.kvEp.Get(key:%s) in GetConsulKV", k))

	return kvPair, nil
}

func (agent *Agent) DeleteKVPair(k string) error {
	startTime := time.Now()
	_, err := agent.kvDelete(k, nil)
	if err != nil {
		agent.logWriter.Errorf("(key:%s) Failed to delete kv from consul:%v", k, err)
		return err
	}
	logIfLongResponseTime(startTime, fmt.Sprintf("agent.kvEp.Delete(key:%s) in DeleteKVPair", k))

	return nil
}

func (agent *Agent) registerServiceWithCatalogMode(su *t.Service) error {
	var err error
	cr := utils.BuildCatalogRegistration(su)
	startTime := time.Now()
	_, err = agent.catalogRegister(cr, nil)
	if err != nil {
		agent.logWriter.Errorf("(serviceId:%s) Failed to register a service:%v", cr.Service.ID, err)
		return err
	}
	logIfLongResponseTime(startTime, fmt.Sprintf("agent.catalogEP.Register(Node:%s,ServiceID:%s) in registerServiceWithCatalogMode", cr.Node, cr.Service.ID))

	return nil
}

func (agent *Agent) GetRouteByName(serviceName string) (*t.Route, error) {
	key := utils.BuildRouteKey(serviceName)
	kvp, err := agent.GetConsulKV(key)
	if err != nil {
		return nil, err
	}

	route := utils.ConsulKVPair2Route(kvp)

	return route, nil
}

func (agent *Agent) UpdateRouteByName(route *t.Route) error {
	value, err := json.Marshal(route)
	if err != nil {
		return err
	}

	return agent.UpdateKVPair(utils.BuildRouteKey(route.Name), string(value))
}

func logIfLongResponseTime(startTime time.Time, interfaceInfo string) {
	responseTime := time.Since(startTime)
	if responseTime > 5*time.Second {
		logger.GetLogger().Warnf("consul response time is too long,response time:%s,interface info:%s", responseTime.String(), interfaceInfo)
	}
}
