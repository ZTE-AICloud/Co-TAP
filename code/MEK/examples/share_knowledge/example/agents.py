from __future__ import annotations

import json
from typing import Dict, List, Any
# 使用前向引用避免循环导入
from typing import TYPE_CHECKING

from mek.core.models import KnowledgeItem
from mek.core.retrieval import RetrievalType, RetrievalOptions

if TYPE_CHECKING:
    from mek.knowledge.service import KnowledgeService
    from mek.memory.base import BaseMemoryService


class ShoppingAgentA:
    """购物代理A - 负责从历史数据构建记忆和共享知识"""

    def __init__(
            self,
            knowledge_service: 'KnowledgeService',
            memory_service: 'BaseMemoryService',
            agent_name: str = "ShoppingAgent-A"
    ) -> None:
        self.knowledge_service = knowledge_service
        self.memory_service = memory_service
        self.agent_name = agent_name
        print(f"[{self.agent_name}]: 初始化完成，使用新架构（Knowledge Service + Memory Service）")

    def build_memory(self, history_file: str):

        print(f"[{self.agent_name}]: 开始从历史数据构建记忆")

        with open(history_file, 'r', encoding='utf-8') as f:
            trajectories = json.load(f)

        if isinstance(trajectories, dict):
            trajectories = [trajectories]

        for trajectory in trajectories:
            self.memory_service.add(trajectory)

        print(f"[{self.agent_name}]: 记忆构建完成")

    def share_knowledge(self, query: str) -> List[KnowledgeItem]:
        """共享知识
        
        使用knowledge service的share方法提取和共享知识
        """
        print(f"[{self.agent_name}]: 开始知识提取和共享...")

        knowledge_items = self.knowledge_service.share(
            query=query,
            context={'source_agent': self.agent_name}
        )

        print(f"[{self.agent_name}]: 知识共享完成 - 共享知识项: {len(knowledge_items)}")
        return knowledge_items


class ShoppingAgentB:
    """购物代理B - 吸收共享知识并执行任务"""

    def __init__(
            self,
            knowledge_service: 'KnowledgeService',
            memory_service: 'BaseMemoryService',
            agent_name: str = "ShoppingAgent-B"
    ) -> None:
        self.knowledge_service = knowledge_service
        self.memory_service = memory_service
        self.agent_name = agent_name
        print(f"[{self.agent_name}]: 初始化完成，使用新架构（Knowledge Service + Memory Service）")

    def absorb_knowledge(self, knowledge_items: List[KnowledgeItem]):
        """吸收共享知识
        
        使用knowledge service的absorb方法吸收外部知识
        """
        print(f"[{self.agent_name}]: 开始吸收共享知识...")

        self.knowledge_service.absorb(knowledge_items)

    def execute_task(self, task_description: str) -> Dict[str, Any]:
        """使用记忆执行任务
        
        使用memory service的retrieve方法查询相关记忆，指导任务执行
        """
        print(f"[{self.agent_name}]: 开始执行任务: {task_description}")
        print(f"[{self.agent_name}]: 在记忆库中搜索相关经验...")

        relevant_memories = self.memory_service.retrieve(
            query=task_description,
            retrieval_type=RetrievalType.SEMANTIC,
            options=RetrievalOptions(limit=10)
        )

        contents = []
        if not relevant_memories:
            print(f"[{self.agent_name}]: 未找到相关记忆")
            formatted_memory = "No available experience found"
        else:
            print(f"[{self.agent_name}]: 找到 {len(relevant_memories)} 条相关记忆，开始执行任务...")
            for memory in relevant_memories:
                steps = memory.content.get('steps', '')
                contents.append(steps)
            newline = '\n'
            formatted_memory = f"The following are the available experiences:\n{newline.join(contents)}"

        task_result = self._execute(task_description, formatted_memory)

        return {
            'task_result': task_result,
            'retrieved_memories': contents,
        }

    def _execute(self, task_description: str, memory: str) -> bool:
        # Add memory to the context and execute task
        return True
