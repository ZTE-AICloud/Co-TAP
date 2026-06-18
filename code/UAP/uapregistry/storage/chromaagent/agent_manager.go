package agent

import (
	"context"
	"os"

	"github.com/amikos-tech/chroma-go/pkg/embeddings"
)

var (
	chromaAgentManager *ChromaAgentManager
	embeddingFunc embeddings.EmbeddingFunction
)

func InitChromaAgent() error {
	var err error
	chromaAgentManager, err = NewChromaAgentManager()
	if err != nil {
		return err
	}

	return initAgentCollection()
}

func initAgentCollection() error {
	var err error
	if os.Getenv("EMBEDDING_TYPE") == "ollama" {
		modelName := os.Getenv("MODEL_NAME")
		baseUrl := os.Getenv("BASE_URL")
		embeddingFunc, err = GetOllamaEmbeddingFunction(modelName, baseUrl)
		if err != nil {
			chromaAgentManager.logWriter.Errorf("failed to init embedding function, err: %v", err)
			return err
		}
	}

	err = chromaAgentManager.chromaClient.CreateCollection(context.Background(), "agents_registry", embeddingFunc, nil)
	if err != nil {
		chromaAgentManager.logWriter.Errorf("failed to init agent collection, err: %v", err)
		return err
	}
	return nil
}


func GetChromaAgentManager() *ChromaAgentManager {
	return chromaAgentManager
}