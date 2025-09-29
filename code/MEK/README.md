# MEK - Memory-Extraction-Knowledge Protocol

MEK协议是一个高效的智能代理记忆管理和知识共享系统。通过**三层架构**：**Extract(转换)**、**Memory(记忆)**、**Knowledge(知识)**，实现多代理间的智能知识共享与协同学习。

## 🎯 核心特性

- **🧠 三层分离架构**: Extract层专门处理数据转换，Memory层纯粹CRUD，Knowledge层负责清理
- **🔄 统一转换服务**: `ExtractService`抽取核心转换逻辑，支持知识↔记忆双向转换
- **🤖 智能记忆生成**: LLM驱动的轨迹分析和工作流提取，包含冲突解决
- **🔍 语义检索**: 智能记忆检索，支持精确、语义、时序等多种检索模式
- **🔄 知识共享**: 标准化`KnowledgeItem`格式，实现跨代理知识传递
- **⚡ 冲突解决**: 自动检测记忆冲突并进行智能融合
- **🔧 配置驱动**: 支持任何OpenAI API兼容的LLM和Embedding服务
- **📦 轻量架构**: 简洁代码结构，无冗余实现，高性能运行

## 🏗️ 系统架构

```
MEK Protocol - 三层分离架构
├── 🔄 Extract Layer (转换层)
│   └── ExtractService           # 统一转换服务
│       ├── 知识→记忆: 感知→理解→融合(含冲突解决)
│       └── 记忆→知识: 筛选→泛化→标准化→提取
├── 🧠 Memory Layer (记忆层)
│   ├── BaseMemoryService        # 统一记忆服务基类
│   ├── WorkflowMemoryService    # 工作流记忆服务 (纯CRUD)
│   └── FileMemoryRepository     # 文件存储仓库
├── 📚 Knowledge Layer (知识层)  
│   └── DefaultKnowledgeService  # 知识清理与标准化
├── 🔍 Core Layer (核心层)
│   ├── MemoryItem              # 记忆数据模型
│   ├── KnowledgeItem           # 知识共享模型
│   ├── RetrievalOptions        # 检索配置
│   └── MemoryType              # 记忆类型枚举
└── 🚀 Engine Layer (引擎层)
    ├── OpenAIClient            # LLM客户端
    └── ConfigManager           # 配置管理器
```

## 🚀 四个核心功能

### 功能1：记忆构建
- **实现**: Memory Service的add方法
- **流程**: 原始轨迹数据 → LLM分析 → 生成MemoryItem → 冲突解决 → 存储
- **使用**: `memory_service.add(trajectory_data)`

### 功能2：知识共享
- **实现**: Knowledge Service的share方法
- **流程**: 查询 → Memory检索 → Extract转换 → 生成KnowledgeItem列表
- **使用**: `knowledge_service.share(query)`

### 功能3：知识吸收
- **实现**: Knowledge Service的absorb方法
- **流程**: KnowledgeItem → Extract转换 → 冲突检测 → Memory存储
- **使用**: `knowledge_service.absorb(knowledge_items)`

### 功能4：记忆使用
- **实现**: Memory Service的retrieve方法
- **流程**: 任务查询 → Memory检索 → 生成执行方案
- **使用**: `memory_service.retrieve(query)`

## 🚀 快速开始

### 环境要求

- Python 3.12+
- 支持OpenAI API格式的LLM服务
- 支持OpenAI API格式的Embedding服务

### 安装依赖

```bash
pip install -r requirements.txt
```

### 配置设置

1. 复制配置模板：
```bash
cp config_example.yaml config.yaml
```

2. 编辑`config.yaml`，填入你的API配置：

```yaml
# LLM Chat Client配置
llm:
  provider: "openai"
  model_name: "gpt-4"
  base_url: "https://api.openai.com/v1"
  api_key: "your-llm-api-key"
  
# Embedding Client配置  
embedding:
  provider: "openai"
  model_name: "text-embedding-3-small"
  base_url: "https://api.openai.com/v1"
  api_key: "your-embedding-api-key"
  embedding_dim: 1024

# MEK系统配置
mek:
  memory_storage_path: "data/memory"
  default_retrieval_limit: 20
  default_importance_weight: 1.0
```

### 运行示例

```bash
# 设置代理（如果需要）
export all_proxy=http://your-proxy:port

# 运行MEK协议演示 - 多代理知识共享
python example/demo.py --query "购物支付流程"

# 清理记忆库重新开始
python example/demo.py --query "购物支付流程" --clean
```

## 📖 使用示例

### 1. 三层架构核心使用

```python
from example.services import ServicesFactory

# 创建服务工厂
factory = ServicesFactory()

# 创建三层架构服务
memory_service = factory.make_memory_service("agent_a")
knowledge_service = factory.make_knowledge_service("agent_a")
extract_service = factory.make_extract_service()

# 四个核心功能
# 功能1：记忆构建
trajectory_data = {
    "session_id": "shopping_001",
    "workflow_name": "购物流程",
    "steps": [{"step_name": "浏览商品", "action": "查看商品列表"}]
}
memory_uuids = memory_service.add(trajectory_data)

# 功能2：知识共享  
knowledge_items = knowledge_service.share("购物支付流程")

# 功能3：知识吸收
absorbed_count = knowledge_service.absorb(knowledge_items)

# 功能4：记忆使用
memories = memory_service.retrieve("购物支付流程")
```

### 2. 多代理知识共享演示

```python
from example.agents import ShoppingAgentA, ShoppingAgentB

# 创建两个独立的代理 (各自拥有独立的三层架构)
agent_a = ShoppingAgentA(knowledge_service_a, memory_service_a, "Agent-A")
agent_b = ShoppingAgentB(knowledge_service_b, memory_service_b, "Agent-B")

# 🔧 Agent-A: 从历史数据构建记忆 (Memory层处理轨迹→记忆)
memory_count = agent_a.build_memory("example/data/workflow_history.json")
print(f"Agent-A构建了 {memory_count} 条工作流记忆")

# 🔄 Agent-A: 提取并共享知识 (Extract层处理记忆→知识)
shared_knowledge = agent_a.share_knowledge("购物支付流程")
print(f"共享了 {len(shared_knowledge)} 项知识")

# 🧠 Agent-B: 吸收共享知识 (Extract层处理知识→记忆，含冲突解决)
success_count = agent_b.absorb_shared_knowledge(shared_knowledge)
print(f"Agent-B成功吸收 {success_count} 项知识")

# 🎯 Agent-B: 使用知识执行任务
result = agent_b.execute_task_with_knowledge("购物支付流程")
print(f"任务执行结果: {result['task_result']}")
```

### 3. 直接使用Extract层转换服务

```python
from mek.extract import ExtractService

# 创建Extract服务
extractor = ExtractService()

# 记忆到知识转换（功能2的核心）
knowledge_items = extractor.memory_to_knowledge(memories)

# 知识到记忆转换（功能3的核心）
memory_item = extractor.knowledge_to_memory(knowledge_item)
```

## 📁 项目结构

```
share_knowledge/               # MEK协议示例项目根目录
├── mek/                       # 🧠 核心MEK协议参考实现
│   ├── extract/               # 🔄 Extract层 - 转换算子
│   │   ├── extractor.py       # ExtractService统一转换服务
│   │   └── __init__.py
│   ├── core/                  # 📊 核心数据模型与类型
│   │   ├── models.py          # MemoryItem & KnowledgeItem
│   │   ├── repository.py      # 记忆库抽象接口
│   │   ├── types.py           # 记忆类型与重要性枚举
│   │   ├── retrieval.py       # 检索类型与配置
│   │   └── __init__.py
│   ├── memory/                # 🧠 Memory层 - 纯粹记忆CRUD
│   │   ├── base.py            # BaseMemoryService基类
│   │   ├── workflow_memory.py # 工作流记忆服务实现
│   │   └── __init__.py
│   ├── knowledge/             # 📚 Knowledge层 - 知识清理
│   │   ├── service.py         # 知识服务接口 (简化版)
│   │   ├── default.py         # 默认知识服务实现
│   │   └── __init__.py
│   ├── storage/               # 💾 存储实现
│   │   ├── file_repository.py # 文件存储仓库
│   │   └── __init__.py
│   ├── llm_client/            # 🤖 LLM客户端
│   │   ├── openai_client.py   # OpenAI兼容客户端
│   │   ├── embedding_client.py # Embedding客户端
│   │   ├── base.py            # 客户端基类
│   │   ├── config.py          # 客户端配置
│   │   ├── models.py          # 消息模型
│   │   └── __init__.py
│   ├── config/                # ⚙️ 配置管理
│   │   ├── manager.py         # 配置管理器
│   │   └── __init__.py
│   ├── utils/                 # 🛠️ 工具函数
│   │   ├── data_util.py       # 数据处理工具
│   │   ├── file_util.py       # 文件处理工具
│   │   ├── log_util.py        # 日志工具
│   │   └── __init__.py
│   └── __init__.py
├── example/                   # 📋 使用示例
│   ├── agents.py              # 示例代理实现
│   ├── services.py            # 服务工厂
│   ├── demo.py                # 完整演示脚本
│   ├── data/
│   │   └── workflow_history.json # 示例轨迹数据
│   └── __init__.py
├── data/                      # 📂 数据目录
│   └── memory/                # 记忆存储目录
├── config.yaml                # ⚙️ 配置文件
├── config_example.yaml        # 配置模板
├── requirements.txt           # 依赖列表
└── README.md                  # 项目说明
```

## 🎮 演示说明

### 演示场景

项目包含完整的**多代理电商购物知识共享**演示，展示MEK三层架构的核心能力：

1. **🔧 记忆构建**: ShoppingAgent-A使用Memory层从轨迹中提取工作流模式
2. **🔄 知识共享**: Agent-A使用Extract层将记忆转化为标准化KnowledgeItem
3. **🧠 知识吸收**: Agent-B使用Extract层吸收共享知识，含冲突检测和融合
4. **🎯 任务执行**: Agent-B利用新知识执行购物任务

### 三层架构工作流程

```
Raw Data → Extract Layer → Memory Layer → Knowledge Layer → Shared Knowledge
    ↓           (转换)        (存储)        (清理)           ↓
 轨迹数据    感知→理解→融合    纯CRUD操作   脱敏→标准化    知识项目
    ↑                                                        ↓
 任务查询  ← Extract Layer ← Memory Layer ← Knowledge Layer ← 其他Agent
         (检索→转换)      (CRUD操作)      (接收知识)
```

### 演示数据

`example/data/workflow_history.json` 包含完整的用户购物轨迹：

```json
{
  "session_id": "shopping_001",
  "user_id": "user_12345", 
  "workflow_name": "智能手机购买流程",
  "domain": "电商购物",
  "steps": [
    {
      "step_name": "商品浏览",
      "action": "浏览手机列表",
      "duration_seconds": 120
    },
    {
      "step_name": "商品详情查看",
      "action": "查看iPhone 15详情", 
      "duration_seconds": 180
    }
  ],
  "result_status": "success"
}
```

### 演示命令

```bash
# 基础演示 - 展示完整的三层架构知识共享流程
python example/demo.py --query "购物支付流程"

# 清理演示 - 重置记忆库
python example/demo.py --query "购物支付流程" --clean

# 自定义查询演示
python example/demo.py --query "用户注册流程"
```

## 🔍 核心概念

### 三层架构设计原则

- **Extract Layer**: 专门处理数据转换，包含所有"理解"和"分析"逻辑
- **Memory Layer**: 纯粹的CRUD操作，不包含业务逻辑
- **Knowledge Layer**: 只负责最后的清理和标准化

### 记忆类型 (Memory Types)

- **WORKFLOW**: 业务流程、操作步骤记忆 (已实现)
- **PROFILE**: 用户画像、偏好记忆 (扩展点)
- **SEMANTIC**: 概念关系、语义记忆 (扩展点)
- **EPISODE**: 具体事件、情节记忆 (扩展点)

### 检索类型 (Retrieval Types)

- **EXACT**: 精确文本匹配
- **SEMANTIC**: 语义相似度检索 (需要embedding服务)
- **TEMPORAL**: 时间顺序检索
- **ASSOCIATIVE**: 关联关系扩展检索

### Extract层转换服务

`ExtractService`是MEK协议的核心转换引擎：

```python
class ExtractService:
    def knowledge_to_memory(self, knowledge_item) -> MemoryItem:
        """知识到记忆转换：感知→理解→融合(含冲突解决)"""
        
    def memory_to_knowledge(self, memories, context=None) -> List[KnowledgeItem]:
        """记忆到知识转换：筛选→泛化→标准化→提取"""
```

## 🛠️ 扩展开发

### 1. 创建新的记忆类型

```python
from mek.memory.base import BaseMemoryService
from mek.core.types import MemoryType

class ProfileMemoryService(BaseMemoryService):
    def __init__(self, repository, llm=None):
        super().__init__(repository, MemoryType.PROFILE)
        self.llm = llm

    def add(self, input_data, **kwargs):
        """注意：input_data应该是Extract层处理后的MemoryItem"""
        # 实现纯粹的存储逻辑，不包含转换逻辑
        pass
        
    def retrieve(self, query, retrieval_type, options):
        # 实现画像检索逻辑
        pass
```

### 2. 扩展Extract层转换逻辑

```python
from mek.extract import ExtractService

class CustomExtractService(ExtractService):
    def _understand(self, perceived_data, context=None):
        """重写理解逻辑，加入自定义业务理解"""
        # 调用父类方法
        base_result = super()._understand(perceived_data, context)
        
        # 添加自定义逻辑
        if perceived_data.get('data_type') == 'custom':
            # 自定义理解逻辑
            pass
            
        return base_result
```

### 3. 自定义知识共享格式

```python
from mek.core.models import KnowledgeItem

knowledge = KnowledgeItem(
    knowledge_id="custom-workflow-001",
    title="自定义业务流程",
    content={
        "workflow_name": "订单处理流程",
        "steps": [
            {"name": "订单验证", "description": "验证订单信息"},
            {"name": "库存检查", "description": "检查商品库存"},
            {"name": "订单确认", "description": "确认订单详情"}
        ],
        "domain": "电商",
        "result_status": "success"
    },
    source_memory_type=MemoryType.WORKFLOW,
    source_agent="OrderAgent"
)
```

## 📈 性能特点

- **🚀 三层分离**: Extract层专门处理转换，Memory层专门处理存储，Knowledge层专门处理清理
- **⚡ 统一转换**: 所有数据转换通过`ExtractService`统一处理，逻辑清晰
- **🧠 智能融合**: Extract层自动检测记忆冲突并进行智能合并
- **🔄 快速检索**: Memory层支持多种检索策略，智能排序
- **📦 简洁架构**: 无冗余代码，模块化设计，易于扩展和维护

## 📄 相关文档

- [使用示例](example/)

## 📋 更新日志

### v0.10.0 (Latest)
- ✨ 三层分离架构：Extract-Memory-Knowledge
- 🔄 统一转换服务：ExtractService处理所有数据转换
- 🔧 完整演示：多代理知识共享演示系统

## 📜 许可证

本项目采用 Apache 2.0 许可证

---

**MEK Protocol** - 三层架构，智能转换，高效共享 🚀

> "分层清晰的架构设计，统一的转换服务。MEK协议让智能代理的记忆管理更简洁、更高效。"
