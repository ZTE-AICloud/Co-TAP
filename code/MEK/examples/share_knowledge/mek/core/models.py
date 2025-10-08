from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from typing import Any, Dict, List, Optional

from .types import MemoryType, ImportanceLevel


@dataclass
class MemoryItem:
    memory_id: str
    memory_type: MemoryType
    created_at: datetime
    content: Dict[str, Any]
    importance: ImportanceLevel = ImportanceLevel.MEDIUM
    related_memory_ids: List[str] = field(default_factory=list)

    def summarize(self) -> str:
        """Return a brief summary of the memory content.

        Protocol requires returning summaries in extraction results. Concrete
        implementations decide how to summarize different memory types.
        """
        raise NotImplementedError


@dataclass
class KnowledgeItem:
    """知识项 - MEK协议中知识共享的标准数据结构"""
    
    knowledge_id: str                           # 知识唯一标识
    title: str                                  # 知识标题
    content: Dict[str, Any]                     # 知识内容
    source_memory_type: MemoryType              # 来源记忆类型
    source_agent: Optional[str] = None          # 源代理标识
    created_at: datetime = field(default_factory=datetime.utcnow)  # 创建时间
    metadata: Dict[str, Any] = field(default_factory=dict)  # 扩展元数据
    
    def __post_init__(self):
        """后初始化处理"""
        if not self.knowledge_id:
            # 如果没有提供ID，基于内容生成一个
            import hashlib
            content_str = f"{self.title}-{self.source_memory_type}-{self.created_at.isoformat()}"
            self.knowledge_id = f"knowledge-{hashlib.md5(content_str.encode()).hexdigest()[:8]}"
    
    def to_dict(self) -> Dict[str, Any]:
        """转换为字典格式，用于序列化"""
        return {
            "knowledge_id": self.knowledge_id,
            "title": self.title,
            "content": self.content,
            "source_memory_type": self.source_memory_type.value,
            "source_agent": self.source_agent,
            "created_at": self.created_at.isoformat(),
            "metadata": self.metadata
        }
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "KnowledgeItem":
        """从字典创建KnowledgeItem实例"""
        return cls(
            knowledge_id=data["knowledge_id"],
            title=data["title"],
            content=data["content"],
            source_memory_type=MemoryType(data["source_memory_type"]),
            source_agent=data.get("source_agent"),
            created_at=datetime.fromisoformat(data["created_at"]),
            metadata=data.get("metadata", {})
        ) 