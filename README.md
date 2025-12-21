# 🔥 Burn-Text

基于 **Go (Gin)** + **Redis** 构建的企业级“阅后即焚”信息分享系统。
数据被读取一次后，立即物理销毁。

## ✨ 核心特性

* **⚡ 阅后即焚**：基于 Redis 原子操作，读后即删，不仅是标记删除。

* **🔒 零信任加密**：AES-GCM (256-bit) 加密，密钥由 URL 持有，服务器不存密钥。

* **🔑 访问密码**：(V1.1) 支持设置二次验证密码，防止爬虫误触销毁。

* **🛡️ 生产级架构**：内置 IP 限流中间件、优雅关机 (Graceful Shutdown)、Zap 结构化日志。

* **🐳 一键部署**：原生支持 Docker Compose。

## 🚀 快速开始

无需配置 Go 环境，需安装 Docker。

```bash
# 1. 克隆项目
git clone [https://github.com/你的用户名/burn-text.git](https://github.com/你的用户名/burn-text.git)
cd burn-text

# 2. 启动 (后台运行)
docker-compose up -d

# 3. 访问
# 浏览器打开 http://localhost:8080