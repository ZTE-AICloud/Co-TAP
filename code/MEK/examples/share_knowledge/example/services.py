from __future__ import annotations

import os
from typing import Optional

from mek.storage.file_repository import FileMemoryRepository
from mek.memory.workflow_memory import WorkflowMemoryService
from mek.extract.extractor import ExtractService
from mek.knowledge.default import DefaultKnowledgeService
from mek.llm_client.openai_client import OpenAIClient
from mek.llm_client.config import LLMConfig
from mek.llm_client.embedding_client import OpenAIEmbedder, OpenAIEmbedderConfig
from mek.config.manager import ConfigManager

# 使用前向引用避免循环导入
from typing import TYPE_CHECKING
if TYPE_CHECKING:
    from mek.core.repository import MemoryRepository
    from mek.memory.base import BaseMemoryService
    from mek.extract.extractor import ExtractService
    from mek.knowledge.service import KnowledgeService


class ServicesFactory:
    """服务工厂 - 按照新架构创建服务实例
    
    依赖关系：
    - Knowledge Service 需要 Extract Service 和 Memory Service
    - Memory Service 需要 Repository 和 LLM
    - Extract Service 独立创建
    """
    
    def __init__(self) -> None:
        # 加载配置
        self.config = ConfigManager()
        
        # 创建基础服务
        self._repository: Optional['MemoryRepository'] = None
        self._llm_client: Optional[OpenAIClient] = None
        self._embedding_client: Optional[OpenAIEmbedder] = None
        self._extract_service: Optional['ExtractService'] = None
    
    def make_repository(self, agent_name: str = "default") -> 'MemoryRepository':
        """创建记忆仓库"""
        if not self._repository:
            mek_config = self.config.get_mek_config()
            storage_path = mek_config.memory_storage_path or 'data/memory'
            # 为每个agent创建独立的存储路径
            agent_storage_path = os.path.join(storage_path, agent_name)
            os.makedirs(agent_storage_path, exist_ok=True)
            
            self._repository = FileMemoryRepository(agent_storage_path)
        
        return self._repository
    
    def make_llm_client(self) -> OpenAIClient:
        """创建LLM客户端"""
        if not self._llm_client:
            llm_config_data = self.config.get_llm_config()
            llm_config = LLMConfig(
                api_key=llm_config_data.api_key,
                provider=llm_config_data.provider,
                model=llm_config_data.model_name,
                base_url=llm_config_data.base_url
            )
            self._llm_client = OpenAIClient(config=llm_config)
        
        return self._llm_client
    
    def make_embedding_client(self) -> OpenAIEmbedder:
        """创建Embedding客户端"""
        if not self._embedding_client:
            embedding_config_data = self.config.get_embedding_config()
            embedding_config = OpenAIEmbedderConfig(
                api_key=embedding_config_data.api_key,
                base_url=embedding_config_data.base_url,
                embedding_model=embedding_config_data.model_name,
                embedding_dim=embedding_config_data.embedding_dim
            )
            self._embedding_client = OpenAIEmbedder(config=embedding_config)
        
        return self._embedding_client
    
    def make_extract_service(self) -> 'ExtractService':
        """创建Extract服务"""
        if not self._extract_service:
            self._extract_service = ExtractService()
        
        return self._extract_service
    
    def make_memory_service(self, agent_name: str) -> 'BaseMemoryService':
        """创建Memory服务
        
        需要repository和llm_client
        """
        repository = self.make_repository(agent_name)
        llm_client = self.make_llm_client()
        
        # 目前只支持工作流记忆
        return WorkflowMemoryService(repository, llm_client)
    
    def make_knowledge_service(self, agent_name: str) -> 'KnowledgeService':
        """创建Knowledge服务
        
        需要extract_service和memory_service
        """
        extract_service = self.make_extract_service()
        memory_service = self.make_memory_service(agent_name)
        
        return DefaultKnowledgeService(extract_service, memory_service)
    
    def create_agent_services(self, agent_name: str) -> tuple['KnowledgeService', 'BaseMemoryService']:
        """为指定agent创建完整的服务实例
        
        Returns:
            (knowledge_service, memory_service) 元组
        """
        print(f"[ServicesFactory]: 为 {agent_name} 创建服务实例...")
        
        knowledge_service = self.make_knowledge_service(agent_name)
        memory_service = self.make_memory_service(agent_name)
        
        print(f"[ServicesFactory]: {agent_name} 服务创建完成")
        return knowledge_service, memory_service

