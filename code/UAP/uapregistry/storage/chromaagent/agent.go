package agent

import (
	"context"
	"fmt"
	"time"
	"uapregistry/logger"
	"uapregistry/types"

	"github.com/a2aproject/a2a-go/a2a"
)

const (
	timeout            = 30
	maxConnectAttempts = 12
	MetaProviderOrg    = "provider_organization"
	MetaCapStreaming   = "cap_streaming"
	MetaCapPushNotify  = "cap_push_notifications"
	MetaAllInputModes  = "all_input_modes"
	MetaAllOutputModes = "all_output_modes"
)

func NewChromaAgentManager() (*ChromaAgentManager, error) {
	var (
		chromaClient *ChromaClient
		err          error
	)
	for attempt := 1; attempt <= maxConnectAttempts; attempt++ {
		chromaClient, err = NewChromaClient(timeout * time.Second)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
	}
	if err != nil {
		return nil, fmt.Errorf("fail to create chroma agent manager, err: %v", err)
	}
	return &ChromaAgentManager{chromaClient: chromaClient, logWriter: logger.GetLogger()}, nil
}

// RegisterAgent 注册单个 Agent
func (c *ChromaAgentManager) RegisterAgent(agentID string, description string, agentCard *a2a.AgentCard) error {
	err := c.chromaClient.UpsertSingleDocument(
        context.Background(),
        "agents_registry",
        agentID,
        description,
        BuildAgentMetadata(agentCard),
    )
	if err != nil {
		return fmt.Errorf("register agent %s failed: %w", agentID, err)
	}

	c.logWriter.Infof("Agent %s registered successfully in collection %s", agentID, "agents_registry")
	return nil
}

// AgentSearchResult 携带相似度的搜索结果
type AgentSearchResult struct {
	AgentID    string
	Similarity float32
}

// QueryTopNAgents 查询最匹配的前 N 个 Agent ID（返回带相似度）
func (c *ChromaAgentManager) QueryTopNAgents(ssq *types.SemanticSearchRequest) ([]AgentSearchResult, error) {
	cqr := &ChromaQueryRequest{
		CollectionName:       "agents_registry",
		Query:                ssq.Query,
		Top_K:                ssq.TopK,
		Threshold:            ssq.MinRelevanceScore, // 传递阈值给底层查询函数统一过滤
		RequiredInputModes:   ssq.RequiredInputModes,
		RequiredOutputModes:  ssq.RequiredOutputModes,
		RequiredCapabilities: ssq.Capabilities,
		ProviderOrganization: ssq.ProviderOrganization,
	}
	res, err := c.chromaClient.Query(context.Background(), cqr)
	if err != nil {
		return nil, fmt.Errorf("query top N agents failed: %w", err)
	}

	idGroups := res.GetIDGroups()
	if len(idGroups) == 0 || len(idGroups[0]) == 0 {
		return []AgentSearchResult{}, nil
	}
	distanceGroups := res.GetDistancesGroups()

	// 校验结果完整性，避免索引越界
	if len(idGroups) != len(distanceGroups) || len(idGroups[0]) != len(distanceGroups[0]) {
		return []AgentSearchResult{}, nil
	}

	// 使用过滤后的距离计算相似度
	ids := idGroups[0]
	distances := distanceGroups[0]
	results := make([]AgentSearchResult, len(ids))
	for i := range ids {
		results[i] = AgentSearchResult{
			AgentID:    string(ids[i]),
			Similarity: 1 - float32(distances[i])/2,
		}
	}

	c.logWriter.Debugf("Found %d matching agents for query: %s", len(results), ssq.Query)
	return results, nil
}

// DeleteAgent 根据 agentID 删除
func (c *ChromaAgentManager) DeleteAgent(agentID string) error {
	return c.chromaClient.DeleteDocuments(context.Background(), "agents_registry", []string{agentID})
}

// BuildAgentMetadata 从 AgentCard 构建符合过滤规范的 metadata
// 预计算全局+所有技能的模式并集，保证查询时 $contains 过滤生效
func BuildAgentMetadata(agentCard *a2a.AgentCard) map[string]interface{} {
	meta := make(map[string]interface{})

	// 1. 提供商组织
	if agentCard.Provider != nil {
		meta[MetaProviderOrg] = agentCard.Provider.Org
	}

	// 2. 能力标志
	meta[MetaCapStreaming] = agentCard.Capabilities.Streaming
	meta[MetaCapPushNotify] = agentCard.Capabilities.PushNotifications

	// 3. 输入模式并集：全局默认 + 所有技能的输入模式
	inputSet := make(map[string]struct{})
	for _, mode := range agentCard.DefaultInputModes {
		inputSet[mode] = struct{}{}
	}
	for _, skill := range agentCard.Skills {
		for _, mode := range skill.InputModes {
			inputSet[mode] = struct{}{}
		}
	}
	inputModes := make([]string, 0, len(inputSet))
	for mode := range inputSet {
		inputModes = append(inputModes, mode)
	}
	meta[MetaAllInputModes] = inputModes

	// 4. 输出模式并集：全局默认 + 所有技能的输出模式
	outputSet := make(map[string]struct{})
	for _, mode := range agentCard.DefaultOutputModes {
		outputSet[mode] = struct{}{}
	}
	for _, skill := range agentCard.Skills {
		for _, mode := range skill.OutputModes {
			outputSet[mode] = struct{}{}
		}
	}
	outputModes := make([]string, 0, len(outputSet))
	for mode := range outputSet {
		outputModes = append(outputModes, mode)
	}
	meta[MetaAllOutputModes] = outputModes

	return meta
}