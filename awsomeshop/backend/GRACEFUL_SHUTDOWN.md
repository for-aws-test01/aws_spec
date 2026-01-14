# 优雅关闭实现文档

## 概述

本文档描述了 AWSomeShop 后端服务的优雅关闭（Graceful Shutdown）实现。优雅关闭确保在服务器关闭时，所有正在进行的请求都能完成处理，并且所有资源（如数据库连接）都能正确释放。

## 实现位置

### 1. 数据库连接关闭

**文件**: `pkg/database/mysql.go`

**函数**: `Close() error`

```go
// Close 关闭数据库连接（优雅关闭）
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	log.Println("Database connection closed")
	return nil
}
```

**功能**:
- 检查数据库连接是否为 nil，避免空指针错误
- 获取底层的 `sql.DB` 实例
- 调用 `Close()` 方法关闭所有数据库连接
- 记录关闭日志
- 返回错误信息（如果有）

### 2. HTTP 服务器优雅关闭

**文件**: `cmd/server/main.go`

**实现步骤**:

1. **创建信号通道**
   ```go
   quit := make(chan os.Signal, 1)
   signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
   ```
   - 创建一个缓冲通道来接收操作系统信号
   - 监听 `SIGINT` (Ctrl+C) 和 `SIGTERM` (kill 命令) 信号

2. **等待关闭信号**
   ```go
   <-quit
   log.Println("Shutting down server...")
   ```
   - 阻塞等待关闭信号
   - 收到信号后开始关闭流程

3. **创建超时上下文**
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()
   ```
   - 创建一个 5 秒超时的 context
   - 确保关闭流程不会无限期等待

4. **关闭 HTTP 服务器**
   ```go
   if err := srv.Shutdown(ctx); err != nil {
       log.Printf("Server forced to shutdown: %v", err)
   }
   ```
   - 调用 `Shutdown()` 方法优雅关闭服务器
   - 等待所有活动连接完成或超时
   - 不再接受新的连接

5. **关闭数据库连接**
   ```go
   if err := database.Close(); err != nil {
       log.Printf("Error closing database: %v", err)
   }
   ```
   - 关闭数据库连接池
   - 释放所有数据库资源

6. **退出程序**
   ```go
   log.Println("Server exited")
   ```
   - 记录退出日志
   - 程序正常退出

## 工作流程

```
1. 服务器正常运行
   ↓
2. 收到 SIGINT/SIGTERM 信号
   ↓
3. 停止接受新连接
   ↓
4. 等待现有请求完成（最多 5 秒）
   ↓
5. 关闭 HTTP 服务器
   ↓
6. 关闭数据库连接
   ↓
7. 程序退出
```

## 测试

### 单元测试

**文件**: `pkg/database/mysql_graceful_test.go`

包含以下测试用例：

1. **TestGracefulShutdown**: 测试完整的优雅关闭流程
   - 连接数据库
   - 验证连接池配置
   - 调用 Close() 关闭连接
   - 验证连接已关闭

2. **TestMultipleClose**: 测试多次关闭的安全性
   - 第一次关闭应该成功
   - 第二次关闭应该安全处理（不崩溃）

3. **TestCloseWithNilDB**: 测试 nil DB 的处理
   - 当 DB 为 nil 时调用 Close()
   - 应该返回 nil 而不是错误

### 运行测试

```bash
cd awsomeshop/backend
go test -v ./pkg/database
```

### 手动测试

1. 启动服务器：
   ```bash
   cd awsomeshop/backend
   go run ./cmd/server/main.go
   ```

2. 发送关闭信号：
   - 按 `Ctrl+C`
   - 或在另一个终端执行: `kill -TERM <pid>`

3. 观察日志输出：
   ```
   Server starting on port 8080
   Database connected successfully
   ^CShutting down server...
   Database connection closed
   Server exited
   ```

## 配置

### 超时时间

当前关闭超时时间设置为 5 秒，可以在 `cmd/server/main.go` 中修改：

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
```

### 连接池配置

数据库连接池配置在 `pkg/database/mysql.go` 中：

```go
sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期
```

## 最佳实践

1. **始终使用优雅关闭**: 避免强制终止进程，使用 SIGINT 或 SIGTERM 信号
2. **设置合理的超时时间**: 根据业务需求调整超时时间
3. **记录关闭日志**: 便于排查问题和监控
4. **按顺序关闭资源**: 先关闭 HTTP 服务器，再关闭数据库连接
5. **处理关闭错误**: 记录但不阻止关闭流程

## 注意事项

1. **长时间运行的请求**: 如果请求处理时间超过超时时间，会被强制终止
2. **数据库事务**: 确保事务在关闭前提交或回滚
3. **后台任务**: 如果有后台 goroutine，需要单独处理其关闭
4. **文件句柄**: 确保所有打开的文件都被正确关闭

## 验收标准

根据 TASK-004 的验收标准：

- [x] 实现 `pkg/database/mysql.go` 数据库连接
- [x] 实现 `pkg/config/config.go` 配置管理
- [x] 从 `.env` 文件读取配置
- [x] 实现数据库连接池
- [x] 实现优雅关闭

所有验收标准已完成。

## 相关文件

- `pkg/database/mysql.go` - 数据库连接和关闭实现
- `cmd/server/main.go` - 服务器启动和优雅关闭实现
- `pkg/database/mysql_graceful_test.go` - 优雅关闭测试
- `pkg/database/mysql_test.go` - 数据库连接测试
- `.env.example` - 环境变量配置示例

## 参考资料

- [Go HTTP Server Graceful Shutdown](https://pkg.go.dev/net/http#Server.Shutdown)
- [GORM Connection Pool](https://gorm.io/docs/generic_interface.html)
- [Go Signal Handling](https://pkg.go.dev/os/signal)
