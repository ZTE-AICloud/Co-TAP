from __future__ import annotations

import hashlib
import json
import time
from datetime import datetime
from typing import Any, Dict, List, Optional, Union

from .base import BaseMemoryService
from ..core.models import MemoryItem
from ..core.repository import MemoryRepository
from ..core.retrieval import RetrievalType, RetrievalOptions
from ..core.types import MemoryType, ImportanceLevel
from ..llm_client.models import Message, RoleType
from ..llm_client.openai_client import OpenAIClient
from ..utils.data_util import json_loads_str


class WorkflowMemoryService(BaseMemoryService):
    """工作流记忆服务
    
    负责工作流记忆的CRUD操作，包括从轨迹数据中提取工作流记忆
    """

    def __init__(self, repository: MemoryRepository, llm: Optional[OpenAIClient] = None) -> None:  # type: ignore[name-defined]
        super().__init__(repository, MemoryType.WORKFLOW)
        self.llm = llm

    def add(self, input_data: Union[str, Dict[str, Any]], **kwargs) -> List[str]:
        """从原始数据中提取生成记忆
        
        使用LLM从轨迹数据中提取通用工作流，生成记忆，并存入数据库
        """
        print(f"[WorkflowMemoryService]: 开始从原始数据生成记忆...")

        trajectory_data = self._parse_input_data(input_data)
        workflow_data = self._extract_workflow_with_llm(trajectory_data)
        memory_item = self._create_memory_item(workflow_data)
        conflicts = self._detect_conflicts(memory_item)
        if conflicts:
            memory_item = self._resolve_conflicts(memory_item, conflicts)
            print(f"[WorkflowMemoryService]: 解决了 {len(conflicts)} 个冲突")
        memory_uuid = self.update(memory_item)

        print(f"[WorkflowMemoryService]: 成功生成记忆 {memory_uuid}")
        return [memory_uuid]

    def retrieve(
            self,
            query: str,
            retrieval_type: RetrievalType = RetrievalType.SEMANTIC,
            options: Optional[RetrievalOptions] = None
    ) -> List[MemoryItem]:
        """检索相关记忆"""
        print(f"[WorkflowMemoryService]: 检索记忆 - query={query}")

        opts = options or RetrievalOptions()

        # 获取所有工作流记忆
        all_memories = [m for m in self.repository.list_all()
                        if m.memory_type == MemoryType.WORKFLOW]

        if retrieval_type == RetrievalType.EXACT:
            # 精确匹配
            results = self._exact_match_retrieve(all_memories, query)
        else:
            # 其他检索类型的默认实现
            results = all_memories

        # 按重要性和时间排序
        results.sort(key=lambda m: (m.importance.value, m.created_at), reverse=True)

        return results[:opts.limit]

    def update(self, memory_item: MemoryItem) -> str:
        """更新或添加记忆条目"""
        existing_memory = self.repository.get(memory_item.memory_id)

        if existing_memory:
            # 更新现有记忆
            self.repository.update(memory_item.memory_id, memory_item)
            print(f"[WorkflowMemoryService]: 更新记忆 {memory_item.memory_id}")
        else:
            # 添加新记忆
            self.repository.add(memory_item)
            print(f"[WorkflowMemoryService]: 添加新记忆 {memory_item.memory_id}")

        return memory_item.memory_id

    def delete(self, memory_id: str) -> bool:
        """删除指定记忆"""
        try:
            self.repository.delete(memory_id)
            print(f"[WorkflowMemoryService]: 删除记忆 {memory_id}")
            return True
        except Exception as e:
            print(f"[WorkflowMemoryService]: 删除记忆失败 {memory_id}: {e}")
            return False

    def _parse_input_data(self, input_data: Union[str, Dict[str, Any]]) -> Dict[str, Any]:
        """解析输入数据"""
        if isinstance(input_data, dict):
            return input_data
        elif isinstance(input_data, str):
            try:
                return json.loads(input_data)
            except json.JSONDecodeError:
                return {"raw_text": input_data}
        else:
            return {"raw_data": str(input_data)}

    def _extract_workflow_with_llm(self, trajectory_data: Dict[str, Any]) -> Dict[str, Any]:
        """使用LLM从轨迹数据中提取通用工作流"""
        if not self.llm:
            return self._simple_workflow_extraction(trajectory_data)

        try:
            prompt = f"""
            请分析以下用户轨迹数据，提取出通用的工作流模式：

            轨迹数据：
            {json.dumps(trajectory_data, ensure_ascii=False, indent=2)}

            请输出一个标准化的工作流，包含：
            1. workflow_name: 工作流名称
            2. steps: 标准化步骤列表，每个步骤包含name和description

            输出格式为JSON：
            {{
                "workflow_name": "具体的工作流名称",
                "steps": [
                    {{"name": "步骤名称", "description": "步骤描述"}},
                    ...
                ]
            }}
            """
            messages = [Message(role=RoleType.User, content=prompt)]
            response = self.llm.generate_response(messages)
            result = json_loads_str(response.get('content', '').strip())
            return result
        except Exception as e:
            print(f"[WorkflowMemoryService]: LLM提取失败: {e}")
            return self._simple_workflow_extraction(trajectory_data)

    def _simple_workflow_extraction(self, trajectory_data: Dict[str, Any]) -> Dict[str, Any]:
        workflow_name = trajectory_data.get('workflow_name', '工作流程')
        steps = []

        for step in trajectory_data.get('steps', []):
            steps.append({
                'name': step.get('step_name', step.get('action', '未知步骤')),
                'description': f"执行{step.get('action', '操作')}"
            })

        return {
            'workflow_name': workflow_name,
            'steps': steps
        }

    def _create_memory_item(self, workflow_data: Dict[str, Any]) -> MemoryItem:
        memory_id = self._generate_memory_id(f"wf-{workflow_data.get('workflow_name', 'unknown')}")

        content = {
            'task_description': f"工作流-{workflow_data.get('workflow_name', '未知')}",
            'workflow_name': workflow_data.get('workflow_name', '未知工作流'),
            'steps': '\n'.join([f'{step["name"]}:{step["description"]}' for step in workflow_data.get('steps', [])]),
            'result_status': 'success',
            'domain': '通用'
        }

        return MemoryItem(
            memory_id=memory_id,
            memory_type=MemoryType.WORKFLOW,
            created_at=datetime.utcnow(),
            content=content,
            importance=ImportanceLevel.MEDIUM,
            related_memory_ids=[]
        )

    def _exact_match_retrieve(self, memories: List[MemoryItem], query: str) -> List[MemoryItem]:
        """精确匹配检索"""
        query_lower = query.lower()
        results = []

        for memory in memories:
            content_text = self._extract_searchable_text(memory)
            if query_lower in content_text.lower():
                results.append(memory)

        return results

    def _extract_searchable_text(self, memory: MemoryItem) -> str:
        """提取记忆的可搜索文本"""
        content = memory.content

        text_parts = [
            content.get('task_description', ''),
            content.get('workflow_name', ''),
            content.get('domain', ''),
            content.get('steps', '')
        ]

        return ' '.join(filter(None, text_parts))

    def _detect_conflicts(self, new_memory: MemoryItem) -> List[MemoryItem]:
        """检测与新记忆冲突的现有记忆"""
        existing_memories = [m for m in self.repository.list_all()
                             if m.memory_type == MemoryType.WORKFLOW]

        conflicts = []
        for memory in existing_memories:
            if self._is_conflict(new_memory, memory):
                conflicts.append(memory)

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

    def _generate_memory_id(self, prefix: str = "memory") -> str:
        """生成记忆ID"""
        timestamp = str(int(time.time() * 1000))
        content_hash = hashlib.md5(timestamp.encode()).hexdigest()[:8]
        return f"{prefix}-{content_hash}"
