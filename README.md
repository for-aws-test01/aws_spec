# AWSomeShop - 内部员工福利电商系统

## 项目概述

AWSomeShop 是一个内部员工福利电商网站，员工可以使用"AWSome积分"浏览和兑换预选产品。

**当前版本**: v2.0 激进简化版  
**状态**: 🚧 开发准备阶段  
**工作量减少**: 65%（从141天减少到50天）

## 快速链接

- 📊 [项目状态](.kiro/PROJECT_STATUS.md) - 当前进度和统计
- 📋 [需求文档](.kiro/specs/requirements.md) - 功能需求
- 🎨 [设计文档](.kiro/specs/design.md) - 技术设计
- ✅ [任务文档](.kiro/specs/tasks.md) - 开发任务
- 📝 [变更日志](.kiro/specs/CHANGELOG.md) - 版本变更

## 技术栈

### 后端
- Go 1.21+ + Gin + GORM + MySQL 8.0+
- JWT 认证 + bcrypt 密码

### 前端
- 原生 JavaScript + Webpack 5 + Express
- Ant Design CSS + History API

## 项目结构

```
awsomeshop/
├── backend/              # Go 后端服务
│   ├── cmd/             # 应用入口
│   ├── internal/        # 业务模块
│   ├── pkg/             # 公共包
│   └── migrations/      # 数据库迁移
├── frontend/            # 前端应用
│   ├── src/            # 源代码
│   └── server.js       # Express 服务器
└── .kiro/              # 项目文档
    ├── specs/          # 规格文档
    └── PROJECT_STATUS.md
```

## 快速开始

### 1. 初始化数据库

```bash
cd awsomeshop/backend/migrations
chmod +x setup_simplified.sh
./setup_simplified.sh
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

### 员工端（7个功能）
- 登录/登出、浏览产品、兑换产品
- 查看订单、查看积分、个人信息

### 管理员端（10个功能）
- 登录/登出、管理员工、管理产品
- 管理积分、审核订单

## 简化说明

本版本删除了15个非核心功能：
- ❌ 审计日志、批量操作、软删除、统计报表
- ❌ 图片上传、单元测试、修改密码、编辑功能
- ❌ 列表筛选、产品上架/下架、取消订单等

详见 [变更日志](.kiro/specs/CHANGELOG.md)

## 开发进度

- ✅ 项目初始化（3/3 任务）
- ⏳ 后端基础设施（0/5 任务）
- ⏳ 用户认证（0/2 任务）
- ⏳ 其他模块（0/59 任务）

**总进度**: 3/69 任务（4%）

## 初始账号

- **工号**: admin
- **密码**: 123456
- **角色**: 管理员

## 文档

- [后端说明](awsomeshop/backend/README_SIMPLIFIED.md)
- [数据库说明](awsomeshop/backend/migrations/README_SIMPLIFIED.md)
- [需求规划](.kiro/requirements_plan.md)
- [设计规划](.kiro/design_plan.md)
- [任务规划](.kiro/tasks_plan.md)

## 许可证

内部项目，仅供公司内部使用。
