# TASK-004 子任务完成报告：从 `.env` 文件读取配置

## 任务状态：✅ 已完成

## 任务描述
实现从 `.env` 文件读取配置的功能，使应用能够在开发环境中使用 `.env` 文件管理配置，在生产环境中使用系统环境变量。

## 实现详情

### 1. 依赖包已安装
在 `go.mod` 文件中已包含 `godotenv` 包：
```go
github.com/joho/godotenv v1.5.1 // indirect
```

### 2. 核心实现代码
在 `pkg/config/config.go` 的 `Load()` 函数中（第 42 行）：

```go
func Load() (*Config, error) {
    // 尝试加载 .env 文件（如果存在）
    // 在生产环境中，可能不存在 .env 文件，直接使用系统环境变量
    _ = godotenv.Load()
    
    config := &Config{
        // 数据库配置
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "3306"),
        DBUser:     getEnv("DB_USER", "root"),
        DBPassword: getEnv("DB_PASSWORD", ""),
        DBName:     getEnv("DB_NAME", "awsomeshop"),
        
        // JWT 配置
        JWTSecret: getEnv("JWT_SECRET", ""),
        
        // 服务器配置
        ServerPort: getEnv("SERVER_PORT", "8080"),
        
        // 文件上传配置
        UploadDir:     getEnv("UPLOAD_DIR", "./uploads/image/"),
        MaxUploadSize: getEnvAsInt64("MAX_UPLOAD_SIZE", 1048576),
        
        // 日志配置
        LogLevel: getEnv("LOG_LEVEL", "INFO"),
    }
    
    // 验证必需的配置项
    if err := config.Validate(); err != nil {
        return nil, err
    }
    
    AppConfig = config
    return config, nil
}
```

### 3. 工作原理

#### 3.1 加载流程
1. **调用 `godotenv.Load()`**：
   - 从当前目录读取 `.env` 文件
   - 将文件中的键值对加载到进程的环境变量中
   - 错误被有意忽略（`_ = godotenv.Load()`），实现优雅降级

2. **读取环境变量**：
   - 使用 `os.Getenv()` 读取环境变量
   - 如果环境变量不存在，使用默认值
   - 支持字符串和整数类型的配置

3. **配置验证**：
   - 验证必需的配置项（如 JWT_SECRET）
   - 如果验证失败，返回错误

#### 3.2 优雅降级设计
```go
_ = godotenv.Load()  // 忽略错误
```

这种设计的好处：
- **开发环境**：存在 `.env` 文件，配置从文件加载
- **生产环境**：不存在 `.env` 文件，配置从系统环境变量加载
- **灵活性**：支持两种配置方式，适应不同部署场景

### 4. 配置文件示例

#### `.env.example` 文件内容：
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=awsomeshop

# JWT Configuration
JWT_SECRET=your-secret-key-here

# Server Configuration
SERVER_PORT=8080

# File Upload Configuration
UPLOAD_DIR=/uploads/image/
MAX_UPLOAD_SIZE=1048576

# Log Configuration
LOG_LEVEL=INFO
```

### 5. 使用方法

#### 5.1 开发环境使用
```bash
# 1. 复制示例文件
cp .env.example .env

# 2. 编辑 .env 文件，填入实际配置
vim .env

# 3. 运行应用，配置将自动从 .env 加载
go run main.go
```

#### 5.2 代码中使用
```go
package main

import (
    "awsomeshop/backend/pkg/config"
    "log"
)

func main() {
    // 加载配置（自动从 .env 或环境变量读取）
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("配置加载失败: %v", err)
    }
    
    // 使用配置
    log.Printf("服务器端口: %s", cfg.ServerPort)
    log.Printf("数据库主机: %s", cfg.DBHost)
    
    // 或使用全局实例
    globalCfg := config.GetConfig()
    log.Printf("上传目录: %s", globalCfg.UploadDir)
}
```

### 6. 测试覆盖

已实现的测试用例（`pkg/config/config_test.go`）：

1. **TestLoad()**：测试从环境变量加载配置
2. **TestLoadWithDefaults()**：测试默认值
3. **TestValidate()**：测试配置验证
4. **TestGetConfig()**：测试全局配置实例

### 7. 支持的配置项

| 配置项 | 环境变量名 | 默认值 | 说明 |
|--------|-----------|--------|------|
| 数据库主机 | DB_HOST | localhost | MySQL 主机地址 |
| 数据库端口 | DB_PORT | 3306 | MySQL 端口 |
| 数据库用户 | DB_USER | root | MySQL 用户名 |
| 数据库密码 | DB_PASSWORD | (空) | MySQL 密码 |
| 数据库名称 | DB_NAME | awsomeshop | 数据库名 |
| JWT 密钥 | JWT_SECRET | (必需) | JWT 签名密钥 |
| 服务器端口 | SERVER_PORT | 8080 | HTTP 服务端口 |
| 上传目录 | UPLOAD_DIR | ./uploads/image/ | 文件上传路径 |
| 最大上传大小 | MAX_UPLOAD_SIZE | 1048576 | 最大文件大小（字节） |
| 日志级别 | LOG_LEVEL | INFO | 日志级别 |

### 8. 验证步骤

要验证 .env 文件加载功能：

```bash
# 1. 创建测试 .env 文件
cat > .env << EOF
DB_HOST=test-host
DB_PORT=3307
DB_USER=test-user
DB_PASSWORD=test-pass
DB_NAME=test-db
JWT_SECRET=test-secret
SERVER_PORT=9090
EOF

# 2. 运行应用或测试
go run main.go

# 3. 检查日志输出，确认配置已从 .env 加载
```

## 技术亮点

1. **使用行业标准库**：`godotenv` 是 Go 生态中最流行的 .env 文件加载库
2. **优雅降级**：支持 .env 文件和系统环境变量两种方式
3. **类型安全**：提供字符串和整数类型的配置读取
4. **默认值支持**：所有配置项都有合理的默认值
5. **配置验证**：确保必需的配置项存在
6. **全局访问**：提供全局配置实例，方便在应用各处使用
7. **完整测试**：包含单元测试，确保功能正确性

## 结论

✅ **任务已完成**

从 `.env` 文件读取配置的功能已经完整实现，包括：
- ✅ 依赖包已安装（godotenv）
- ✅ 核心代码已实现（godotenv.Load()）
- ✅ 支持优雅降级（.env 文件可选）
- ✅ 提供默认值
- ✅ 配置验证
- ✅ 单元测试覆盖

**无需额外的代码修改**，此任务可以标记为完成。

## 下一步

建议继续完成 TASK-004 的其他子任务：
- [ ] 实现数据库连接池
- [ ] 实现优雅关闭
