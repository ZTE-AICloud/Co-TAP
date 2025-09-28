package main

import (
	a2a "ai-protocol-conversion/pkg/mcp"
	"bytes"
	"context"
	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	apiclient "github.com/ZTE-AICloud/Co-TAP/code/UAP/uapregistrysdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
)

const (
	Version  = "0.2"
	Priority = 100

	NamingServerAddrEnvKey = "NAMING-SERVER-ADDR"
	DefaultErrorMessage    = `{"message": "internal error"}`
	DefaultNaHost          = "127.0.0.1"
	DefaultNaSchema        = "http"
)

var NaHost = DefaultNaHost
var NaSchema = DefaultNaSchema
var defaultApiClient *apiclient.APIClient

type Config struct {
	Name    string
	SrvName string
}

func init() {
	NaAddr := os.Getenv(NamingServerAddrEnvKey)
	if len(NaAddr) > 0 {
		naUrl, err := url.Parse(NaAddr)
		if err != nil {
			return
		}

		NaSchema = naUrl.Scheme
		NaHost = naUrl.Host
	}

	configuration := apiclient.NewConfiguration()
	configuration.Host = NaHost
	configuration.Scheme = NaSchema
	defaultApiClient = apiclient.NewAPIClient(configuration)
}

func main() {
	err := server.StartServer(New, Version, Priority)
	if err != nil {
		return
	}
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(kong *pdk.PDK) {
	defaultContext := context.Background()
	if defaultApiClient == nil {
		_ = kong.Log.Info("registration center address is empty, not support ai-protocol-conversion")
		return
	}

	//TODO watch service update in background
	srv, _, err := defaultApiClient.DefaultAPI.ServicesServiceNameGet(defaultContext, conf.SrvName).Execute()
	if err != nil || len(srv) == 0 {
		_ = kong.Log.Err("fail to get service info: %s", err.Error())
		kong.Response.Exit(500, []byte(`{"message": "fail to get service info"}`), nil)
		return
	}

	agentPro := srv[0].AgentProtocol
	//only support mcp-to-a2a
	if agentPro == nil || *agentPro != "a2a" {
		// do nothing
		return
	}

	mcpServer := mcp.NewServer(&mcp.Implementation{Name: conf.Name}, nil)
	mcp.AddTool(mcpServer, &mcp.Tool{Name: "list_agents"}, a2a.ListAgent)
	//mcp.AddTool(mcpServer, &mcp.Tool{Name: "run_agent_task"}, mcp.RunAgentTask)

	h := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return mcpServer
	}, &mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true})

	method, err := kong.Request.GetMethod()
	if err != nil {
		_ = kong.Log.Err("fail to get method: %s", err.Error())
		kong.Response.Exit(500, []byte(DefaultErrorMessage), nil)
		return
	}
	path, err := kong.Request.GetPath()
	if err != nil {
		_ = kong.Log.Err("fail to get path: %s", err.Error())
		kong.Response.Exit(500, []byte(DefaultErrorMessage), nil)
		return
	}
	body, err := kong.Request.GetRawBody()
	if err != nil {
		_ = kong.Log.Err("fail to get body: %s", err.Error())
		kong.Response.Exit(500, []byte(DefaultErrorMessage), nil)
		return
	}
	valueCtx := context.WithValue(defaultContext, "agentInfo", &srv[0])
	wrapReq, err := http.NewRequestWithContext(valueCtx, method, path,
		bytes.NewReader(body))
	wrapReq.Header, err = kong.Request.GetHeaders(100)
	wrapResp := httptest.NewRecorder()
	h.ServeHTTP(wrapResp, wrapReq)

	kong.Response.Exit(wrapResp.Code, wrapResp.Body.Bytes(), wrapResp.Header())
}
