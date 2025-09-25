package types

type Route struct {
	ID                      string    `json:"id"`
	CreatedAt               Timestamp `json:"created_at"`
	UpdatedAt               Timestamp `json:"updated_at"`
	Name                    string    `json:"name"`
	GatewayId               string    `json:"gateway_id"`
	Protocols               []string  `json:"protocols"`
	Methods                 []string  `json:"methods"`
	Hosts                   []string  `json:"hosts"`
	Paths                   []string  `json:"paths"`
	Headers                 []string  `json:"headers"`
	HTTPSRedirectStatusCode int       `json:"https_redirect_status_code"`
	RegexPriority           int       `json:"regex_priority"`
	StripPath               *bool     `json:"strip_path"`
	PreserveHost            *bool     `json:"preserve_host"`
	Request_buffering       *bool     `json:"request_buffering"`
	Response_buffering      *bool     `json:"response_buffering"`
	SNIs                    []string  `json:"snis"`
	Sources                 []string  `json:"sources"`
	Destinations            []string  `json:"destinations"`
	Tags                    []string  `json:"tags"`
	Service                 string    `json:"service"`
	AgentProtocol           string    `json:"agent_protocol"`
	Index                   uint64    `json:"index"`
}
