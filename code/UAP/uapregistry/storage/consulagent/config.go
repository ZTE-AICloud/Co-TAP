package agent

import (
	"strings"
	"uapregistry/config"
)

type Config struct {
	ConsulAgent string
}

func DefaultConfig() *Config {
	return &Config{
		ConsulAgent: defaultConsulAgentConfig(),
	}
}

func defaultConsulAgentConfig() string {
	consulIP := config.GetConsulIP()

	// ip_v6
	if strings.Contains(consulIP, ":") {
		return "http://[" + consulIP + "]:" + config.GetConsulHTTPPort()
	}

	return "http://" + consulIP + ":" + config.GetConsulHTTPPort()
}

func (c *Config) GetConsulAgent() string {
	return c.ConsulAgent
}

func (c *Config) SetConsulAgent(ca string) {
	c.ConsulAgent = ca
}
