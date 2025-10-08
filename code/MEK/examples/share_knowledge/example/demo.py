#!/usr/bin/env python3
"""
MEK协议多代理知识共享演示

新架构演示：
1. Agent注入knowledge service和memory service
2. Knowledge service注入extract service和memory service  
3. 四个核心功能：
   - 功能1：从原始数据生成记忆（Memory Service.add）
   - 功能2：共享知识（Knowledge Service.share）
   - 功能3：吸收知识（Knowledge Service.absorb）
   - 功能4：使用记忆（Memory Service.retrieve）
"""

import argparse
import os
import shutil
import sys

# 添加项目根目录到Python路径
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from example.services import ServicesFactory
from example.agents import ShoppingAgentA, ShoppingAgentB


def clear_memory_stores():
    """清理记忆库"""
    print("🧹 清理记忆库阶段")
    print("-" * 50)

    base_memory_path = "data/memory"

    # 清理agent_a和agent_b的记忆库
    for agent_name in ["agent_a", "agent_b"]:
        agent_path = os.path.join(base_memory_path, agent_name)
        if os.path.exists(agent_path):
            shutil.rmtree(agent_path)
            print(f"🗑️  已清理 {agent_name} 的记忆库")

    print("✅ 记忆库清理完成\n")


def create_agents(factory: ServicesFactory):
    """创建代理实例"""
    print("🔧 创建独立的Agent...")
    print("-" * 50)

    # 为Agent-A创建服务实例
    print("[System]: 为ShoppingAgent-A创建独立服务...")
    knowledge_service_a, memory_service_a = factory.create_agent_services("agent_a")
    agent_a = ShoppingAgentA(knowledge_service_a, memory_service_a, "ShoppingAgent-A")

    # 为Agent-B创建服务实例
    print("[System]: 为ShoppingAgent-B创建独立服务...")
    knowledge_service_b, memory_service_b = factory.create_agent_services("agent_b")
    agent_b = ShoppingAgentB(knowledge_service_b, memory_service_b, "ShoppingAgent-B")

    print("[System]: 两个购物代理初始化完成")
    print("✅ Agent创建完成\n")

    return agent_a, agent_b


def main():
    parser = argparse.ArgumentParser(description="MEK协议多代理知识共享演示")
    parser.add_argument("--query", default="购物支付流程", help="任务查询描述")
    parser.add_argument("--clean", action="store_true", help="清理记忆库")
    args = parser.parse_args()

    print("=== MEK协议 Multi-Agent知识共享演示 ===\n")

    # 清理记忆库（如果需要）
    if args.clean:
        clear_memory_stores()

    factory = ServicesFactory()
    agent_a, agent_b = create_agents(factory)

    print("📚 Step 1 - 记忆构建")
    print("-" * 50)
    agent_a.build_memory("example/data/workflow_history.json")
    print("✅ Agent-A完成记忆构建\n")

    print("🔄 Step 2 - 知识共享")
    print("-" * 50)
    shared_knowledge = agent_a.share_knowledge(args.query)

    print("🧠 Step 3 - 知识吸收")
    print("-" * 50)
    agent_b.absorb_knowledge(shared_knowledge)
    print("✅ Agent-B成功吸收共享知识\n")

    print("🎯 Step 4 - 任务执行")
    print("-" * 50)
    execution_result = agent_b.execute_task(args.query)

    print("📊 任务执行结果:")
    print(f"  - 任务: {args.query}")
    print(f"  - 任务结果: {execution_result['task_result']}")
    print("✅ Agent-B完成任务执行阶段\n")


if __name__ == "__main__":
    main()
