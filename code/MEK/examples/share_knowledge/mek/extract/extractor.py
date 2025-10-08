from __future__ import annotations

import hashlib
import time
from datetime import datetime
from typing import Any, Dict, List, Optional

from ..core.models import MemoryItem, KnowledgeItem
from ..core.types import MemoryType, ImportanceLevel


class ExtractService:
    """Extract服务
    
    专门负责记忆和知识之间的转换，包含过滤、脱敏、泛化、标准化等流程
    """

    def __init__(self) -> None:
        pass

    def knowledge_to_memory(self, knowledge_item: KnowledgeItem) -> List[MemoryItem]:
        """将知识项转换为记忆项
        
        将外部KnowledgeItem转换为本地MemoryItem，包含感知和理解阶段
        
        Args:
            knowledge_item: 外部知识项
            
        Returns:
            转换后的记忆项
        """
        print(f"[ExtractService]: 开始知识到记忆转换...")

        knowledge = self._perception(knowledge_item)

        memory_items = self._understand(knowledge)
        return memory_items

    def memory_to_knowledge(
            self,
            memories: List[MemoryItem],
            context: Optional[Dict[str, Any]] = None
    ) -> List[KnowledgeItem]:
        """将记忆转换为知识
        
        包含过滤、脱敏、泛化、标准化等流程
        
        Args:
            memories: 记忆列表
            context: 上下文信息
            
        Returns:
            知识项列表
        """
        print(f"[ExtractService]: 开始记忆到知识转换...")

        # 过滤阶段：筛选可共享的记忆
        filtered_memories = self._filter_shareable_memories(memories)

        # 脱敏阶段：移除敏感信息
        desensitized_memories = self._desensitize_memories(filtered_memories)

        # 泛化阶段：去除个人上下文
        generalized_memories = self._generalize_memories(desensitized_memories)

        # 标准化阶段：转换为KnowledgeItem
        knowledge_items = self._standardize_to_knowledge(generalized_memories, context)

        print(f"[ExtractService]: 记忆到知识转换完成，生成 {len(knowledge_items)} 个知识项")
        return knowledge_items

    def _generate_memory_id(self, prefix: str = "memory") -> str:
        """生成记忆ID"""
        timestamp = str(int(time.time() * 1000))
        content_hash = hashlib.md5(timestamp.encode()).hexdigest()[:8]
        return f"{prefix}-{content_hash}"

    def _filter_shareable_memories(self, memories: List[MemoryItem]) -> List[MemoryItem]:
        """过滤阶段：筛选可共享的记忆"""
        filtered = []
        for memory in memories:
            # 基本条件：不是低重要性记忆，且是成功的记忆
            if memory.importance.value != 'low':
                if memory.memory_type == MemoryType.WORKFLOW:
                    result_status = memory.content.get('result_status', 'unknown')
                    if result_status in ['success', 'completed']:
                        filtered.append(memory)
                else:
                    filtered.append(memory)

        print(f"[ExtractService]: 过滤阶段 - 从 {len(memories)} 条记忆中筛选出 {len(filtered)} 条可共享记忆")
        return filtered

    def _desensitize_memories(self, memories: List[MemoryItem]) -> List[MemoryItem]:
        """脱敏阶段：移除敏感信息"""
        desensitized = []

        for memory in memories:
            # 创建副本进行处理
            memory_copy = MemoryItem(
                memory_id=memory.memory_id,
                memory_type=memory.memory_type,
                created_at=memory.created_at,
                content=memory.content.copy(),
                importance=memory.importance,
                related_memory_ids=memory.related_memory_ids.copy()
            )

            # 移除敏感信息
            self._remove_sensitive_info(memory_copy)
            desensitized.append(memory_copy)

        print(f"[ExtractService]: 脱敏阶段 - 处理了 {len(desensitized)} 条记忆的脱敏")
        return desensitized

    def _generalize_memories(self, memories: List[MemoryItem]) -> List[MemoryItem]:
        """泛化阶段：去除个人上下文"""
        generalized = []

        for memory in memories:
            # 创建副本进行处理
            memory_copy = MemoryItem(
                memory_id=memory.memory_id,
                memory_type=memory.memory_type,
                created_at=memory.created_at,
                content=memory.content.copy(),
                importance=memory.importance,
                related_memory_ids=memory.related_memory_ids.copy()
            )

            # 去除个人上下文
            self._remove_personal_context(memory_copy)
            generalized.append(memory_copy)

        print(f"[ExtractService]: 泛化阶段 - 处理了 {len(generalized)} 条记忆的泛化")
        return generalized

    def _standardize_to_knowledge(
            self,
            memories: List[MemoryItem],
            context: Optional[Dict[str, Any]] = None
    ) -> List[KnowledgeItem]:
        """标准化阶段：转换为KnowledgeItem"""
        knowledge_items = []
        source_agent = context.get('source_agent', 'unknown') if context else 'unknown'

        for memory in memories:
            # 标准化内容
            standardized_content = {
                "workflow_name": memory.content.get('workflow_name', memory.content.get('task_description', '未知工作流')),
                "steps": memory.content.get('steps', []),
                "result_status": memory.content.get('result_status', 'unknown'),
                "domain": memory.content.get('domain', '通用'),
            }

            # 生成知识ID
            knowledge_id = f"knowledge-{memory.memory_id}"

            # 创建KnowledgeItem
            knowledge_item = KnowledgeItem(
                knowledge_id=knowledge_id,
                title=f"工作流-{standardized_content['workflow_name']}",
                content=standardized_content,
                source_memory_type=memory.memory_type,
                source_agent=source_agent,
                created_at=datetime.utcnow(),
                metadata={
                    'original_memory_id': memory.memory_id,
                    'extraction_method': 'extract_service'
                }
            )

            knowledge_items.append(knowledge_item)

        print(f"[ExtractService]: 标准化阶段 - 生成了 {len(knowledge_items)} 个知识项")
        return knowledge_items

    def _remove_sensitive_info(self, memory: MemoryItem) -> None:
        """移除敏感信息"""
        sensitive_fields = [
            'user_id', 'session_id', 'ip_address', 'device_id',
            'personal_info', 'private_data', 'credentials'
        ]

        for field in sensitive_fields:
            if field in memory.content:
                del memory.content[field]

        # 清理嵌套的source_info中的敏感信息
        if 'source_info' in memory.content and isinstance(memory.content['source_info'], dict):
            source_info = memory.content['source_info']
            if 'session_id' in source_info:
                source_info['session_id'] = 'anonymized'

    def _remove_personal_context(self, memory: MemoryItem) -> None:
        """去除个人上下文"""
        # 移除个人化的元数据
        if 'metadata' in memory.content:
            del memory.content['metadata']

        # 移除特定的用户信息
        personal_fields = ['user_name', 'user_profile', 'personal_preferences']
        for field in personal_fields:
            if field in memory.content:
                del memory.content[field]

    def _perception(self, knowledge_item: KnowledgeItem):
        return knowledge_item

    def _understand(self, knowledge: KnowledgeItem):
        # 生成新的记忆ID
        memory_id = self._generate_memory_id("absorbed-记忆")

        # 当前只支持工作流类型，知识理解为工作流记忆
        knowledge_content = knowledge.content
        memory_content = {
            'task_description': knowledge.title,
            'workflow_name': knowledge_content.get('workflow_name', knowledge.title),
            'steps': knowledge_content.get('steps', ''),
            'result_status': knowledge_content.get('result_status', 'absorbed'),
            'domain': knowledge_content.get('domain', '通用'),
            'source_info': {
                'source_type': 'absorbed_knowledge',
                'source_agent': knowledge.source_agent,
                'original_knowledge_id': knowledge.knowledge_id
            }
        }

        memory_item = MemoryItem(
            memory_id=memory_id,
            memory_type=MemoryType.WORKFLOW,  
            created_at=datetime.utcnow(),
            content=memory_content,
            importance=ImportanceLevel.MEDIUM,
            related_memory_ids=[]
        )

        print(f"[ExtractService]: 知识到记忆转换完成，记忆ID: {memory_item.memory_id}")

        return [memory_item]
