# UAP注册中心

UAP注册中心是服务治理核心组件，维护服务全生命周期状态

功能定位‌:

* **服务管理** - 处理注册/注销请求，维护服务元数据（含协议扩展字段）。

* **数据分层存储** - 
  * **缓存层‌**：内存数据库实现快速响应。
  * **‌持久层**‌：分布式KV存储保障数据强一致性。

* **健康监控** - 定期检测服务可用性，剔除异常实例。


# 快速开始

## 1、在启动注册中心前，请确保consul集群已经启动，并通过命令行设置所需的环境变量：
```bash
export HEALTH_CHECK_ENABLE=true
export CONSUL_IP=127.0.0.1
export CONSUL_HTTP_PORT=8500
export UAPREGISTRY_CENTER_PORT=8080
```
其中：
HEALTH_CHECK_ENABLE用于控制是否开启健康检查功能。
CONSUL_IP用于配置要连接的Consul集群中某个服务器的IP地址。
CONSUL_HTTP_PORT用于配置要连接的Consul集群中某个服务器的HTTP端口。
UAPREGISTRY_CENTER_PORT用于配置注册中心自身的HTTP端口。


## 2、设置日志配置文件log.yml
```yaml
#level:Emergency、Alert、Critical、Error、
#Warn、Notice、Info、Debug
console:
  level: Info
file:
  filename: ../uapregistry-works/logs/uapregistry.log
  level: Info
  maxlines: 300000
  #metric:M
  maxsize: 20
  daily: false
  maxdays: 3
  rotate: true
  perm: 0640
```

## 3、启动服务
根据您的部署方式，执行启动命令，例如：
```bash
#!/bin/sh

export HEALTH_CHECK_ENABLE=true
export CONSUL_IP=10.234.88.163
export CONSUL_HTTP_PORT=8500
export UAPREGISTRY_CENTER_PORT=8080

./uapregistry
```
## 4、验证服务
注册中心启动后，您可以通过以下方式验证：
```bash
curl http://127.0.0.1:8080/health -v
```

