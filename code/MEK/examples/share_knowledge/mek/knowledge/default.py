from __future__ import annotations

from typing import Dict, List, Optional
from typing import TYPE_CHECKING

from .service import KnowledgeService
from ..core.models import KnowledgeItem, MemoryItem
from ..core.retrieval import RetrievalType, RetrievalOptions
from ..core.types import ImportanceLevel

if TYPE_CHECKING:
    from ..extract.extractor import ExtractService
    from ..memory.base import BaseMemoryService


class DefaultKnowledgeService(KnowledgeService):
    """默认知识服务实现
    
    依赖ExtractService和MemoryService实现知识的共享和吸收
    """

    def __init__(
            self,
            extract_service: 'ExtractService',
            memory_service: 'BaseMemoryService'
    ) -> None:
        self.extract_service = extract_service
        self.memory_service = memory_service

    def share(self, query: str, context: Optional[Dict[str, object]] = None) -> List[KnowledgeItem]:
        """共享知识
        
        1. 调用memory service的retrieve方法查询获取相关记忆条目
        2. 调用extract service的memory_to_knowledge方法转换为标准知识对象
        """
        print(f"[DefaultKnowledgeService]: 检索相关记忆...")
        memories = self.memory_service.retrieve(
            query=query,
            retrieval_type=RetrievalType.SEMANTIC,
            options=RetrievalOptions(limit=10)
        )

        if not memories:
            print(f"[DefaultKnowledgeService]: 未找到相关记忆")
            return []
        print(f"[DefaultKnowledgeService]: 转换 {len(memories)} 条记忆为知识...")
        knowledge_items = self.extract_service.memory_to_knowledge(
            memories=memories,
            context=context
        )

        print(f"[DefaultKnowledgeService]: 知识共享完成，生成 {len(knowledge_items)} 个知识项")
        return knowledge_items

    def absorb(self, knowledge_items: List[KnowledgeItem]) -> List[str]:
        """吸收知识
        
        1. 调用extract service的knowledge_to_memory方法将知识转换为记忆
        2. 调用memory service的retrieve方法检索相关记忆条目
        3. 判断是否与旧记忆存在冲突，决定需要更新或保留的条目
        4. 调用memory service的update方法入库记忆
        """
        print(f"[DefaultKnowledgeService]: 开始知识吸收，知识项数量: {len(knowledge_items)}")

        absorbed_memory_ids = []

        for knowledge_item in knowledge_items:
            try:
                print(f"[DefaultKnowledgeService]: 处理知识项: {knowledge_item.title}")

                memory_items = self.extract_service.knowledge_to_memory(knowledge_item)
                # 当前只支持工作流类型，知识理解为工作流记忆，仅一条记录。
                memory_item = memory_items[0]

                existing_memories = self.memory_service.retrieve(
                    query=memory_item.content.get('workflow_name', ''),
                    retrieval_type=RetrievalType.EXACT,
                    options=RetrievalOptions(limit=5)
                )

                if existing_memories:
                    conflicts = self._detect_conflicts(memory_item, existing_memories)
                    if conflicts:
                        memory_item = self._resolve_conflicts(memory_item, conflicts)

                memory_uuid = self.memory_service.update(memory_item)
                absorbed_memory_ids.append(memory_uuid)

                print(f"[DefaultKnowledgeService]: 成功吸收知识 {knowledge_item.knowledge_id} -> 记忆 {memory_uuid}")

            except Exception as e:
                print(f"[DefaultKnowledgeService]: 吸收知识失败 {knowledge_item.knowledge_id}: {e}")
                continue

        return absorbed_memory_ids

    def _detect_conflicts(self, new_memory: MemoryItem, existing_memories: List[MemoryItem]) -> List[MemoryItem]:
        """检测与新记忆冲突的现有记忆"""
        conflicts = []

        for existing_memory in existing_memories:
            if self._is_conflict(new_memory, existing_memory):
                conflicts.append(existing_memory)

        return conflicts

    def _is_conflict(self, memory1: MemoryItem, memory2: MemoryItem) -> bool:
        """判断两个记忆是否冲突"""
        name1 = memory1.content.get('workflow_name', '')
        name2 = memory2.content.get('workflow_name', '')

        if name1 and name2:
            return name1.lower() in name2.lower() or name2.lower() in name1.lower()

        return False

    def _resolve_conflicts(self, new_memory: MemoryItem, conflicts: List[MemoryItem]) -> MemoryItem:
        """解决记忆冲突"""
        if not conflicts:
            return new_memory

        related_ids = set(new_memory.related_memory_ids)
        for conflict in conflicts:
            related_ids.update(conflict.related_memory_ids)
            related_ids.add(conflict.memory_id)

        new_memory.related_memory_ids = list(related_ids)

        if any(c.content.get('result_status') == 'success' for c in conflicts):
            if new_memory.importance != ImportanceLevel.HIGH:
                new_memory.importance = ImportanceLevel.MEDIUM

        return new_memory
