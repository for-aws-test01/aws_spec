# AWSomeShop - 内部员工福利电商系统

## 项目简介

AWSomeShop 是一个内部员工福利电商系统，允许员工使用"AWSome积分"兑换预选产品。系统采用前后端分离架构，提供员工端和管理员端两套功能界面。

### 主要功能

**员工端**
- 浏览和兑换产品
- 查看积分余额和明细
- 管理兑换订单
- 修改个人密码

**管理员端**
- 员工管理（创建、编辑、重置密码）
- 产品管理（CRUD、上下架）
- 积分管理（发放、扣除、批量发放）
- 订单审核（核销、拒绝）
- 系统概况统计

## 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **认证**: JWT
- **密码加密**: bcrypt

### 前端
- **语言**: JavaScript (ES6+)
- **UI 框架**: Ant Design CSS
- **构建工具**: Webpack 5
- **服务器**: Express
- **路由**: History API

### 开发工具
- **版本控制**: Git
- **容器化**: Docker Compose
- **测试框架**: Go testing + testify
- **API 文档**: Swagger

## 项目结构

```
awsomeshop/
├── backend/                        # 后端项目（Go）
│   ├── cmd/server/                 # 应用入口
│   ├── internal/                   # 内部模块
│   │   ├── auth/                   # 认证模块
│   │   ├── user/                   # 用户模块
│   │   ├── product/                # 产品模块
│   │   ├── point/                  # 积分模块
│   │   ├── order/                  # 订单模块
│   │   ├── audit/                  # 审计日志
│   │   ├── logger/                 # 应用日志
│   │   └── common/                 # 公共模块
│   ├── pkg/                        # 公共包
│   │   ├── database/               # 数据库连接
│   │   └── config/                 # 配置管理
│   ├── migrations/                 # 数据库迁移脚本
│   ├── uploads/                    # 文件上传目录
│   └── README.md                   # 后端说明文档
├── frontend/                       # 前端项目（JavaScript）
│   ├── src/                        # 源代码
│   │   ├── js/                     # JavaScript 文件
│   │   ├── css/                    # 样式文件
│   │   └── index.html              # HTML 模板
│   ├── dist/                       # 构建产物
│   ├── server.js                   # Express 服务器
│   ├── webpack.config.js           # Webpack 配置
│   └── README.md                   # 前端说明文档
├── docker-compose.yml              # Docker Compose 配置
└── README.md                       # 项目总体说明（本文件）
```

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 16+
- MySQL 8.0+
- Docker & Docker Compose（可选）

### 使用 Docker Compose（推荐）

1. **克隆项目**
```bash
git clone <repository-url>
cd awsomeshop
```

2. **启动所有服务**
```bash
docker-compose up -d
```

3. **访问应用**
- 前端: http://localhost:8000
- 后端 API: http://localhost:8080
- 默认管理员账号: admin / Admin@123

### 本地开发

#### 后端开发

1. **安装依赖**
```bash
cd backend
go mod download
```

2. **配置环境变量**
```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库连接等信息
```

3. **初始化数据库**
```bash
mysql -u root -p < migrations/schema.sql
mysql -u root -p < migrations/seed.sql
```

4. **启动后端服务**
```bash
go run cmd/server/main.go
```

后端服务将在 http://localhost:8080 启动

#### 前端开发

1. **安装依赖**
```bash
cd frontend
npm install
```

2. **开发模式**
```bash
npm run dev
```

3. **生产构建**
```bash
npm run build
npm start
```

前端服务将在 http://localhost:8000 启动

## API 文档

启动后端服务后，访问 Swagger API 文档：
- http://localhost:8080/swagger/index.html

## 测试

### 后端测试

```bash
cd backend
# 运行所有测试
go test ./...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 前端测试

```bash
cd frontend
npm test
```

## 默认账号

系统初始化后会创建以下默认账号：

| 角色 | 工号 | 密码 | 说明 |
|------|------|------|------|
| 管理员 | admin | Admin@123 | 系统管理员账号 |

## 核心功能模块

### 1. 用户认证
- JWT Token 认证
- 基于角色的访问控制（员工/管理员）
- 密码加密存储（bcrypt）

### 2. 用户管理
- 自动生成工号（格式：YYYYMMDD-NNN）
- 员工信息管理
- 密码重置

### 3. 产品管理
- 产品 CRUD 操作
- 图片上传（最大 1MB）
- 产品上架/下架
- 软删除

### 4. 积分管理
- 积分发放/扣除
- 批量发放（按部门/岗位）
- 积分日志记录
- 入职自动发放 1000 积分

### 5. 订单管理
- 产品兑换（扣除积分）
- 订单审核（核销/拒绝）
- 订单取消（退回积分）
- 产品快照（JSON 存储）

### 6. 审计日志
- 记录所有管理员操作
- 操作前后数据对比
- 完整的审计追踪

## 开发规范

### Git 提交规范

```
<type>: <subject>

<body>

<footer>
```

**类型（type）**:
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建/工具相关

### 分支策略

- `main`: 主分支，稳定版本
- `develop`: 开发分支
- `feature/*`: 功能分支
- `bugfix/*`: 修复分支

## 规格说明文档

本项目包含完整的规格说明文档，位于 `.kiro/specs/` 目录：

- **requirements.md** - 需求规格说明（EARS 模式，46 个需求）
- **design.md** - 技术设计方案（架构、数据库、API、前端）
- **architecture.md** - C4 模型架构图（Mermaid）
- **tasks.md** - 任务分解（85 个任务，14 个阶段）
- **tasks-parallel.md** - 并行任务执行计划（优化工期 54%）

### 项目统计

- **总任务数**: 85 个
- **API 接口数**: 40+
- **数据库表数**: 6 张
- **前端页面数**: 18 个
- **预估工期**: 13 周（并行优化后）
- **测试覆盖率目标**: 60%

## 贡献指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'feat: Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

本项目仅供内部使用。

## 联系方式

如有问题，请联系项目维护团队。

---

**项目状态**: 开发中 🚧  
**当前版本**: v1.0.0  
**最后更新**: 2026-01-14
