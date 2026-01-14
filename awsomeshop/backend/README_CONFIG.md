# 配置管理说明

## 概述

AWSomeShop 后端使用环境变量进行配置管理，支持从 `.env` 文件或系统环境变量加载配置。

## 快速开始

### 1. 创建配置文件

```bash
# 复制示例配置文件
cp .env.example .env

# 编辑配置文件
vim .env
```

### 2. 配置示例

```env
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=awsomeshop

# JWT 配置
JWT_SECRET=your-secret-key-change-this-in-production

# 服务器配置
SERVER_PORT=8080

# 文件上传配置
UPLOAD_DIR=./uploads/image/
MAX_UPLOAD_SIZE=1048576

# 日志配置
LOG_LEVEL=INFO
```

### 3. 在代码中使用

```go
package main

import (
    "awsomeshop/backend/pkg/config"
    "log"
)

func main() {
    // 加载配置
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // 使用配置
    log.Printf("Server starting on port %s", cfg.ServerPort)
    log.Printf("Connecting to database at %s:%s", cfg.DBHost, cfg.DBPort)
}
```

## 配置项说明

### 数据库配置

| 配置项 | 环境变量 | 默认值 | 必需 | 说明 |
|--------|---------|--------|------|------|
| 主机地址 | DB_HOST | localhost | 是 | MySQL 服务器地址 |
| 端口 | DB_PORT | 3306 | 是 | MySQL 服务器端口 |
| 用户名 | DB_USER | root | 是 | 数据库用户名 |
| 密码 | DB_PASSWORD | (空) | 否 | 数据库密码 |
| 数据库名 | DB_NAME | awsomeshop | 是 | 数据库名称 |

### JWT 配置

| 配置项 | 环境变量 | 默认值 | 必需 | 说明 |
|--------|---------|--------|------|------|
| 密钥 | JWT_SECRET | (无) | 是 | JWT 签名密钥，生产环境必须设置强密钥 |

### 服务器配置

| 配置项 | 环境变量 | 默认值 | 必需 | 说明 |
|--------|---------|--------|------|------|
| 端口 | SERVER_PORT | 8080 | 是 | HTTP 服务监听端口 |

### 文件上传配置

| 配置项 | 环境变量 | 默认值 | 必需 | 说明 |
|--------|---------|--------|------|------|
| 上传目录 | UPLOAD_DIR | ./uploads/image/ | 否 | 文件上传保存路径 |
| 最大大小 | MAX_UPLOAD_SIZE | 1048576 | 否 | 最大上传文件大小（字节），默认 1MB |

### 日志配置

| 配置项 | 环境变量 | 默认值 | 必需 | 说明 |
|--------|---------|--------|------|------|
| 日志级别 | LOG_LEVEL | INFO | 否 | 日志级别：DEBUG, INFO, WARNING, ERROR |

## 部署说明

### 开发环境

使用 `.env` 文件：

```bash
# 1. 创建 .env 文件
cp .env.example .env

# 2. 修改配置
vim .env

# 3. 运行应用
go run cmd/server/main.go
```

### 生产环境

使用系统环境变量：

```bash
# 设置环境变量
export DB_HOST=prod-db-host
export DB_PORT=3306
export DB_USER=prod_user
export DB_PASSWORD=secure_password
export DB_NAME=awsomeshop
export JWT_SECRET=very-secure-secret-key
export SERVER_PORT=8080

# 运行应用
./awsomeshop-server
```

或使用 Docker：

```dockerfile
# Dockerfile
FROM golang:1.18-alpine

WORKDIR /app
COPY . .

RUN go build -o server cmd/server/main.go

# 通过环境变量传递配置
ENV DB_HOST=db
ENV DB_PORT=3306
ENV JWT_SECRET=change-me

CMD ["./server"]
```

```bash
# docker-compose.yml
version: '3.8'
services:
  app:
    build: .
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=password
      - DB_NAME=awsomeshop
      - JWT_SECRET=your-secret-key
      - SERVER_PORT=8080
    ports:
      - "8080:8080"
```

## 安全建议

1. **不要提交 .env 文件到版本控制**
   - `.env` 文件已添加到 `.gitignore`
   - 只提交 `.env.example` 作为模板

2. **使用强 JWT 密钥**
   ```bash
   # 生成随机密钥
   openssl rand -base64 32
   ```

3. **生产环境配置**
   - 使用环境变量而不是 .env 文件
   - 使用密钥管理服务（如 AWS Secrets Manager）
   - 定期轮换密钥

4. **数据库密码**
   - 使用强密码
   - 不要使用默认密码
   - 限制数据库访问权限

## 故障排查

### 配置加载失败

```
Error: JWT_SECRET is required
```

**解决方案**：确保设置了 JWT_SECRET 环境变量

```bash
export JWT_SECRET=your-secret-key
```

### 数据库连接失败

```
Error: DB_HOST is required
```

**解决方案**：检查数据库配置是否正确

```bash
# 检查环境变量
echo $DB_HOST
echo $DB_PORT

# 或检查 .env 文件
cat .env
```

### .env 文件未加载

**可能原因**：
1. .env 文件不在当前工作目录
2. .env 文件格式错误

**解决方案**：
```bash
# 确保 .env 文件在正确位置
ls -la .env

# 检查文件格式（不要有多余空格）
cat .env
```

## API 参考

### config.Load()

加载配置并返回 Config 实例。

```go
func Load() (*Config, error)
```

**返回值**：
- `*Config`: 配置实例
- `error`: 如果配置验证失败，返回错误

**示例**：
```go
cfg, err := config.Load()
if err != nil {
    log.Fatal(err)
}
```

### config.GetConfig()

获取全局配置实例。

```go
func GetConfig() *Config
```

**返回值**：
- `*Config`: 全局配置实例

**示例**：
```go
cfg := config.GetConfig()
fmt.Println(cfg.ServerPort)
```

### Config.Validate()

验证配置的有效性。

```go
func (c *Config) Validate() error
```

**返回值**：
- `error`: 如果配置无效，返回错误

## 更多信息

- [godotenv 文档](https://github.com/joho/godotenv)
- [12-Factor App 配置原则](https://12factor.net/config)
