package http

import (
	"uapregistry/config"
)

type Config struct {
	BindIP   string
	HTTPIPs  []string
	HTTPPort string
}

func DefaultConfig() *Config {
	return &Config{
		BindIP:   config.GetHTTPBindIP(),
		HTTPIPs:  config.GetHTTPIPs(),
		HTTPPort: defaultHTTPListenPort(),
	}
}

func defaultHTTPListenPort() string {
	return config.GetHTTPListenPort()
}

func (c *Config) GetHTTPIPs() []string {
	return c.HTTPIPs
}

func (c *Config) GetHTTPBindIP() string {
	return c.BindIP
}

func (c *Config) GetHTTPPort() string {
	return c.HTTPPort
}

func (c *Config) SetHTTPPort(port string) {
	c.HTTPPort = port
}

func (c *Config) SetHTTPBindIP(bindIP string) {
	c.BindIP = bindIP
}
