from __future__ import annotations

from abc import ABC, abstractmethod
from typing import Any, Dict, List, Optional, Union

from ..core.models import MemoryItem
from ..core.repository import MemoryRepository
from ..core.types import MemoryType
from ..core.retrieval import RetrievalType, RetrievalOptions


class BaseMemoryService(ABC):
    """记忆服务抽象基类
    
    负责记忆的CRUD操作，包括从原始数据生成记忆、检索记忆、更新记忆等
    """
    
    def __init__(self, repository: MemoryRepository, memory_type: MemoryType) -> None:
        self.repository = repository
        self.memory_type = memory_type
    
    @abstractmethod
    def add(self, input_data: Union[str, Dict[str, Any]], **kwargs) -> List[str]:
        """从原始数据中提取生成记忆
        
        使用LLM从输入数据中提取通用工作流，生成记忆，并存入数据库
        包含与记忆库中旧记忆条目的冲突解决过程
        
        Args:
            input_data: 原始输入数据（轨迹数据等）
            
        Returns:
            新记忆的UUID列表（或需要更新的记忆UUID）
        """
        pass
    
    @abstractmethod
    def retrieve(
        self,
        query: str,
        retrieval_type: RetrievalType = RetrievalType.SEMANTIC,
        options: Optional[RetrievalOptions] = None
    ) -> List[MemoryItem]:
        """检索相关记忆
        
        Args:
            query: 查询字符串
            retrieval_type: 检索类型
            options: 检索选项
            
        Returns:
            相关记忆列表
        """
        pass
    
    @abstractmethod
    def update(self, memory_item: MemoryItem) -> str:
        """更新或添加记忆条目
        
        如果UUID是新的，则添加记忆对象
        如果UUID已存在，则更新对应的旧记录
        
        Args:
            memory_item: 记忆对象
            
        Returns:
            记忆UUID
        """
        pass
    
    @abstractmethod
    def delete(self, memory_id: str) -> bool:
        """删除指定记忆
        
        Args:
            memory_id: 记忆ID
            
        Returns:
            是否删除成功
        """
        pass 