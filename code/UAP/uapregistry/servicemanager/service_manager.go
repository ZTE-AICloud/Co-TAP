package servicemanager

import (
	"uapregistry/logger"
	agent "uapregistry/storage/consulagent"
	t "uapregistry/types"
)

var (
	svcManager *ServiceManager
)

type ServiceManager struct {
	localHandler *ServiceHandler
	log          logger.Logger
}

var NewServiceManager = func() *ServiceManager {
	if svcManager != nil {
		return svcManager
	}

	localHandler := NewServiceHandler(agent.GetLocalAgent())
	svcManager = BuildServiceManagerByHandler(localHandler, nil)

	return svcManager
}

func BuildServiceManagerByHandler(localHandler, remoteHandler *ServiceHandler) *ServiceManager {
	return &ServiceManager{
		log:          logger.GetLogger(),
		localHandler: localHandler,
	}
}

func (svcManager *ServiceManager) GetService(serviceName string) (svcs []*t.Service, err error) {
	return svcManager.localHandler.GetServicesByName(serviceName)
}

func (svcManager *ServiceManager) PostService(su *t.Service, overWrite bool) (*t.Service, int, error) {
	serviceUnit, statusCode, err := svcManager.localHandler.PostService(su, overWrite)
	return serviceUnit, statusCode, err
}

func (svcManager *ServiceManager) PatchService(id string, svc *t.Service) (*t.Service, int, error) {
	return svcManager.localHandler.PatchService(id, svc)
}

func (svcManager *ServiceManager) DeleteServiceByID(id string) error {
	return svcManager.localHandler.DeleteServiceByID(id)
}

func (svcManager *ServiceManager) UpdateServiceStatus(id, status string) error {
	return svcManager.localHandler.UpdateServiceStatus(id, status)
}

func (svcManager *ServiceManager) PostRoute(route *t.Route) (*t.Route, int, error) {
	return svcManager.localHandler.PostRoute(route)
}

func (svcManager *ServiceManager) UpdateRoute(route *t.Route) (*t.Route, int, error) {
	return svcManager.localHandler.UpdateRoute(route)
}

func (svcManager *ServiceManager) DeleteRoute(routeName string) error {
	return svcManager.localHandler.DeleteRoute(routeName)
}
