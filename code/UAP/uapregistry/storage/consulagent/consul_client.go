package agent

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
	"uapregistry/logger"

	consulapi "github.com/hashicorp/consul/api"
)

var NewConsulClient = func(consulAgent string, maxConnectAttempts int) (*consulapi.Client, error) {
	var (
		client *consulapi.Client
		err    error
	)

	consulConfig := consulapi.DefaultConfig()
	consulConfig.Transport.MaxIdleConns = 2000
	consulConfig.Transport.MaxIdleConnsPerHost = 2000
	consulConfig.HttpClient = &http.Client{
		Transport: consulConfig.Transport,
		// Set blocking query timeout
		Timeout: 11 * time.Minute,
	}

	consulAgentURL, err := url.Parse(consulAgent)
	logger.GetLogger().Info("start to init consul-agent:", consulAgent)
	if err != nil {
		logger.GetLogger().Errorf("Failed to Parse %s:%v", consulAgent, err)
		return nil, err
	}

	if consulAgentURL.Host != "" {
		consulConfig.Address = consulAgentURL.Host
	}

	if consulAgentURL.Scheme != "" {
		consulConfig.Scheme = consulAgentURL.Scheme
	}

	client, err = consulapi.NewClient(consulConfig)
	if err != nil {
		logger.GetLogger().Errorf("Failed to NewClient:%v", err)
		return nil, err
	}

	for attempt := 1; attempt <= maxConnectAttempts; attempt++ {
		if _, err = client.Agent().Self(); err == nil {
			break
		}

		if attempt == maxConnectAttempts {
			break
		}

		logger.GetLogger().Infof("[Attempt: %d] Attempting access to Consul after 5 second sleep", attempt)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to Consul agent: %v, error: %v", consulAgent, err)
	}
	logger.GetLogger().Infof("Consul agent init success: %v", consulAgent)

	return client, nil
}
