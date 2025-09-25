package types

import (
	"encoding/json"
	"time"
)

// 自定义时间类型
type Timestamp time.Time

// 实现 MarshalJSON 接口，返回毫秒级时间戳
func (t Timestamp) MarshalJSON() ([]byte, error) {
	millis := time.Time(t).UnixNano() / int64(time.Millisecond)
	return json.Marshal(millis)
}

// 实现 UnmarshalJSON 接口（如果需要反序列化）
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var millis int64
	if err := json.Unmarshal(data, &millis); err != nil {
		return err
	}
	*t = Timestamp(time.Unix(0, millis*int64(time.Millisecond)))
	return nil
}

type Service struct {
	ID              string               `json:"id"`
	Ephemeral       bool                 `json:"ephemeral"`
	CreatedAt       Timestamp            `json:"created_at"`
	UpdatedAt       Timestamp            `json:"updated_at"`
	Name            string               `json:"name"`
	Retries         int                  `json:"retries"`
	Protocol        string               `json:"protocol"`
	Host            string               `json:"host"`
	Port            int                  `json:"port"`
	Path            string               `json:"path"`
	ConnectTimeout  uint64               `json:"connect_timeout"`
	WriteTimeout    uint64               `json:"write_timeout"`
	ReadTimeout     uint64               `json:"read_timeout"`
	Tags            []string             `json:"tags,omitempty"`
	AgentProtocol   string               `json:"agent_protocol"`
	AgentInfo       *AgentInfo           `json:"agent_info"`
	AgentInfoUrl    string               `json:"agent_info_url"`
	Index           uint64               `json:"index"`
	HealthStatus    string               `json:"health_status"`
	EphemeralCheck  *EphemeralCheckInfo  `json:"ephemeral_check,omitempty"`
	PersistentCheck *PersistentCheckInfo `json:"persistent_check,omitempty"`
}

type AgentInfo struct {
	A2AAgentCard     map[string]interface{} `json:"a2a_agent_card,omitempty"`
	McpServer        map[string]interface{} `json:"mcp_server_info,omitempty"`
	AcpAgentManifest map[string]interface{} `json:"acp_agent_manifest,omitempty"`
}

type EphemeralCheckInfo struct {
	CheckType               string `json:"check_type"`
	RenewalDeleteTimeout    string `json:"renewal_delete_timeout"`
	RenewalInterval         string `json:"renewal_interval"`
	RenewalUnhealthyTimeout string `json:"renewal_unhealthy_timeout"`
}

type PersistentCheckInfo struct {
	CheckType             string `json:"check_type"`
	CheckInterval         string `json:"check_interval"`
	CheckTimeout          string `json:"check_timeout"`
	CheckHTTPURL          string `json:"check_http_url"`
	CheckHTTPMethod       string `json:"check_http_method"`
	CheckUnhealthyTimeout string `json:"check_unhealthy_timeout"`
}
