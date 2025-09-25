package utils

import (
	"encoding/json"
	"errors"
	"hash/fnv"
	"slices"
	"strconv"
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"

	"uapregistry/logger"
	"uapregistry/types"
)

func BuildRouteKey(name string) string {
	return "v1/routes/" + name
}

func FillRouteDefaultValue(route *types.Route) {
	route.ID = BuildRouteUID(route.Name)
	route.CreatedAt = types.Timestamp(time.Now())
	route.UpdatedAt = types.Timestamp(time.Now())
	for i, protocol := range route.Protocols {
		route.Protocols[i] = strings.ToLower(protocol)
	}
	if route.HTTPSRedirectStatusCode == 0 {
		route.HTTPSRedirectStatusCode = 426
	}
	if route.StripPath == nil {
		v := true
		route.StripPath = &v
	}
	if route.PreserveHost == nil {
		v := false
		route.PreserveHost = &v
	}
	if route.Request_buffering == nil {
		v := true
		route.Request_buffering = &v
	}
	if route.Response_buffering == nil {
		v := true
		route.Response_buffering = &v
	}
}

func CheckAndModifyRoute(route *types.Route) error {
	if route.Service == "" {
		return errors.New("service is empty")
	}

	if err := checkRouteProtocol(route.Protocols); err != nil {
		return err
	}
	if err := checkHTTPSRedirectStatusCode(route.HTTPSRedirectStatusCode); err != nil {
		return err
	}
	if err := checkAgentProtocol(route.AgentProtocol); err != nil {
		return err
	}

	return nil
}

func checkRouteProtocol(protocols []string) error {
	if len(protocols) == 0 {
		return errors.New("protocols is empty")
	}

	validProtocols := []string{"http", "https", "tcp", "udp", "tls", "grpc", "grpcs", "tls_passthrough"}
	for _, protocol := range protocols {
		if !slices.Contains(validProtocols, protocol) {
			return errors.New("invalid protocol: " + protocol)
		}
	}

	return nil
}

func checkHTTPSRedirectStatusCode(code int) error {
	validCodes := []int{426, 301, 302, 307, 308}
	if !slices.Contains(validCodes, code) {
		return errors.New("invalid https_redirect_status_code: " + strconv.Itoa(code))
	}
	return nil
}

func ConsulKVPairs2Routes(kvps []*consulapi.KVPair) []*types.Route {
	var routes []*types.Route

	for _, kvp := range kvps {
		route := ConsulKVPair2Route(kvp)
		if route != nil {
			route.Index = kvp.ModifyIndex
			routes = append(routes, route)
		}
	}

	return routes
}

func ConsulKVPair2Route(kvp *consulapi.KVPair) *types.Route {
	if kvp == nil {
		return nil
	}

	var route *types.Route
	err := json.Unmarshal(kvp.Value, &route)
	if err != nil {
		logger.GetLogger().Errorf("Unmarshal consul kv pair failed: %v", err)
		return nil
	}

	return route
}

func CalculateRoutesIndex(routes []*types.Route) uint64 {
	var (
		index  uint64
		indexs = make(map[string]uint64)
		h      = fnv.New32a()
	)

	for _, route := range routes {
		indexs[route.ID] = route.Index
	}

	data, _ := json.Marshal(indexs)
	h.Write(data)
	index = uint64(h.Sum32())
	return index
}

func GetDiffRoutes(routesOld, routesNew []*types.Route) (added, updated, deleted []string) {
	oldMap := make(map[string]*types.Route)
	newMap := make(map[string]*types.Route)

	for _, r := range routesOld {
		oldMap[r.ID] = r
	}

	for _, r := range routesNew {
		newMap[r.ID] = r
	}

	for id, newRoute := range newMap {
		if oldRoute, exists := oldMap[id]; !exists {
			added = append(added, newRoute.Name)
		} else {
			if oldRoute.Index != newRoute.Index {
				updated = append(updated, newRoute.Name)
			}
		}
	}

	for id, oldRoute := range oldMap {
		if _, exists := newMap[id]; !exists {
			deleted = append(deleted, oldRoute.Name)
		}
	}

	return added, updated, deleted
}
