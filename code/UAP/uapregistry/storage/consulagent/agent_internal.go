package agent

import (
	"time"
	"uapregistry/logger"

	consulapi "github.com/hashicorp/consul/api"
)

const (
	retryCount    = 2
	retryInterval = 2
)

func (agent *Agent) catalogRegister(reg *consulapi.CatalogRegistration, q *consulapi.WriteOptions) (wm *consulapi.WriteMeta, err error) {
	for i := 1; i <= retryCount; i++ {
		wm, err = agent.catalogEP.Register(reg, q)
		if err != nil {
			time.Sleep(retryInterval * time.Second)
			continue
		}
		break
	}

	logCatalogRegister(reg, err)
	return
}

func logCatalogRegister(reg *consulapi.CatalogRegistration, err error) {
	var (
		serviceID string
		checkID   string
	)

	if reg.Service != nil {
		serviceID = reg.Service.ID
	}
	if reg.Check != nil {
		checkID = reg.Check.CheckID
	}

	if err != nil {
		logger.GetLogger().Errorf("failed catalogRegister(node:%s,serviceID:%s,checkID:%s):%v", reg.Node, serviceID, checkID, err)
		return
	}
	logger.GetLogger().Infof("catalogRegister(node:%s,serviceID:%s,checkID:%s) successfully", reg.Node, serviceID, checkID)
}

func (agent *Agent) catalogDeregister(dereg *consulapi.CatalogDeregistration, q *consulapi.WriteOptions) (wm *consulapi.WriteMeta, err error) {
	for i := 1; i <= retryCount; i++ {
		wm, err = agent.catalogEP.Deregister(dereg, q)
		if err != nil {
			time.Sleep(retryInterval * time.Second)
			continue
		}
		agent.logWriter.Infof("catalogDeregister(node:%s,id:%s) successfully", dereg.Node, dereg.ServiceID)
		break
	}

	if err != nil {
		agent.logWriter.Errorf("failed catalogDeregister(node:%s,id:%s):%v", dereg.Node, dereg.ServiceID, err)
	}

	return
}

func (agent *Agent) catalogService(service, tag string, q *consulapi.QueryOptions) (out []*consulapi.CatalogService, qm *consulapi.QueryMeta, err error) {
	for i := 1; i <= retryCount; i++ {
		out, qm, err = agent.catalogEP.Service(service, tag, q)
		if err != nil {
			if q != nil {
				return
			}
			agent.logWriter.Warnf("catalogService(%s) failed, retrying...: %v", service, err)
			time.Sleep(retryInterval * time.Second)
			continue
		}
		break
	}

	return
}

func (agent *Agent) CatalogNode(nodeName string, q *consulapi.QueryOptions) (out *consulapi.CatalogNode, qm *consulapi.QueryMeta, err error) {
	for i := 1; i <= retryCount; i++ {
		out, qm, err = agent.catalogEP.Node(nodeName, q)
		if err != nil {
			if q != nil {
				return
			}
			agent.logWriter.Warnf("catalogServices failed, retrying...: %v", err)
			time.Sleep(retryInterval * time.Second)
			continue
		}
		break
	}

	return
}

func (agent *Agent) kvPut(p *consulapi.KVPair, q *consulapi.WriteOptions) (wm *consulapi.WriteMeta, err error) {
	for i := 1; i <= retryCount; i++ {
		wm, err = agent.kvEp.Put(p, q)
		if err != nil {
			time.Sleep(retryInterval * time.Second)
			continue
		}
		agent.logWriter.Infof("kvPut(key:%s) successfully", p.Key)
		break
	}

	if err != nil {
		agent.logWriter.Errorf("failed kvPut(key:%s):%v", p.Key, err)
	}

	return
}

func (agent *Agent) kvGet(key string, q *consulapi.QueryOptions) (out *consulapi.KVPair, qm *consulapi.QueryMeta, err error) {
	for i := 1; i <= retryCount; i++ {
		out, qm, err = agent.kvEp.Get(key, q)
		if err != nil {
			if q != nil {
				return
			}
			agent.logWriter.Warnf("kvGet(key:%s) failed, retrying...: %v", key, err)
			time.Sleep(retryInterval * time.Second)
			continue
		}
		break
	}

	return
}

func (agent *Agent) kvDelete(key string, w *consulapi.WriteOptions) (wm *consulapi.WriteMeta, err error) {
	for i := 1; i <= retryCount; i++ {
		wm, err = agent.kvEp.Delete(key, w)
		if err != nil {
			time.Sleep(retryInterval * time.Second)
			continue
		}
		agent.logWriter.Infof("kvDelete(key:%s) successfully", key)
		break
	}

	if err != nil {
		agent.logWriter.Errorf("failed kvDelete(key:%s):%v", key, err)
	}

	return
}

func (agent *Agent) KVList(prefix string, q *consulapi.QueryOptions) (out consulapi.KVPairs, qm *consulapi.QueryMeta, err error) {
	for i := 1; i <= retryCount; i++ {
		out, qm, err = agent.kvEp.List(prefix, q)
		if err != nil {
			if q != nil {
				return
			}
			agent.logWriter.Warnf("kvList(prefix:%s) failed, retrying...: %v", prefix, err)
			time.Sleep(retryInterval * time.Second)
			continue
		}
		break
	}

	return
}
