from pydantic import BaseModel, Field

from mek.llm_client.base import EmbedderClient

"""
Copyright 2024, Zep Software, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

from openai import OpenAI
from openai.types import EmbeddingModel

DEFAULT_EMBEDDING_MODEL = 'text-embedding-3-small'
EMBEDDING_DIM = 1024


class EmbedderConfig(BaseModel):
    embedding_dim: int = Field(default=EMBEDDING_DIM, frozen=True)


class OpenAIEmbedderConfig(EmbedderConfig):
    embedding_model: EmbeddingModel | str = DEFAULT_EMBEDDING_MODEL
    api_key: str | None = None
    base_url: str | None = None


class OpenAIEmbedder(EmbedderClient):
    """
    OpenAI Embedder Client

    This client supports both AsyncOpenAI and AsyncAzureOpenAI clients.
    """

    def __init__(
            self,
            config: OpenAIEmbedderConfig | None = None
    ):
        if config is None:
            config = OpenAIEmbedderConfig()
        self.config = config
        self.client = OpenAI(api_key=config.api_key, base_url=config.base_url)

    def create(self, input_data: str | list[str]) -> list[float]:
        result = self.client.embeddings.create(
            input=input_data, model=self.config.embedding_model
        )
        return result.data[0].embedding[: self.config.embedding_dim]

    def create_batch(self, input_data_list: list[str]) -> list[list[float]]:
        result = self.client.embeddings.create(
            input=input_data_list, model=self.config.embedding_model, encoding_format="float"
        )
        return [embedding.embedding[: self.config.embedding_dim] for embedding in result.data]
