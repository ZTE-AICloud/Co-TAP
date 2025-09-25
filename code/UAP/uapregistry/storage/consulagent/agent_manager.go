package agent

import (
	"uapregistry/logger"
)

var (
	localAgent    *Agent
	localAgentCfg string
)

func InitLocalAgent(cfg *Config) error {
	localAgentCfg = cfg.GetConsulAgent()
	agent, err := Create(cfg)
	if err != nil {
		logger.GetLogger().Errorf("failed to create local agent:%v", err)
		return err
	}

	localAgent = agent
	return nil
}

var GetLocalAgent = func() *Agent {
	return localAgent
}

func GetLocalAgentCfg() string {
	return localAgentCfg
}
