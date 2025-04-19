# WXProxy服务
微信 Open Api代理服务，提供稳定的Api调用。

[![Go Version](https://img.shields.io/badge/Go-1.20+-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

---

## 项目简介

WXProxy是一个基于Go语言开发的微信 Api 代理服务，旨在提供更高效的微信Open Api调用，并支持分布式缓存，提供高性能和低延迟的服务。


---

## 主要功能

- **微信 Api 代理**
  提供微信 Open Api代理服务，支持最新的Open Api的调用。

- **Redis分布式缓存支持**
  使用Redis作为分布式缓存，支持高并发场景下的Token共享与访问。

- **GRPC接口**
  提供高效的GRPC接口，方便其他服务调用和集成。

- **健康检查**
  内置健康检查机制，确保服务的稳定性和可用性。

---

## 技术栈

- **编程语言**: Go (Golang)
- **缓存**: Redis
- **通信协议**: GRPC
- **日志**: Zap
- **依赖注入**: 自定义DI容器

---

## 快速开始

### 环境要求

- Go 1.20+
- Redis 6.x+
- Protobuf 编译器 (用于生成GRPC代码)

### 安装与运行

1. **克隆项目**
  ```bash
   git clone https://github.com/seth16888/wxproxy.git
   cd wxproxy
  ```

2. **安装依赖**
```bash
go mod download
```

3. **配置文件 在config.yaml中设置以下参数：**
```yaml
server:
  addr: 0.0.0.0:9000
  timeout: 15
log:
  level: debug
  filename: app.log
  max_size: 100
  max_age: 30
  max_backups: 3
  compress: false
redis:
  addr: 127.0.0.1:6379
  bd: 0
  password:
  username:
  read_timeout: 3
  write_timeout: 3
```

4. **构建pb文件**
```bash
./scripts/gen_pb.cmd
```

5. **编译**
```bash
./scripts/build.cmd
```

6. **运行**
```bash
./bin/wxproxy.exe -c conf/conf.yaml
```

## 健康检查
方法: Check

请求参数: 无

响应:
```json
{
  "status": "SERVING"
}
```

## 贡献指南
欢迎提交PR或Issue！以下是贡献步骤：

Fork 本仓库。

创建新分支 (git checkout -b feature/your-feature)。

提交更改 (git commit -m 'Add your feature')。

推送分支 (git push origin feature/your-feature)。

提交PR。

## 许可证
本项目采用 [MIT许可证](./LICENSE)。

## 联系我们
如有任何问题，请通过以下方式联系我们：

Email: [ucipl0094@hotmail.com](mailto:ucipl0094@hotmail.com)

GitHub Issues: [Issues Page](https://github.com/seth16888/wxproxy/issues)
