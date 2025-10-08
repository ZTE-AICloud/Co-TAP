from __future__ import annotations

from abc import ABC, abstractmethod
from typing import Dict, List, Optional

from ..core.models import KnowledgeItem


class KnowledgeService(ABC):
    """知识服务抽象基类
    
    负责知识的共享和吸收，是MEK协议的核心接口
    """
    
    @abstractmethod
    def share(self, query: str, context: Optional[Dict[str, object]] = None) -> List[KnowledgeItem]:
        """共享知识
        
        根据查询从记忆中提取相关知识并进行标准化处理
        
        Args:
            query: 查询字符串
            context: 上下文信息
            
        Returns:
            标准化的知识项列表
        """
        pass
    
    @abstractmethod
    def absorb(self, knowledge_items: List[KnowledgeItem]) -> List[str]:
        """吸收知识
        
        将外部知识转换为本地记忆，包含冲突检测和解决
        
        Args:
            knowledge_items: 外部知识项列表
            
        Returns:
            新生成或更新的记忆UUID列表
        """
        pass 