# AgentGraph

agent 知识图谱，用于agent节点定义以及agent 之间关系管理

功能定位‌:

* **agent 类型定义** - 根据agent 类型定义agent 节点
* **agent 关系定义** - 指定不同类型agent间的关系

* **数据存储** - 默认使用外置的 neo4j 数据库存储节点关系图谱， 不设置数据库时不支持知识图谱功能


# 快速开始

## 1、安装neo4j：

参考官网安装教程
https://neo4j.com/docs/operations-manual/current/installation/


## 3、启动服务并设置neo4j连接信息
根据您的部署方式，执行启动命令，例如：
```bash
./uapregistry --agentgraph-uri neo4j://127.0.0.1:7687 --agentgraph-username neo4j --agentgraph-password mypsw --agentgraph-database databasename
```
## 4、验证服务
注册中心启动后，您可以通过以下方式验证, 该接口返回当前环境中的所有agent 以及其关系：
```bash
curl -X GET http://127.0.0.1:8080/agentgraph/graph
```

# 功能列表

## 整系统

### 导出系统数据
不带分页参数则导出所有数据
```bash
GET	/agentgraph/graph?page=0&limit=1000
```

### 导入数据

导入agent 节点及节点间的关系
```bash
POST	/agentgraph/graph
```

## 节点

### 创建节点

```bash
POST	/agentgraph/nodes
```

### 批量创建节点

```bash
POST	/agentgraph/nodes/bulk
```

### 更新节点

```bash
PUT	/agentgraph/nodes/{elementId}
```

### 删除节点

```bash
DELETE	/agentgraph/nodes/{elementId}
```

### 查询节点

```bash
GET /agentgraph/nodes/{elementId}
```

### 分页查询节点

```bash
GET /agentgraph/nodes?page=0&limit=100&label=Person
```



## 关系

### 创建关系
```bash
POST	/agentgraph/relationships
```

### 批量创建关系

```bash
POST	/agentgraph/relationships/bulk
```

### 更新关系属性

```bash
PUT	/agentgraph/relationships/{elementId}
```

### 删除关系

```bash
DELETE	/agentgraph/relationships/{elementId}
```

### 查询关系

```bash
GET	/agentgraph/relationships/{elementId}
```


### 分页查询关系

```bash
GET	agentgraph/relationships?page=0&limit=100&type=depend
```