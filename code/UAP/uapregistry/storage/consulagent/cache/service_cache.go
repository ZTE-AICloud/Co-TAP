package cache

import (
	"sync"
	"time"

	consulapi "github.com/hashicorp/consul/api"

	"uapregistry/logger"
	agent "uapregistry/storage/consulagent"
	"uapregistry/types"
	"uapregistry/utils"
	"uapregistry/watch"
)

var (
	c *ServiceCache
)

type ServiceCache struct {
	sync.RWMutex
	index       uint64
	services    map[string]*types.Service
	agt         *agent.Agent
	w           *watch.NotifyGroup
	oneServiceW map[string]*watch.NotifyGroup
}

func InitServiceCache() {
	c = &ServiceCache{
		index:       0,
		services:    make(map[string]*types.Service),
		agt:         agent.GetLocalAgent(),
		w:           watch.NewNotifyGroup(),
		oneServiceW: make(map[string]*watch.NotifyGroup),
	}
	go c.watchAllServices()
}

func GetServicesCache() *ServiceCache {
	return c
}

func (cs *ServiceCache) addAWatch(serviceName string) {
	if cs.isWatchExist(serviceName) {
		return
	}

	cs.Lock()
	defer cs.Unlock()

	cs.oneServiceW[serviceName] = watch.NewNotifyGroup()
}

func (cs *ServiceCache) notifyAWatch(serviceName string) {
	if cs.isWatchExist(serviceName) {
		cs.oneServiceW[serviceName].NoBlockingNotify()
	}
}

func (cs *ServiceCache) isWatchExist(serviceName string) bool {
	cs.RLock()
	defer cs.RUnlock()

	_, exist := cs.oneServiceW[serviceName]
	return exist

}

func (cs *ServiceCache) watchAllServices() {
	q := &consulapi.QueryOptions{RequireConsistent: true, WaitIndex: 0}

	for {
		node, qm, err := cs.agt.CatalogNode("default-node", q)
		if err != nil {
			logger.GetLogger().Errorf("Failed to get node: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if qm.LastIndex == q.WaitIndex {
			continue
		}

		if node == nil {
			q.WaitIndex = qm.LastIndex
			continue
		}

		svcs := utils.AgentServices2Services(node.Services)
		addSvcs, updateSvcs, delSvcs := utils.GetDiffServiceNames(cs.services, svcs)
		printDiffServices(addSvcs, updateSvcs, delSvcs)

		c.updateServices(svcs)
		q.WaitIndex = qm.LastIndex

		cs.w.NoBlockingNotify()
		for _, serviceName := range addSvcs {
			cs.notifyAWatch(serviceName)
		}
		for _, serviceName := range updateSvcs {
			cs.notifyAWatch(serviceName)
		}
		for _, serviceName := range delSvcs {
			cs.notifyAWatch(serviceName)
		}

		time.Sleep(1 * time.Second)
	}
}

func printDiffServices(addSvcs, updateSvcs, delSvcs []string) {
	logger.GetLogger().Infof("add:%d, update:%d, del:%d services", len(addSvcs), len(updateSvcs), len(delSvcs))
	for _, serviceName := range addSvcs {
		logger.GetLogger().Infof("add service: %s", serviceName)
	}
	for _, serviceName := range updateSvcs {
		logger.GetLogger().Infof("update service: %s", serviceName)
	}
	for _, serviceName := range delSvcs {
		logger.GetLogger().Infof("del service: %s", serviceName)
	}
}

func (cs *ServiceCache) updateServices(svcs map[string]*types.Service) {
	cs.Lock()
	defer cs.Unlock()

	cs.index = utils.CalculateServiceMapIndex(svcs)
	cs.services = svcs
}

func (cs *ServiceCache) GetServicesByServiceName(name string) []*types.Service {
	cs.RLock()
	defer cs.RUnlock()

	svcs := make([]*types.Service, 0)
	for _, svc := range cs.services {
		if svc.Name == name {
			svcs = append(svcs, svc)

		}
	}
	return svcs
}

func (cs *ServiceCache) GetServiceByID(id string) *types.Service {
	cs.RLock()
	defer cs.RUnlock()

	return cs.services[id]
}

func (cs *ServiceCache) GetAllService() []*types.Service {
	cs.RLock()
	defer cs.RUnlock()

	svcs := make([]*types.Service, 0)
	for _, svc := range cs.services {
		svcs = append(svcs, svc)
	}
	return svcs
}

func (cs *ServiceCache) BlockingQueryAllServices(waitTime time.Duration, index uint64) ([]*types.Service, uint64) {
	if index != cs.index {
		return cs.GetAllService(), cs.index
	}

	var (
		ch = make(chan struct{})
	)

	cs.w.Wait(ch)
	defer cs.w.Clear(ch)

	timeAfater := time.NewTimer(waitTime)
	defer timeAfater.Stop()

	select {
	case <-ch:
		logger.GetLogger().Infof("BlockingQueryAllServices:services has changed")
		break
	case <-timeAfater.C:
		logger.GetLogger().Infof("BlockingQueryAllServices:return all service after %v", waitTime)
	}

	return cs.GetAllService(), cs.index
}

func (cs *ServiceCache) BlockingQueryOneServices(serviceName string, waitTime time.Duration, index uint64) ([]*types.Service, uint64) {
	svcs := cs.GetServicesByServiceName(serviceName)
	lastIndex := utils.CalculateServicesIndex(svcs)

	if index != lastIndex {
		return svcs, lastIndex
	}

	var (
		ch = make(chan struct{})
	)

	cs.addAWatch(serviceName)

	cs.oneServiceW[serviceName].Wait(ch)
	defer cs.oneServiceW[serviceName].Clear(ch)

	timeAfater := time.NewTimer(waitTime)
	defer timeAfater.Stop()

	select {
	case <-ch:
		logger.GetLogger().Infof("BlockingQueryOneServices(%s):services has changed", serviceName)
		break
	case <-timeAfater.C:
		logger.GetLogger().Infof("BlockingQueryOneServices(%s):return all service after %v", serviceName, waitTime)
	}

	svcs = cs.GetServicesByServiceName(serviceName)
	lastIndex = utils.CalculateServicesIndex(svcs)
	return svcs, lastIndex
}
