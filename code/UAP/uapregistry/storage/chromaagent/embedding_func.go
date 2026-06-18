package agent

import (
	"fmt"

	"github.com/amikos-tech/chroma-go/pkg/embeddings"
	"github.com/amikos-tech/chroma-go/pkg/embeddings/ollama"
)

func GetOllamaEmbeddingFunction(modelName string, baseUrl string) (*ollama.OllamaEmbeddingFunction, error) {
	embeddingFunc, err := ollama.NewOllamaEmbeddingFunction(
		ollama.WithModel(embeddings.EmbeddingModel(modelName)),
		ollama.WithBaseURL(baseUrl),
	)
	if err != nil {
		return nil, fmt.Errorf("init ollama embedding function failed, err: %v", err)
	}

	return embeddingFunc, nil
}
