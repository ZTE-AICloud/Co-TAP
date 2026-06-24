package http

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"uapregistry/http/graph"
	"uapregistry/logger"
)

const (
	operationTimeOut = 60
)

var (
	locker               = &sync.RWMutex{}
	updateSvcLocker      = &sync.RWMutex{}
	operationSyncChannel = make(map[string]chan struct{})
)

func LockUpdateSvc() {
	updateSvcLocker.Lock()
	logger.GetLogger().Warn("update service is locked")
}

func AllowUpdateSvc() {
	updateSvcLocker.Unlock()
	logger.GetLogger().Warn("update service is unlocked")
}

func CheckUpdateSvc() {
	updateSvcLocker.RLock()
}

func UnlockUpdateSvc() {
	updateSvcLocker.RUnlock()
}

func preDo(key string) (timeout bool) {
	var syncChan chan struct{}
	locker.RLock()
	if _, ok := operationSyncChannel[key]; ok {
		syncChan = operationSyncChannel[key]
	}
	locker.RUnlock()

	if syncChan == nil {
		locker.Lock()
		if _, ok := operationSyncChannel[key]; ok {
			locker.Unlock()
		} else {
			syncChan = make(chan struct{}, 1)
			operationSyncChannel[key] = syncChan
			locker.Unlock()
			return
		}
	}

	timer := time.NewTimer(operationTimeOut * time.Second)
	select {
	case <-syncChan:
	case <-timer.C:
		timeout = true
	}
	return
}

func postDo(key string) {
	locker.RLock()
	syncChan := operationSyncChannel[key]
	locker.RUnlock()
	syncChan <- struct{}{}
}

type HTTPServer struct {
	log    logger.Logger
	router *mux.Router
	addr   string
}

func StartHTTPServer(conf *Config, errCh chan error) []*http.Server {
	servers := make([]*http.Server, 0, 5)
	for _, value := range conf.HTTPIPs {
		httpAddr := value + ":" + conf.HTTPPort

		srv := &HTTPServer{
			log:    logger.GetLogger(),
			router: mux.NewRouter().StrictSlash(false),
			addr:   httpAddr,
		}

		srv.registerHandlers()

		srv.addMiddlewares()

		server := &http.Server{
			Addr:              httpAddr,
			Handler:           srv.router,
			MaxHeaderBytes:    16 * (1 << 10),
			IdleTimeout:       120 * time.Second,
			ReadHeaderTimeout: 60 * time.Second,
		}

		go func() {
			if err := server.ListenAndServe(); err != nil {
				errCh <- err
				return
			}
		}()
		servers = append(servers, server)
	}
	return servers
}

func (s *HTTPServer) registerHandlers() {
	s.registGraphHandlers()

	s.router.HandleFunc("/health", s.HealthCheckHandler).Methods("GET")

	s.router.HandleFunc("/services", s.GetAllServiceHandler).Methods("GET")
	s.router.HandleFunc("/services", s.PostServiceHandler).Methods("POST")
	s.router.HandleFunc("/services/{id}", s.DeleteServiceHandler).Methods("DELETE")
	s.router.HandleFunc("/services/{id}", s.PatchServiceHandler).Methods("PATCH")
	s.router.HandleFunc("/services/{id}/renewal", s.RenewalHandler).Methods("PUT")
	s.router.HandleFunc("/services/{id}/healthcheck", s.ServiceHealthCheckHandler).Methods("GET")
	s.router.HandleFunc("/services/{serviceName}", s.GetServiceHandler).Methods("GET")
	s.router.HandleFunc("/api/v1/agents/search", s.SemanticSearchHandler).Methods("POST")

	s.router.HandleFunc("/routes", s.PostRouteHandler).Methods("POST")
	s.router.HandleFunc("/routes", s.GetAllRoutesHandler).Methods("GET")
	s.router.HandleFunc("/routes/{routeName}", s.GetRouteHandler).Methods("GET")
	s.router.HandleFunc("/routes/{routeName}", s.DeleteRouteHandler).Methods("DELETE")
	s.router.HandleFunc("/routes/{routeName}", s.UpdateRouteHandler).Methods("PUT")
}

func (s *HTTPServer) registGraphHandlers() {
	logger.GetLogger().Info("regist agent graph handlers")
	graphController := &graph.GraphController{}
	// /knowledgegraph/graph?page=0&limit=1000
	s.router.HandleFunc("/knowledgegraph/graph", graphController.Export).Methods(http.MethodGet)
	s.router.HandleFunc("/knowledgegraph/graph", graphController.Import).Methods(http.MethodPost)

	nodeController := &graph.NodeController{}
	s.router.HandleFunc("/knowledgegraph/nodes", nodeController.Create).Methods(http.MethodPost)
	s.router.HandleFunc("/knowledgegraph/nodes/bulk", nodeController.CreateBulk).Methods(http.MethodPost)
	s.router.HandleFunc("/knowledgegraph/nodes/{elementId}", nodeController.Put).Methods(http.MethodPut)
	s.router.HandleFunc("/knowledgegraph/nodes/{elementId}", nodeController.Delete).Methods(http.MethodDelete)

	// /knowledgegraph/nodes/{agentCardName}/namespace/{ns1}?cluster=c1
	s.router.HandleFunc("/knowledgegraph/nodes/{agentCardName}/namespace/{ns1}", nodeController.GetNodesByName).Methods(http.MethodGet)
	s.router.HandleFunc("/knowledgegraph/nodes/{elementId}", nodeController.GetNodeByID).Methods(http.MethodGet)
	// /knowledgegraph/nodes?page=0&limit=100&label=Person
	s.router.HandleFunc("/knowledgegraph/nodes", nodeController.GetNodes).Methods(http.MethodGet)
	// /knowledgegraph/nodes/{elementId}/relations
	s.router.HandleFunc("/knowledgegraph/nodes/{elementId}/relations", nodeController.GetRelatedNodes).Methods(http.MethodGet)

	relationshipController := &graph.RelationshipController{}
	s.router.HandleFunc("/knowledgegraph/relationships", relationshipController.Create).Methods(http.MethodPost)
	s.router.HandleFunc("/knowledgegraph/relationships/bulk", relationshipController.CreateBulk).Methods(http.MethodPost)
	s.router.HandleFunc("/knowledgegraph/relationships/{elementId}", relationshipController.Put).Methods(http.MethodPut)
	s.router.HandleFunc("/knowledgegraph/relationships/{elementId}", relationshipController.Delete).Methods(http.MethodDelete)

	s.router.HandleFunc("/knowledgegraph/relationships/{elementId}", relationshipController.GetRelationship).Methods(http.MethodGet)
	// /knowledgegraph/relationships?page=0&limit=100&type=depend
	s.router.HandleFunc("/knowledgegraph/relationships", relationshipController.GetRelationships).Methods(http.MethodGet)
}
