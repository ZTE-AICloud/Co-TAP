package utils

import (
	"encoding/json"
	"hash/fnv"
	"net"
	"net/http"
	"strings"

	"uapregistry/logger"
	"uapregistry/types"
)

// getRemoteIP
func GetRemoteIP(r *http.Request) string {
	remoteIP := r.Header.Get("X-Forwarded-For")
	if remoteIP != "" {
		return strings.Split(remoteIP, ",")[0]
	}
	remoteIP = r.Header.Get("X-Real-IP")
	if remoteIP != "" {
		return remoteIP
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		logger.GetLogger().Warnf("failed to get ip from RemoteAddr:%v,ignore it", err)
		return ""
	}
	return ip
}

func CalculateServicesIndex(svcs []*types.Service) uint64 {
	var (
		index  uint64
		indexs = make(map[string]uint64)
		h      = fnv.New32a()
	)

	for _, svc := range svcs {
		indexs[svc.ID] = svc.Index
	}

	data, _ := json.Marshal(indexs)
	h.Write(data)
	index = uint64(h.Sum32())
	return index
}

func CalculateServiceMapIndex(svcs map[string]*types.Service) uint64 {
	var (
		index  uint64
		indexs = make(map[string]uint64)
		h      = fnv.New32a()
	)

	for _, svc := range svcs {
		indexs[svc.ID] = svc.Index
	}

	data, _ := json.Marshal(indexs)
	h.Write(data)
	index = uint64(h.Sum32())
	return index
}

func GetDiffServiceNames(old, new map[string]*types.Service) (add, update, del []string) {
	for id, oldSvc := range old {
		newSvc, ok := new[id]
		if !ok {
			del = append(del, oldSvc.Name)
		} else {
			if oldSvc.Index != newSvc.Index {
				update = append(update, oldSvc.Name)
			}
		}
	}
	for id, newSvc := range new {
		if _, ok := old[id]; !ok {
			add = append(add, newSvc.Name)
		}
	}
	return
}
