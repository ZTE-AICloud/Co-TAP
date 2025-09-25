package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"uapregistry/healthcheckmanager"
	"uapregistry/storage/consulagent/cache"

	"github.com/gorilla/mux"
)

// HTTP - PUT - /services/{id}/renewal
func (s *HTTPServer) RenewalHandler(w http.ResponseWriter, r *http.Request) {
	CheckUpdateSvc()
	defer UnlockUpdateSvc()

	id := mux.Vars(r)["id"]

	s.log.Infow("renewal a service", "id", id)
	if timeout := preDo(id); timeout {
		s.log.Errorf("(id:%s) Waiting for the completion of last operation is timed out", id)
		Response408(w)
		return
	}
	defer postDo(id)

	statusCode, err := healthcheckmanager.GetCache().UpdateServiceTTL(id)
	if err != nil {
		w.WriteHeader(statusCode)
		fmt.Fprintf(w, "Failed to update ttl")
		return
	}

	Response200(w, nil)
}

// HTTP - GET - /health
func (s *HTTPServer) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	Response200(w, []byte(`{"status":"healthy"}`))
}

func (s *HTTPServer) ServiceHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	s.log.Infof("do health check for service(id:%s)", id)

	svc := cache.GetServicesCache().GetServiceByID(id)
	if svc == nil {
		s.log.Infof("service(id:%s) not found", id)
		Response404(w, fmt.Sprintf("service(id:%s) not found", id))
		return
	}

	hcr := healthcheckmanager.DoServiceHealthCheckForOnce(svc)
	j, err := json.Marshal(hcr)
	if err != nil {
		s.log.Errorf("Failed to marshal json data when build response: %v", err)
		Response500(w, "Failed to marshal json data when build response")
		return
	}

	Response200(w, j)
}
