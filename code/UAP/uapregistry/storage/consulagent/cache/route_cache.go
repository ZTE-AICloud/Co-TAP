package cache

import (
	"sync"
	"time"
	"uapregistry/logger"
	agent "uapregistry/storage/consulagent"
	"uapregistry/types"
	"uapregistry/utils"
	"uapregistry/watch"

	consulapi "github.com/hashicorp/consul/api"
)

const (
	prefix = "v1/routes"
)

var (
	routeCache *RouteCache
)

type RouteCache struct {
	sync.RWMutex
	index  uint64
	routes []*types.Route
	agt    *agent.Agent
	w      *watch.NotifyGroup
}

func InitRouteCache() {
	routeCache = &RouteCache{
		index:  0,
		routes: make([]*types.Route, 0),
		agt:    agent.GetLocalAgent(),
		w:      watch.NewNotifyGroup(),
	}
	go routeCache.watchRoutes()
}

func GetRoutesCache() *RouteCache {
	return routeCache
}

func (rc *RouteCache) watchRoutes() {
	q := &consulapi.QueryOptions{RequireConsistent: true, WaitIndex: 0}

	for {
		kvps, qm, err := rc.agt.KVList(prefix, q)
		if err != nil {
			logger.GetLogger().Errorf("Failed to get kvs(prefix:%s): %v", prefix, err)
			time.Sleep(5 * time.Second)
			continue
		}

		if qm.LastIndex == q.WaitIndex {
			continue
		}

		routes := utils.ConsulKVPairs2Routes(kvps)

		added, updated, deleted := utils.GetDiffRoutes(rc.routes, routes)
		printDiffRoutes(added, updated, deleted)

		q.WaitIndex = qm.LastIndex
		rc.updateRoutes(routes)

		rc.w.NoBlockingNotify()

		time.Sleep(1 * time.Second)
	}
}

func (rc *RouteCache) updateRoutes(routes []*types.Route) {
	rc.Lock()
	defer rc.Unlock()

	rc.index = utils.CalculateRoutesIndex(routes)
	rc.routes = routes
}

func (rc *RouteCache) BlockingQueryAllRoutes(waitTime time.Duration, index uint64) ([]*types.Route, uint64) {
	if index != rc.index {
		return rc.GetAllRoutes(), rc.index
	}

	var (
		ch = make(chan struct{})
	)

	rc.w.Wait(ch)
	defer rc.w.Clear(ch)

	timeAfater := time.NewTimer(waitTime)
	defer timeAfater.Stop()

	select {
	case <-ch:
		logger.GetLogger().Infof("BlockingQueryAllRoutes:services has changed")
		break
	case <-timeAfater.C:
		logger.GetLogger().Infof("BlockingQueryAllRoutes:return all service after %v", waitTime)
	}

	return rc.GetAllRoutes(), rc.index
}

func (rc *RouteCache) GetAllRoutes() []*types.Route {
	rc.RLock()
	defer rc.RUnlock()

	return rc.routes
}

func (rc *RouteCache) GetRouteByName(name string) *types.Route {
	rc.RLock()
	defer rc.RUnlock()

	for _, route := range rc.routes {
		if route.Name == name {
			return route
		}
	}

	return nil
}

func printDiffRoutes(add, update, del []string) {
	logger.GetLogger().Infof("add:%d, update:%d, del:%d routes", len(add), len(update), len(del))
	for _, name := range add {
		logger.GetLogger().Infof("add route: %s", name)
	}
	for _, name := range update {
		logger.GetLogger().Infof("update route: %s", name)
	}
	for _, name := range del {
		logger.GetLogger().Infof("del route: %s", name)
	}
}
