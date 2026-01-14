# AWSomeShop 系统架构图（C4 Model）

## 文档信息

- **项目名称**: AWSomeShop 内部员工福利电商系统
- **文档版本**: 1.0
- **创建日期**: 2026-01-14
- **架构模型**: C4 Model
- **图表工具**: Mermaid

---

## Level 1: System Context Diagram（系统上下文图）

系统上下文图展示了 AWSomeShop 系统与外部用户和系统的关系。

```mermaid
C4Context
    title System Context Diagram - AWSomeShop 内部员工福利电商系统

    Person(employee, "员工", "内部员工，使用积分兑换产品")
    Person(admin, "管理员", "系统管理员，管理产品、员工和积分")

    System(awsomeshop, "AWSomeShop 系统", "内部员工福利电商平台，提供积分兑换功能")

    Rel(employee, awsomeshop, "浏览产品、兑换产品、查看积分", "HTTPS")
    Rel(admin, awsomeshop, "管理员工、管理产品、管理积分、审核订单", "HTTPS")

    UpdateLayoutConfig($c4ShapeInRow="2", $c4BoundaryInRow="1")
```

**说明**:
- **员工**: 使用系统浏览产品、兑换产品、查看积分余额和兑换历史
- **管理员**: 管理系统的所有资源（员工、产品、积分、订单）
- **AWSomeShop 系统**: 核心业务系统，处理所有业务逻辑

---

## Level 2: Container Diagram（容器图）

容器图展示了系统内部的高层技术构建块（应用程序、数据存储等）。

```mermaid
C4Container
    title Container Diagram - AWSomeShop 系统容器

    Person(employee, "员工", "内部员工")
    Person(admin, "管理员", "系统管理员")

    Container_Boundary(awsomeshop, "AWSomeShop 系统") {
        Container(frontend, "前端应用", "Node.js, Express, JavaScript", "提供用户界面，托管静态文件")
        Container(backend, "后端 API", "Go, Gin, GORM", "处理业务逻辑，提供 RESTful API")
        ContainerDb(database, "数据库", "MySQL 8.0", "存储用户、产品、订单、积分日志等数据")
        Container(filesystem, "文件系统", "本地存储", "存储产品图片")
    }

    Rel(employee, frontend, "访问", "HTTPS, 端口 8000")
    Rel(admin, frontend, "访问", "HTTPS, 端口 8000")
    
    Rel(frontend, backend, "调用 API", "HTTP, 端口 8080, JSON/REST")
    Rel(backend, database, "读写数据", "MySQL Protocol, 端口 3306")
    Rel(backend, filesystem, "读写文件", "文件 I/O")

    UpdateLayoutConfig($c4ShapeInRow="2", $c4BoundaryInRow="1")
```

**说明**:
- **前端应用**: Express 服务器托管 Webpack 打包的静态文件（HTML/CSS/JS）
- **后端 API**: Go + Gin 框架，提供 RESTful API，处理所有业务逻辑
- **数据库**: MySQL 存储所有业务数据
- **文件系统**: 本地存储产品图片（/uploads/image/）

---

## Level 3: Component Diagram - 后端 API（组件图）

组件图展示了后端 API 容器内部的组件结构。

```mermaid
C4Component
    title Component Diagram - 后端 API 组件

    Container_Boundary(backend, "后端 API (Go + Gin)") {
        Component(api_gateway, "API 网关层", "Gin Router", "路由分发、中间件处理")
        
        Component(auth_middleware, "认证中间件", "JWT", "验证 Token、权限检查")
        Component(cors_middleware, "CORS 中间件", "Gin CORS", "处理跨域请求")
        Component(log_middleware, "日志中间件", "Logger", "记录请求日志")
        
        Component(auth_handler, "认证处理器", "Handler", "登录、登出、修改密码")
        Component(user_handler, "用户处理器", "Handler", "员工和管理员管理")
        Component(product_handler, "产品处理器", "Handler", "产品 CRUD、上下架")
        Component(point_handler, "积分处理器", "Handler", "积分发放、扣除、查询")
        Component(order_handler, "订单处理器", "Handler", "兑换、取消、审核")
        Component(dashboard_handler, "概况处理器", "Handler", "系统统计数据")
        
        Component(user_service, "用户服务", "Service", "用户业务逻辑")
        Component(product_service, "产品服务", "Service", "产品业务逻辑")
        Component(point_service, "积分服务", "Service", "积分业务逻辑")
        Component(order_service, "订单服务", "Service", "订单业务逻辑")
        Component(audit_service, "审计服务", "Service", "审计日志记录")
        Component(logger_service, "日志服务", "Service", "应用日志记录")
        
        Component(user_repo, "用户仓储", "Repository", "用户数据访问")
        Component(product_repo, "产品仓储", "Repository", "产品数据访问")
        Component(point_repo, "积分仓储", "Repository", "积分数据访问")
        Component(order_repo, "订单仓储", "Repository", "订单数据访问")
        Component(audit_repo, "审计仓储", "Repository", "审计日志数据访问")
        Component(log_repo, "日志仓储", "Repository", "应用日志数据访问")
    }

    ContainerDb(database, "MySQL 数据库", "MySQL 8.0")
    Container(filesystem, "文件系统", "本地存储")

    Rel(api_gateway, auth_middleware, "使用")
    Rel(api_gateway, cors_middleware, "使用")
    Rel(api_gateway, log_middleware, "使用")
    
    Rel(api_gateway, auth_handler, "路由到")
    Rel(api_gateway, user_handler, "路由到")
    Rel(api_gateway, product_handler, "路由到")
    Rel(api_gateway, point_handler, "路由到")
    Rel(api_gateway, order_handler, "路由到")
    Rel(api_gateway, dashboard_handler, "路由到")
    
    Rel(auth_handler, user_service, "调用")
    Rel(user_handler, user_service, "调用")
    Rel(product_handler, product_service, "调用")
    Rel(point_handler, point_service, "调用")
    Rel(order_handler, order_service, "调用")
    
    Rel(user_service, audit_service, "记录审计")
    Rel(product_service, audit_service, "记录审计")
    Rel(point_service, audit_service, "记录审计")
    Rel(order_service, audit_service, "记录审计")
    
    Rel(user_service, user_repo, "调用")
    Rel(product_service, product_repo, "调用")
    Rel(point_service, point_repo, "调用")
    Rel(order_service, order_repo, "调用")
    Rel(audit_service, audit_repo, "调用")
    Rel(logger_service, log_repo, "调用")
    
    Rel(user_repo, database, "读写", "SQL")
    Rel(product_repo, database, "读写", "SQL")
    Rel(point_repo, database, "读写", "SQL")
    Rel(order_repo, database, "读写", "SQL")
    Rel(audit_repo, database, "写入", "SQL")
    Rel(log_repo, database, "写入", "SQL")
    
    Rel(product_service, filesystem, "读写图片", "File I/O")

    UpdateLayoutConfig($c4ShapeInRow="3", $c4BoundaryInRow="1")
```

**组件说明**:

### API 网关层
- **API 网关**: Gin Router，负责路由分发和中间件处理

### 中间件层
- **认证中间件**: JWT Token 验证和权限检查
- **CORS 中间件**: 处理跨域请求
- **日志中间件**: 记录所有 API 请求

### 处理器层（Handler）
- **认证处理器**: 处理登录、登出、修改密码
- **用户处理器**: 处理员工和管理员的 CRUD 操作
- **产品处理器**: 处理产品的 CRUD、上下架操作
- **积分处理器**: 处理积分发放、扣除、查询
- **订单处理器**: 处理产品兑换、订单取消、订单审核
- **概况处理器**: 提供系统统计数据

### 服务层（Service）
- **用户服务**: 用户相关业务逻辑（工号生成、密码哈希等）
- **产品服务**: 产品相关业务逻辑（图片上传、软删除等）
- **积分服务**: 积分相关业务逻辑（批量发放、日志记录等）
- **订单服务**: 订单相关业务逻辑（订单号生成、事务处理等）
- **审计服务**: 记录所有管理员操作
- **日志服务**: 记录应用运行日志

### 仓储层（Repository）
- **用户仓储**: 用户数据的 CRUD 操作
- **产品仓储**: 产品数据的 CRUD 操作
- **积分仓储**: 积分和积分日志的 CRUD 操作
- **订单仓储**: 订单数据的 CRUD 操作
- **审计仓储**: 审计日志的写入操作
- **日志仓储**: 应用日志的写入操作

---

## Level 3: Component Diagram - 前端应用（组件图）

```mermaid
C4Component
    title Component Diagram - 前端应用组件

    Container_Boundary(frontend, "前端应用 (JavaScript)") {
        Component(express_server, "Express 服务器", "Node.js", "托管静态文件、支持 History 路由")
        
        Component(router, "路由管理器", "Router", "管理前端路由、页面导航")
        Component(api_client, "API 客户端", "Fetch API", "封装 HTTP 请求")
        Component(auth_manager, "认证管理器", "Auth", "管理登录状态、Token")
        Component(storage_manager, "存储管理器", "Storage", "管理 localStorage")
        
        Component(navbar, "导航栏组件", "Component", "顶部导航、用户信息")
        Component(toast, "提示组件", "Component", "Toast 消息提示")
        Component(pagination, "分页组件", "Component", "分页控件")
        
        Component(login_page, "登录页", "Page", "用户登录界面")
        Component(product_list_page, "产品列表页", "Page", "员工端产品浏览")
        Component(product_detail_page, "产品详情页", "Page", "产品详情和兑换")
        Component(profile_page, "个人中心页", "Page", "个人信息、积分、订单")
        Component(admin_dashboard_page, "系统概况页", "Page", "管理员首页")
        Component(admin_employee_page, "员工管理页", "Page", "员工列表和管理")
        Component(admin_product_page, "产品管理页", "Page", "产品列表和管理")
        Component(admin_point_page, "积分管理页", "Page", "积分发放和日志")
        Component(admin_order_page, "订单管理页", "Page", "订单列表和审核")
    }

    Container(backend, "后端 API", "Go, Gin")

    Rel(express_server, router, "提供")
    
    Rel(router, login_page, "路由到")
    Rel(router, product_list_page, "路由到")
    Rel(router, product_detail_page, "路由到")
    Rel(router, profile_page, "路由到")
    Rel(router, admin_dashboard_page, "路由到")
    Rel(router, admin_employee_page, "路由到")
    Rel(router, admin_product_page, "路由到")
    Rel(router, admin_point_page, "路由到")
    Rel(router, admin_order_page, "路由到")
    
    Rel(login_page, auth_manager, "使用")
    Rel(login_page, api_client, "使用")
    Rel(login_page, toast, "使用")
    
    Rel(product_list_page, api_client, "使用")
    Rel(product_list_page, pagination, "使用")
    Rel(product_list_page, navbar, "使用")
    
    Rel(auth_manager, storage_manager, "使用")
    Rel(api_client, backend, "调用 API", "HTTP/JSON")

    UpdateLayoutConfig($c4ShapeInRow="3", $c4BoundaryInRow="1")
```

**组件说明**:

### 服务器层
- **Express 服务器**: 托管静态文件，支持 History 路由

### 核心模块
- **路由管理器**: 管理前端路由，实现 SPA 导航
- **API 客户端**: 封装所有 HTTP 请求，统一错误处理
- **认证管理器**: 管理用户登录状态，处理 JWT Token
- **存储管理器**: 管理 localStorage 数据

### 公共组件
- **导航栏组件**: 顶部导航菜单、用户信息显示
- **提示组件**: Toast 消息提示
- **分页组件**: 列表分页控件

### 员工端页面
- **登录页**: 员工和管理员登录
- **产品列表页**: 浏览所有产品
- **产品详情页**: 查看产品详情、兑换产品
- **个人中心页**: 个人信息、积分余额、兑换历史

### 管理员端页面
- **系统概况页**: 系统统计数据
- **员工管理页**: 员工列表、创建、编辑
- **产品管理页**: 产品列表、创建、编辑、上下架
- **积分管理页**: 积分发放、扣除、日志查询
- **订单管理页**: 订单列表、核销、拒绝

---

## 数据流图

### 用户登录流程

```mermaid
sequenceDiagram
    actor User as 用户
    participant FE as 前端应用
    participant BE as 后端 API
    participant DB as MySQL 数据库

    User->>FE: 输入工号和密码
    FE->>BE: POST /api/auth/login
    BE->>DB: 查询用户信息
    DB-->>BE: 返回用户数据
    BE->>BE: 验证密码 (bcrypt)
    BE->>BE: 生成 JWT Token
    BE-->>FE: 返回用户信息 + Set-Cookie
    FE->>FE: 保存用户信息到 localStorage
    FE-->>User: 跳转到首页
```

### 产品兑换流程

```mermaid
sequenceDiagram
    actor Employee as 员工
    participant FE as 前端应用
    participant BE as 后端 API
    participant DB as MySQL 数据库

    Employee->>FE: 点击兑换按钮
    FE->>BE: POST /api/orders/redeem
    BE->>DB: 开始事务
    BE->>DB: 验证产品状态
    BE->>DB: 验证用户积分
    BE->>DB: 扣除积分
    BE->>DB: 创建订单
    BE->>DB: 记录积分日志
    BE->>DB: 提交事务
    DB-->>BE: 事务成功
    BE-->>FE: 返回订单信息
    FE->>FE: 更新积分余额
    FE-->>Employee: 显示兑换成功
```

### 订单审核流程

```mermaid
sequenceDiagram
    actor Admin as 管理员
    participant FE as 前端应用
    participant BE as 后端 API
    participant DB as MySQL 数据库

    Admin->>FE: 选择订单并核销/拒绝
    alt 核销订单
        FE->>BE: PUT /api/admin/orders/:id/approve
        BE->>DB: 更新订单状态为 approved
        BE->>DB: 记录审核信息
        BE->>DB: 记录审计日志
    else 拒绝订单
        FE->>BE: PUT /api/admin/orders/:id/reject
        BE->>DB: 开始事务
        BE->>DB: 更新订单状态为 rejected
        BE->>DB: 退回积分
        BE->>DB: 记录积分日志
        BE->>DB: 记录审计日志
        BE->>DB: 提交事务
    end
    DB-->>BE: 操作成功
    BE-->>FE: 返回结果
    FE-->>Admin: 显示操作成功
```

---

## 部署架构图

```mermaid
graph TB
    subgraph "本地服务器"
        subgraph "前端服务 - 端口 8000"
            FE[Express Server<br/>静态文件托管]
        end
        
        subgraph "后端服务 - 端口 8080"
            BE[Go Application<br/>Gin + GORM]
        end
        
        subgraph "数据库 - 端口 3306"
            DB[(MySQL 8.0<br/>awsomeshop)]
        end
        
        subgraph "文件存储"
            FS[/uploads/image/<br/>产品图片/]
        end
    end
    
    User[用户浏览器] -->|HTTP :8000| FE
    FE -->|HTTP :8080<br/>CORS| BE
    BE -->|MySQL Protocol| DB
    BE -->|File I/O| FS
    
    style User fill:#e1f5ff
    style FE fill:#b3e5fc
    style BE fill:#81c784
    style DB fill:#ffb74d
    style FS fill:#fff59d
```

---

## 数据库 ER 图

```mermaid
erDiagram
    USERS ||--o{ ORDERS : "creates"
    USERS ||--o{ POINT_LOGS : "has"
    USERS ||--o{ AUDIT_LOGS : "performs"
    PRODUCTS ||--o{ ORDERS : "ordered_in"
    ORDERS ||--o| POINT_LOGS : "generates"
    
    USERS {
        bigint id PK
        varchar employee_id UK
        varchar name
        varchar department
        varchar position
        varchar email
        varchar phone
        varchar password
        enum role
        enum status
        int points
        date hire_date
        timestamp created_at
        timestamp updated_at
    }
    
    PRODUCTS {
        bigint id PK
        varchar name
        text description
        varchar image_url
        int points_required
        enum status
        tinyint is_deleted
        timestamp created_at
        timestamp updated_at
    }
    
    ORDERS {
        bigint id PK
        varchar order_no UK
        bigint user_id FK
        bigint product_id FK
        json product_snapshot
        int points_cost
        enum status
        timestamp applied_at
        timestamp reviewed_at
        bigint reviewer_id FK
        text reject_reason
        text approval_note
        timestamp created_at
        timestamp updated_at
    }
    
    POINT_LOGS {
        bigint id PK
        bigint user_id FK
        bigint operator_id FK
        int amount
        text reason
        enum type
        bigint order_id FK
        timestamp created_at
    }
    
    AUDIT_LOGS {
        bigint id PK
        bigint operator_id FK
        varchar operation_type
        varchar target_type
        bigint target_id
        json before_data
        json after_data
        timestamp created_at
    }
    
    APP_LOGS {
        bigint id PK
        enum level
        text message
        varchar source
        bigint user_id FK
        timestamp created_at
    }
```

---

## 技术栈总览

```mermaid
mindmap
  root((AWSomeShop<br/>技术栈))
    后端
      语言: Go 1.21+
      框架: Gin
      ORM: GORM
      数据库: MySQL 8.0
      认证: JWT
      密码: bcrypt
    前端
      语言: JavaScript ES6+
      UI: Ant Design CSS
      构建: Webpack 5
      服务器: Express
      路由: History API
      HTTP: Fetch API
    开发工具
      版本控制: Git
      容器化: Docker Compose
      测试: Go testing + testify
      API文档: Swagger
    部署
      环境: 本地服务器
      前端端口: 8000
      后端端口: 8080
      数据库端口: 3306
```

---

## 架构特点

### 优点
1. **前后端分离**: 职责清晰，便于独立开发和部署
2. **RESTful API**: 标准化接口，易于理解和维护
3. **分层架构**: Handler → Service → Repository，职责明确
4. **事务支持**: 关键操作使用数据库事务保证数据一致性
5. **审计日志**: 完整记录管理员操作，便于追溯
6. **JWT 认证**: 无状态认证，易于扩展

### 限制
1. **无缓存**: MVP 版本未使用缓存，性能可能受限
2. **无并发控制**: 未处理并发场景，可能出现数据不一致
3. **单数据库**: 所有服务共享一个数据库，耦合度较高
4. **本地存储**: 图片存储在本地文件系统，不支持分布式

### 扩展性
1. **添加缓存**: 可引入 Redis 缓存热点数据
2. **数据库分离**: 可为不同模块使用独立数据库
3. **消息队列**: 可引入消息队列处理异步任务
4. **对象存储**: 可将图片迁移到对象存储服务
5. **负载均衡**: 可部署多个后端实例，使用负载均衡

---

**文档结束**
