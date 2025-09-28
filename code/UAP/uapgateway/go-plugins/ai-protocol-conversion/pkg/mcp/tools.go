package mcp

import (
	http2 "ai-protocol-conversion/pkg/http"
	"context"
	"encoding/json"
	apiclient "github.com/ZTE-AICloud/Co-TAP/code/UAP/uapregistrysdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"net/url"
)

const (
	CardInfoUrl = "/.well-known/agent-card.json"
)

func ListAgent(ctx context.Context, req *mcp.CallToolRequest, param map[string]any) (*mcp.CallToolResult, any, error) {
	callResult := &mcp.CallToolResult{}
	if agentInfo, ok := ctx.Value("agentInfo").(*apiclient.Service); !ok {
		callResult.Content = []mcp.Content{&mcp.TextContent{Text: "fail to get agent info"}}
		callResult.IsError = true
	} else {
		if agentInfo.HasAgentInfo() {
			_info, err := json.Marshal(agentInfo.AgentInfo)
			if err != nil {
				return nil, nil, err
			}
			callResult.Content = []mcp.Content{&mcp.TextContent{Text: string(_info)}}
		} else {
			agentUrl, err := url.JoinPath(*agentInfo.AgentInfoUrl, CardInfoUrl)
			if err != nil {
				return nil, nil, err
			}

			reqParam := &http2.ReqParams{
				Url:    agentUrl,
				Method: "GET",
			}
			result, err := reqParam.Call()
			if err != nil {
				return nil, nil, err
			}
			callResult.Content = result
		}
	}
	return callResult, nil, nil
}

//func RunAgentTask(ctx context.Context, req *mcp.CallToolRequest, param map[string]any) (*mcp.CallToolResult, any, error) {
//	callResult := &mcp.CallToolResult{}
//	if agentInfo, ok := ctx.Value("agentInfo").(*apiclient.Service); !ok {
//		callResult.Content = []mcp.Content{&mcp.TextContent{Text: "fail to get agent info"}}
//		callResult.IsError = true
//	} else {
//		if agentInfo.HasAgentInfo() {
//			_info, err := json.Marshal(agentInfo.AgentInfo)
//			if err != nil {
//				return nil, nil, err
//			}
//			callResult.Content = []mcp.Content{&mcp.TextContent{Text: string(_info)}}
//		} else {
//
//		}
//	}
//
//	return callResult, nil, nil
//}
