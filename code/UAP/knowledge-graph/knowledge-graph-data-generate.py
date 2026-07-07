#!/usr/bin/env python3
"""
处理文档并存储到Neo4j，然后导出到知识图谱服务

这是一个简化版本，没有命令行参数。
所有配置参数都直接定义在这个文件中。

处理流程：
1. 启动Neo4j Docker容器（如果未运行）
2. 处理文档并将图数据存储到Neo4j中
3. 从Neo4j导出数据，通过REST API上传到知识图谱服务
   （第3步展示需要做什么，具体实现待定）
"""

import os
import subprocess
import time
import json
import requests
from neo4j import GraphDatabase
from datetime import datetime

# ============================================================================
# 配置区域 - 所有参数都在这里定义（无命令行参数）
# ============================================================================

# Neo4j Docker容器配置
NEO4J_CONTAINER_NAME = "neo4j-server"
NEO4J_IMAGE = "neo4j:2026.05.0"
NEO4J_HTTP_PORT = 7474
NEO4J_BOLT_PORT = 7687

# 工作目录
MODEL_CACHE_DIR = "./models"           # 缓存嵌入模型的目录
WORKING_DIR = "./neo4j-output"         # GraphRAG临时工作目录

# Neo4j连接配置
NEO4J_URL = f"bolt://localhost:{NEO4J_BOLT_PORT}"    # Neo4j Bolt协议URL
NEO4J_USER = "neo4j"                   # Neo4j用户名
NEO4J_PASSWORD = "password"            # Neo4j密码


# 知识图谱服务配置
KG_SERVICE_URL = "http://localhost:8080/knowledgegraph"  # KG Service REST API基础URL
KG_SERVICE_API_KEY = None             # 可选的认证API密钥

# LLM配置（必需）
LLM_BASE_URL = "..."
LLM_API_KEY = "..."
LLM_MODEL = "..."


# 输入文档
INPUT_DOCUMENT_FILE = "./demo.txt"     # 输入文本文档路径

# ============================================================================
# Neo4j Docker容器管理函数
# ============================================================================

def check_docker_installed():
    """检查Docker是否已安装"""
    try:
        subprocess.run(["docker", "--version"], capture_output=True, text=True, check=True)
        return True
    except (subprocess.CalledProcessError, FileNotFoundError):
        return False

def check_container_running(container_name):
    """检查指定名称的容器是否正在运行"""
    try:
        result = subprocess.run(
            ["docker", "ps", "-q", "-f", f"name={container_name}"],
            capture_output=True, text=True, check=True
        )
        return bool(result.stdout.strip())
    except subprocess.CalledProcessError:
        return False

def check_container_exists(container_name):
    """检查指定名称的容器是否存在（包括已停止的）"""
    try:
        result = subprocess.run(
            ["docker", "ps", "-aq", "-f", f"name={container_name}"],
            capture_output=True, text=True, check=True
        )
        return bool(result.stdout.strip())
    except subprocess.CalledProcessError:
        return False

def start_neo4j_container():
    """
    启动Neo4j Docker容器
    
    如果容器已经在运行，则等待其就绪
    如果容器存在但未运行，则删除后重新启动
    如果容器不存在，则创建并启动新容器
    """
    print("\n" + "="*70)
    print("启动Neo4j Docker容器")
    print("="*70)
    
    # 检查Docker是否安装
    if not check_docker_installed():
        print("错误：未检测到Docker。请先安装Docker。")
        return False
    
    container_name = NEO4J_CONTAINER_NAME
    
    # 检查容器状态
    if check_container_running(container_name):
        print(f"\n✓ Neo4j容器 '{container_name}' 已在运行")
        print("等待Neo4j完全启动...")
        wait_for_neo4j_ready()
        return True
    
    # 容器存在但未运行，删除它
    if check_container_exists(container_name):
        print(f"\n发现已存在的容器 '{container_name}'（未运行）")
        print("删除旧容器...")
        subprocess.run(["docker", "rm", container_name], check=True)
    
    # 启动新容器
    print(f"\n启动新的Neo4j容器 '{container_name}'...")
    
    docker_cmd = [
        "docker", "run", "-d",  # 后台运行
        "--name", container_name,
        "--publish", f"{NEO4J_HTTP_PORT}:7474",
        "--publish", f"{NEO4J_BOLT_PORT}:7687",
        "--user", f"{os.getuid()}:{os.getgid()}",
        "--env", f"NEO4J_AUTH={NEO4J_USER}/{NEO4J_PASSWORD}",
        "--env", 'NEO4J_PLUGINS=["graph-data-science"]',
        NEO4J_IMAGE
    ]
    
    print(f"执行命令: {' '.join(docker_cmd)}")
    
    try:
        result = subprocess.run(docker_cmd, capture_output=True, text=True, check=True)
        container_id = result.stdout.strip()
        print(f"✓ 容器已启动，ID: {container_id[:12]}")
        
        # 等待Neo4j完全启动
        print("\n等待Neo4j服务就绪（这可能需要1-2分钟）...")
        wait_for_neo4j_ready()
        
        return True
        
    except subprocess.CalledProcessError as e:
        print(f"✗ 启动容器失败: {e}")
        print(f"错误输出: {e.stderr}")
        return False

def wait_for_neo4j_ready(timeout=120):
    """
    等待Neo4j服务完全就绪
    
    参数:
        timeout: 最大等待时间（秒），默认120秒
    """
    import socket
    
    start_time = time.time()
    print(f"尝试连接到Neo4j端口 {NEO4J_BOLT_PORT}...")
    
    while time.time() - start_time < timeout:
        try:
            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sock.settimeout(2)
            result = sock.connect_ex(('localhost', NEO4J_BOLT_PORT))
            sock.close()
            
            if result == 0:
                print(f"✓ Neo4j端口 {NEO4J_BOLT_PORT} 已开放")
                
                # 额外等待几秒确保服务完全就绪
                time.sleep(5)
                print("✓ Neo4j服务已完全就绪")
                return True
                
        except Exception as e:
            pass
        
        print(".", end="", flush=True)
        time.sleep(2)
    
    print(f"\n✗ 等待Neo4j启动超时（{timeout}秒）")
    print("请检查Docker日志: docker logs neo4j-server")
    return False

def stop_neo4j_container():
    """停止并删除Neo4j容器"""
    container_name = NEO4J_CONTAINER_NAME
    
    if check_container_running(container_name):
        print(f"\n停止Neo4j容器 '{container_name}'...")
        subprocess.run(["docker", "stop", container_name], check=True)
        print("✓ 容器已停止")
    
    if check_container_exists(container_name):
        print(f"删除Neo4j容器 '{container_name}'...")
        subprocess.run(["docker", "rm", container_name], check=True)
        print("✓ 容器已删除")



# ============================================================================
# 导入和设置
# ============================================================================

import numpy as np
from nano_graphrag import GraphRAG
from nano_graphrag._storage import Neo4jStorage
from sentence_transformers import SentenceTransformer
from nano_graphrag._utils import wrap_embedding_func_with_attrs, compute_args_hash
from openai import AsyncOpenAI
from ratelimit import limits, RateLimitException
from backoff import on_exception, expo

def setup_directories():
    """创建必要的目录（如果不存在）"""
    os.makedirs(MODEL_CACHE_DIR, exist_ok=True)
    os.makedirs(WORKING_DIR, exist_ok=True)
    print(f"已创建目录：{MODEL_CACHE_DIR}, {WORKING_DIR}")

# ============================================================================
# LLM配置（带缓存）
# ============================================================================

@on_exception(expo, RateLimitException, max_tries=8)
@limits(calls=20, period=60)  # 每分钟最多20次调用
async def my_llm_if_cache(
    prompt, system_prompt=None, history_messages=[], **kwargs
):
    """带有缓存的LLM函数，避免重复的API调用"""
    from datetime import datetime
    now = datetime.now()
    formatted_time = now.strftime("%Y-%m-%d %H:%M:%S")
    print(f"{formatted_time}: call my_llm_if_cache")

    openai_async_client = AsyncOpenAI(
        api_key=LLM_API_KEY, base_url=LLM_BASE_URL
    )
    messages = []
    if system_prompt:
        messages.append({"role": "system", "content": system_prompt})

    # 首先检查缓存
    hashing_kv = kwargs.pop("hashing_kv", None)
    messages.extend(history_messages)
    messages.append({"role": "user", "content": prompt})
    
    if hashing_kv is not None:
        args_hash = compute_args_hash(LLM_MODEL, messages)
        if_cache_return = await hashing_kv.get_by_id(args_hash)
        if if_cache_return is not None:
            return if_cache_return["return"]

    # 调用LLM API
    response = await openai_async_client.chat.completions.create(
        model=LLM_MODEL, messages=messages, **kwargs
    )

    time.sleep(5)  # 速率限制

    # 缓存响应
    if hashing_kv is not None:
        await hashing_kv.upsert(
            {args_hash: {"return": response.choices[0].message.content, "model": LLM_MODEL}}
        )
    
    return response.choices[0].message.content

# ============================================================================
# 嵌入模型设置
# ============================================================================

def setup_embedding_model():
    """初始化句子转换器模型用于嵌入"""
    print(f"\n加载嵌入模型（缓存：{MODEL_CACHE_DIR}）...")
    EMBED_MODEL = SentenceTransformer(
        "sentence-transformers/all-MiniLM-L6-v2", 
        cache_folder=MODEL_CACHE_DIR, 
        device="cpu"
    )
    print(f"模型加载成功")
    
    @wrap_embedding_func_with_attrs(
        embedding_dim=EMBED_MODEL.get_sentence_embedding_dimension(),
        max_token_size=EMBED_MODEL.max_seq_length,
    )
    async def local_embedding(texts: list[str]) -> np.ndarray:
        return EMBED_MODEL.encode(texts, normalize_embeddings=True)
    
    return local_embedding

# ============================================================================
# 步骤0：启动Neo4j容器
# ============================================================================

def step0_start_neo4j():
    """
    步骤0：启动Neo4j Docker容器
    
    此函数：
    1. 检查Docker是否安装
    2. 检查Neo4j容器是否已在运行
    3. 如果未运行，则启动新容器
    4. 等待Neo4j服务完全就绪
    """
    print("\n" + "="*70)
    print("步骤0：启动Neo4j Docker容器")
    print("="*70)
    
    success = start_neo4j_container()
    
    if success:
        print(f"\n✓ Neo4j容器已成功启动")
        print(f"  HTTP界面：http://localhost:{NEO4J_HTTP_PORT}")
        print(f"  Bolt连接：bolt://localhost:{NEO4J_BOLT_PORT}")
        print(f"  用户名：{NEO4J_USER}")
        print(f"  密码：{NEO4J_PASSWORD}")
    else:
        print("\n✗ Neo4j容器启动失败")
        print("请确保:")
        print("  1. Docker已安装并正在运行")
        print("  2. 端口7474和7687未被占用")
        print("  3. 有足够的磁盘空间")
    
    return success

# ============================================================================
# 步骤1：处理文档并存储到Neo4j
# ============================================================================

def step1_process_and_store_in_neo4j():
    """
    步骤1：使用GraphRAG处理文档并将图数据存储到Neo4j
    
    此函数：
    1. 初始化带有Neo4j存储后端的GraphRAG
    2. 读取输入文档
    3. 处理文档以提取实体和关系
    4. 将图数据存储到Neo4j数据库
    """
    print("\n" + "="*70)
    print("步骤1：处理文档并存储到Neo4j")
    print("="*70)
    
    # 使用Neo4j存储初始化GraphRAG
    neo4j_config = {
        "neo4j_url": NEO4J_URL,
        "neo4j_auth": (NEO4J_USER, NEO4J_PASSWORD),
    }
    
    print(f"\n初始化带有Neo4j存储的GraphRAG...")
    print(f"  Neo4j URL: {NEO4J_URL}")
    print(f"  用户名：{NEO4J_USER}")
    print(f"  工作目录：{WORKING_DIR}")
    
    # 初始化GraphRAG
    graph_func = GraphRAG(
        working_dir=WORKING_DIR,
        enable_llm_cache=True,
        best_model_func=my_llm_if_cache,
        cheap_model_func=my_llm_if_cache,
        embedding_func=setup_embedding_model(),
        graph_storage_cls=Neo4jStorage,
        addon_params=neo4j_config,
    )
    
    # 读取输入文档
    print(f"\n读取输入文档：{INPUT_DOCUMENT_FILE}")
    if not os.path.exists(INPUT_DOCUMENT_FILE):
        print(f"错误：未找到输入文件 '{INPUT_DOCUMENT_FILE}'!")
        return False
        
    with open(INPUT_DOCUMENT_FILE, "r", encoding="utf-8") as f:
        content = f.read()
    print(f"已加载 {len(content)} 个字符")
    
    # 处理文档并插入到GraphRAG
    print("\n处理文档（提取实体和关系）...")
    print("根据文档大小，这可能需要几分钟...")
    graph_func.insert(content)
    
    print("\n✓ 步骤1成功完成！")
    print(f"图数据已存储在Neo4j，地址：{NEO4J_URL}")
    print(f"临时文件创建于：{WORKING_DIR}")
    
    return True

# ============================================================================
# 步骤2：从Neo4j导出并上传到KG Service
# ============================================================================

def export_entities(tx):
    # 查询所有的节点和关系
    query = """
    MATCH (n)
    RETURN elementId(n) AS elementId, 
           labels(n) AS labels, 
           properties(n) AS properties
    """
    result = tx.run(query)
    return [record.data() for record in result]

def export_entities(tx):
    # 查询所有的节点和关系
    query = """
    MATCH (n)
    RETURN elementId(n) AS elementId, 
           labels(n) AS labels, 
           properties(n) AS Props
    """
    result = tx.run(query)
    return [record.data() for record in result]

def export_relationships(tx):
    # 查询所有关系，并获取起点和终点的唯一 ID
    query = """
    MATCH (source)-[r]->(target)
    RETURN elementId(r) AS ElementId,
           type(r) AS Type,
           properties(r) AS Props,
           elementId(source) AS StartElementId,
           elementId(target) AS EndElementId
    """
    result = tx.run(query)
    return [record.data() for record in result]


def step2_export_and_upload_to_kg_service():
    """
    步骤2：从Neo4j导出数据并上传到知识图谱服务
    
    此函数演示了从Neo4j导出图数据并通过REST API上传到KG Service的工作流程。
    
    注意：这是概念性实现，展示需要做什么。
          实际的REST API客户端代码需要单独实现。
    """
    print("\n" + "="*70)
    print("步骤2：从Neo4j导出并上传到KG Service")
    print("="*70)
    
    print(f"\nKG Service配置：")
    print(f"  基础URL: {KG_SERVICE_URL}")
    if KG_SERVICE_API_KEY:
        print(f"  API密钥：[已配置]")
    else:
        print(f"  API密钥：[未设置]")
    
    print("\n--- 概念性实现 ---")
    print("(实际代码实现待定)")
    print("\n需要实现以下步骤：\n")
    
    # 步骤2.1：连接Neo4j
    print("步骤2.1：连接到Neo4j数据库")
    print("  - 使用Neo4j Python驱动或Cypher查询")
    print("  - 连接URL:", NEO4J_URL)
    print("  - 认证:", NEO4J_USER, "/[密码]")
    print()
    
    # 步骤2.2：导出节点（实体）
    print("步骤2.2：从Neo4j导出所有节点")
    print("  Cypher查询示例：")
    print("  ```cypher")
    print("  MATCH (n) RETURN n.id AS id, labels(n) AS labels, properties(n) AS props")
    print("  ```")
    print("  期望输出格式：")
    print("  [")
    print("    {\"id\": \"实体1\", \"labels\": [\"PERSON\"], \"props\": {...}},")
    print("    {\"id\": \"实体2\", \"labels\": [\"ORGANIZATION\"], \"props\": {...}},")
    print("    ...")
    print("  ]")
    print()
    # 初始化连接
    driver = GraphDatabase.driver(NEO4J_URL, auth=(NEO4J_USER, NEO4J_PASSWORD))
    with driver.session() as session:
        nodes = session.execute_read(export_entities)
        relationships = session.execute_read(export_relationships)
         
    # 保存为本地 JSON 文件
    with open("graphrag_exported.json", "w", encoding="utf-8") as f:
        json.dump(nodes, f, ensure_ascii=False, indent=4)

    driver.close()

    graph = {
        "nodes": nodes,
        "relationships": relationships
    }

    output_filename = "graphrag_nodes_edges.json"
    with open(output_filename, "w", encoding="utf-8") as f:
        json.dump(graph, f, ensure_ascii=False, indent=4)

    # 步骤2.3：导出关系
    print("步骤2.3：从Neo4j导出所有关系")
    print("  Cypher查询示例：")
    print("  ```cypher")
    print("  MATCH (s)-[r]->(t) ")
    print("  RETURN s.id AS source, type(r) AS type, t.id AS target, properties(r) AS props")
    print("  ```")
    print("  期望输出格式：")
    print("  [")
    print("    {\"source\": \"实体1\", \"type\": \"KNOWS\", \"target\": \"实体2\", \"props\": {...}},")
    print("    ...")
    print("  ]")
    print()
   

    # 步骤2.4：为KG Service API格式化数据
    print("步骤2.4：根据KG Service API要求格式化数据")
    print("  - 查看KG Service API文档了解确切的有效载荷格式")
    print("  - 通常是JSON格式，带有特定的字段名")
    print("  - 可能需要将Neo4j格式转换为KG Service格式")
    print()

    
    
    
    # 步骤2.5：上传节点到KG Service
    print("步骤2.5：通过REST API上传节点到KG Service")
    print(f"  端点：POST {KG_SERVICE_URL}/nodes")
    print("  请求头：")
    print("    Content-Type: application/json")
    if KG_SERVICE_API_KEY:
        print(f"    Authorization: Bearer {KG_SERVICE_API_KEY}")
    print("  请求体：节点对象的JSON数组")
    print("  示例代码请参考之前的版本")
    print()
    
    # 步骤2.6：上传关系到KG Service
    print("步骤2.6：通过REST API上传关系到KG Service")
    print(f"  端点：POST {KG_SERVICE_URL}/relationships")
    print("  （与节点上传类似）")
    print()

    response = requests.post(f"{KG_SERVICE_URL}/graph", json=graph)

    # 打印响应结果（解析为 JSON 字典）
    print("状态码:", response.status_code)
    print("响应内容:", response.json())
    
    # 步骤2.7：错误处理和重试
    print("步骤2.7：实现错误处理和重试逻辑")
    print("  - 优雅地处理网络错误")
    print("  - 为重试实现指数退避")
    print("  - 考虑大数据集的批量上传")
    print()
    
    # 步骤2.8：验证
    print("步骤2.8：验证成功上传")
    print("  - 检查API响应状态码")
    print("  - 查询KG Service确认数据已上传")
    print("  - 比较Neo4j和KG Service之间的节点/关系数量")
    print()
    response = requests.get(f"{KG_SERVICE_URL}/graph")
    print("状态码:", response.status_code)
    # 如果返回的是 JSON 数据，直接解析为 Python 字典
    print("响应数据:", response.json()) 
    
    print("--- 概念性实现结束 ---")
    print("\n注意：步骤2的实际实现需要：")
    print("  1. 安装Neo4j Python驱动：pip install neo4j")
    print("  2. HTTP客户端库：aiohttp或requests")
    print("  3. KG Service API文档了解确切的有效载荷格式")
    print("  4. 错误处理和重试机制")
    
    return True

# ============================================================================
# 主执行流程
# ============================================================================

def main():
    """主执行流程"""
    print("="*70)
    print("GraphRAG：Neo4j到KG Service数据迁移")
    print("="*70)
    print(f"开始时间：{datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    
    # 步骤0：启动Neo4j容器
    # step0_success = step0_start_neo4j()
    
    # if not step0_success:
    #     print("\n✗ 步骤0失败（Neo4j容器启动失败）。中止。")
    #     return
    
    # 设置目录
    setup_directories()
    
    # 执行步骤1：处理并存储到Neo4j
    step1_success = step1_process_and_store_in_neo4j()
    
    if not step1_success:
        print("\n✗ 步骤1失败。中止。")
        return
    
    # 在进入步骤2之前等待片刻
    print("\n等待5秒后继续步骤2...")
    time.sleep(5)
    
    # 执行步骤2：导出并上传到KG Service
    step2_success = step2_export_and_upload_to_kg_service()
    
    if not step2_success:
        print("\n✗ 步骤2失败。")
        return
    
    # 完成
    print("\n" + "="*70)
    print("处理完成")
    print("="*70)
    print(f"结束时间：{datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print("\n摘要：")
    print("  ✓ 步骤0：Neo4j容器已启动")
    print("  ✓ 步骤1：文档已处理并存储在Neo4j")
    print("  ✓ 步骤2：导出和上传工作流程已文档化")
    print(f"\n后续步骤：")
    print(f"  1. 为步骤2实现实际的REST API客户端代码")
    print(f"  2. 用您的KG Service实例测试")
    print(f"  3. 监控上传进度并处理任何错误")
    print(f"\nNeo4j浏览器访问：http://localhost:{NEO4J_HTTP_PORT}")
    print(f"用户名：{NEO4J_USER}，密码：{NEO4J_PASSWORD}")

if __name__ == "__main__":
    main()