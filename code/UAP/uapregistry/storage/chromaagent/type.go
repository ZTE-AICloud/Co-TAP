package agent

import (
	"sync"
	"uapregistry/logger"

	"github.com/a2aproject/a2a-go/a2a"
	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"
)

type ChromaAgentManager struct {
	chromaClient *ChromaClient
	logWriter    logger.Logger
}

type ChromaClient struct {
	sync.RWMutex // 用于保护 collection map 的并发安全
	client       chroma.Client
	collection   map[string]chroma.Collection
}

type ChromaQueryRequest struct {
	CollectionName string
	Query          string
	Top_K          int
	Threshold      float32 // 最小相似度阈值，低于该值的结果直接过滤

	// 新增：硬约束过滤参数
	RequiredInputModes   []string              `json:"required_input_modes,omitempty"`
	RequiredOutputModes  []string              `json:"required_output_modes,omitempty"`
	RequiredCapabilities a2a.AgentCapabilities `json:"required_capabilities,omitempty"`
	ProviderOrganization string                `json:"provider_organization,omitempty"`
}

// ChromaUpsertRequest 幂等写入请求参数
type ChromaUpsertRequest struct {
    CollectionName string                   `json:"collection_name"`
    IDs            []string                 `json:"ids"`
    Documents      []string                 `json:"documents"`
    Metadatas      []map[string]interface{} `json:"metadatas,omitempty"`
}