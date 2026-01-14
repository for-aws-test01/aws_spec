# AWSomeShop - 内部员工福利电商系统

## 项目概述

AWSomeShop 是一个内部员工福利电商网站，员工可以使用"AWSome积分"浏览和兑换预选产品。

**当前版本**: v3.0 超级简化版  
**状态**: 🚧 开发准备阶段  
**工作量减少**: 70%（从141天减少到49天）

## 快速链接

- 📊 [项目状态](.kiro/PROJECT_STATUS.md) - 当前进度和统计
- 📋 [需求文档](.kiro/specs/requirements.md) - 功能需求（v3.0）
- ✅ [任务文档](.kiro/specs/tasks-v3.md) - 开发任务（v3.0）
- 🗄️ [数据库架构](awsomeshop/backend/migrations/schema_v3.sql) - 数据库设计（v3.0）

## 技术栈

### 后端
- Go 1.21+ + Gin + GORM + MySQL 8.0+
- JWT 认证 + bcrypt 密码

### 前端
- 原生 JavaScript + Webpack 5 + Express
- 原生 CSS + History API

## 项目结构

```
awsomeshop/
├── backend/              # Go 后端服务
│   ├── cmd/             # 应用入口
│   ├── internal/        # 业务模块
│   ├── pkg/             # 公共包
│   └── migrations/      # 数据库迁移（v3.0）
├── frontend/            # 前端应用
│   ├── src/            # 源代码
│   └── server.js       # Express 服务器
└── .kiro/              # 项目文档
    ├── specs/          # 规格文档（v3.0）
    └── PROJECT_STATUS.md
```

## 快速开始

### 1. 初始化数据库

```bash
cd awsomeshop/backend/migrations
mysql -u root -p < schema_v3.sql
mysql -u root -p awsomeshop < seed_v3.sql
```

### 2. 启动后端（待开发）

```bash
cd awsomeshop/backend
cp .env.example .env
# 编辑 .env 配置数据库连接
go run cmd/server/main.go
```

### 3. 启动前端（待开发）

```bash
cd awsomeshop/frontend
npm install
npm start
```

## 核心功能

### 员工端（4个页面）
- 登录页、产品列表页、产品详情页、订单列表页
- 积分余额显示（导航栏）

### 管理员端（5个页面）
- 登录页、员工管理、产品管理、订单管理、积分发放

## 简化说明（v3.0）

本版本删除了35个任务，工作量减少70%：
- ❌ 所有单元测试
- ❌ 审计日志、应用日志系统
- ❌ 图片上传（改用URL）
- ❌ 系统概况、个人信息页
- ❌ 修改/重置密码
- ❌ 批量操作、扣除积分
- ❌ 积分日志查看
- ❌ 详情页（员工/产品/订单）
- ❌ 编辑功能
- ❌ 上架/下架、软删除
- ❌ 订单拒绝/取消
- ❌ Docker、Swagger、性能测试

详见 [任务文档](.kiro/specs/tasks-v3.md)

## 开发进度

- ✅ 项目初始化（3/3 任务）
- ⏳ 后端基础设施（2/4 任务）
- ⏳ 其他模块（0/27 任务）

**总进度**: 3/34 任务（9%）

## 数据库（4个表）

1. **users** - 用户表（员工+管理员）
2. **products** - 产品表
3. **orders** - 订单表
4. **point_logs** - 积分日志表

## 初始账号

- **工号**: admin
- **密码**: 123456
- **角色**: 管理员

## 版本对比

| 版本 | 任务数 | 工作量 | 减少比例 |
|------|--------|--------|----------|
| v1.0 原始版 | 85 | 141天 | - |
| v2.0 激进简化 | 69 | 50天 | 65% |
| v3.0 超级简化 | 34 | 49天 | 70% |

## 许可证

内部项目，仅供公司内部使用。
