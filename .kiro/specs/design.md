# AWSomeShop 设计文档

## 文档信息

- **项目名称**: AWSomeShop 内部员工福利电商系统
- **文档版本**: 1.0
- **创建日期**: 2026-01-14
- **文档状态**: 已批准
- **批准日期**: 2026-01-14
- **基于需求文档**: requirements.md v1.0

## 1. 设计概述

### 1.1 设计目标
本设计文档基于已批准的需求文档，提供 AWSomeShop 系统的详细技术设计方案，包括系统架构、数据库设计、API 接口规范、前端设计和部署方案。

### 1.2 技术栈选型

#### 后端技术栈
- **编程语言**: Go 1.21+
- **Web 框架**: Gin
- **ORM 框架**: GORM
- **数据库**: MySQL 8.0+
- **认证方式**: JWT (JSON Web Token)
- **密码哈希**: bcrypt
- **API 文档**: Swagger/OpenAPI

#### 前端技术栈
- **编程语言**: 原生 JavaScript (ES6+)
- **UI 框架**: Ant Design CSS
- **构建工具**: Webpack 5
- **路由方式**: History API
- **状态管理**: localStorage + Cookie
- **HTTP 客户端**: Fetch API

#### 开发工具
- **版本控制**: Git
- **容器化**: Docker + Docker Compose (仅用于本地开发)
- **测试框架**: Go testing + testify + sqlmock
- **代码覆盖率目标**: 60%


## 2. 系统架构设计

### 2.1 整体架构

系统采用前后端分离架构，包含两个主要服务：

```
┌─────────────────────────────────────────────────────────────┐
│                         用户浏览器                            │
│                    (Chrome/Firefox/Safari/Edge)              │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP/HTTPS
                         │
┌────────────────────────▼────────────────────────────────────┐
│                    前端服务 (Node.js)                        │
│                      端口: 8000                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Express Server                                       │  │
│  │  - 托管静态文件 (Webpack 打包产物)                    │  │
│  │  - 支持 History 路由                                  │  │
│  │  - 原生 JavaScript + Ant Design CSS                  │  │
│  └──────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP (localhost:8080)
                         │ CORS 配置
                         │
┌────────────────────────▼────────────────────────────────────┐
│                    后端服务 (Go + Gin)                       │
│                      端口: 8080                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  API 层                                               │  │
│  │  - RESTful API                                        │  │
│  │  - JWT 认证中间件                                     │  │
│  │  - CORS 中间件                                        │  │
│  │  - 日志中间件                                         │  │
│  │  - 错误处理中间件                                     │  │
│  └──────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  业务逻辑层                                           │  │
│  │  - 用户管理 (员工 + 管理员)                          │  │
│  │  - 产品管理                                           │  │
│  │  - 积分管理                                           │  │
│  │  - 订单管理                                           │  │
│  │  - 审计日志                                           │  │
│  └──────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  数据访问层 (GORM)                                    │  │
│  │  - Model 定义                                         │  │
│  │  - Repository 模式                                    │  │
│  │  - 事务管理                                           │  │
│  └──────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────┘
                         │ MySQL Protocol
                         │
┌────────────────────────▼────────────────────────────────────┐
│                    MySQL 数据库                              │
│                      端口: 3306                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  数据库: awsomeshop                                   │  │
│  │  - users (用户表)                                     │  │
│  │  - products (产品表)                                  │  │
│  │  - orders (订单表)                                    │  │
│  │  - point_logs (积分日志表)                           │  │
│  │  - audit_logs (审计日志表)                           │  │
│  │  - app_logs (应用日志表)                             │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 服务职责划分

#### 前端服务职责
- 托管 Webpack 打包后的静态文件 (HTML/CSS/JS)
- 支持 History 路由（所有路由返回 index.html）
- 页面刷新时重定向到首页
- 不处理业务逻辑

#### 后端服务职责
- 提供 RESTful API 接口
- 处理所有业务逻辑
- 用户认证和授权（JWT）
- 数据库操作
- 文件存储管理（产品图片）
- 日志记录（应用日志 + 审计日志）
- 错误处理和响应


### 2.3 认证授权机制

#### JWT 认证流程

```
员工/管理员登录流程:
┌──────────┐                ┌──────────┐                ┌──────────┐
│  前端    │                │  后端    │                │  数据库  │
└────┬─────┘                └────┬─────┘                └────┬─────┘
     │                           │                           │
     │ POST /api/auth/login      │                           │
     │ {employee_id, password}   │                           │
     ├──────────────────────────>│                           │
     │                           │ 查询用户                   │
     │                           ├──────────────────────────>│
     │                           │                           │
     │                           │ 返回用户信息               │
     │                           │<──────────────────────────┤
     │                           │                           │
     │                           │ 验证密码 (bcrypt)          │
     │                           │                           │
     │                           │ 生成 JWT Token            │
     │                           │                           │
     │ Set-Cookie: token=xxx     │                           │
     │ {user_info}               │                           │
     │<──────────────────────────┤                           │
     │                           │                           │
     │ 存储用户信息到 localStorage│                           │
     │                           │                           │

API 请求认证流程:
┌──────────┐                ┌──────────┐
│  前端    │                │  后端    │
└────┬─────┘                └────┬─────┘
     │                           │
     │ GET /api/products         │
     │ Cookie: token=xxx         │
     ├──────────────────────────>│
     │                           │ 验证 JWT Token
     │                           │ 解析用户信息和权限
     │                           │
     │                           │ 检查权限
     │                           │
     │ 返回数据                   │
     │<──────────────────────────┤
     │                           │
```

#### JWT 配置
- **存储位置**: HTTP Cookie
- **Cookie 属性**:
  - `HttpOnly`: true (防止 XSS 攻击)
  - `Secure`: false (本地开发环境)
  - `SameSite`: Lax (防止 CSRF 攻击)
  - `Path`: /
  - `MaxAge`: 86400 秒 (24 小时)
- **Token 有效期**: 24 小时
- **Token 内容**:
  ```json
  {
    "user_id": "用户ID",
    "employee_id": "工号",
    "role": "employee/admin",
    "exp": "过期时间戳",
    "iat": "签发时间戳"
  }
  ```
- **密钥管理**: 通过环境变量 `JWT_SECRET` 配置

#### 权限控制
- **角色定义**: 
  - `employee`: 员工角色
  - `admin`: 管理员角色
- **权限检查**: 
  - 通过 JWT 中的 `role` 字段判断
  - 中间件在路由层面进行权限拦截
  - 员工只能访问员工端 API
  - 管理员可以访问所有 API


### 2.4 跨域配置

前端服务（localhost:8000）访问后端服务（localhost:8080）需要配置 CORS。

#### 后端 CORS 配置
```go
// Gin CORS 中间件配置
router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:8000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,  // 允许携带 Cookie
    MaxAge:           12 * time.Hour,
}))
```

### 2.5 错误处理机制

#### 统一错误响应格式
```json
{
  "error": "错误描述信息",
  "code": "ERROR_CODE",
  "details": {}
}
```

#### 错误码定义规范
- 命名规则: 大写字母 + 下划线
- 示例错误码:
  - `USER_NOT_FOUND`: 用户不存在
  - `INVALID_PASSWORD`: 密码错误
  - `INSUFFICIENT_POINTS`: 积分不足
  - `PRODUCT_NOT_FOUND`: 产品不存在
  - `PRODUCT_OFFLINE`: 产品已下架
  - `ORDER_NOT_FOUND`: 订单不存在
  - `ORDER_CANNOT_CANCEL`: 订单无法取消
  - `UNAUTHORIZED`: 未授权
  - `FORBIDDEN`: 无权限
  - `INVALID_INPUT`: 输入参数错误
  - `INTERNAL_ERROR`: 内部服务器错误
  - `EMPLOYEE_LIMIT_REACHED`: 当天员工创建数量已达上限

#### HTTP 状态码映射
- `200 OK`: 请求成功
- `201 Created`: 资源创建成功
- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 未认证
- `403 Forbidden`: 无权限
- `404 Not Found`: 资源不存在
- `409 Conflict`: 资源冲突（如工号重复）
- `500 Internal Server Error`: 服务器内部错误

#### 前端错误提示
- 使用 Toast 提示显示错误信息
- 显示 `error` 字段的内容
- 自动消失时间: 3 秒


## 3. 数据库设计

### 3.1 数据库概述
- **数据库名称**: awsomeshop
- **字符集**: utf8mb4
- **排序规则**: utf8mb4_unicode_ci
- **存储引擎**: InnoDB
- **事务支持**: 是（用于兑换操作）

### 3.2 表结构设计

#### 3.2.1 用户表 (users)

存储员工和管理员账号信息（统一表）。

```sql
CREATE TABLE `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `employee_id` VARCHAR(50) NOT NULL COMMENT '工号',
  `name` VARCHAR(100) NOT NULL COMMENT '姓名',
  `department` VARCHAR(100) NOT NULL COMMENT '部门',
  `position` VARCHAR(100) NOT NULL COMMENT '岗位',
  `email` VARCHAR(100) NOT NULL COMMENT '邮箱',
  `phone` VARCHAR(20) NOT NULL COMMENT '手机号',
  `password` VARCHAR(255) NOT NULL COMMENT '密码哈希',
  `role` ENUM('employee', 'admin') NOT NULL DEFAULT 'employee' COMMENT '角色',
  `status` ENUM('active', 'resigned', 'probation') NOT NULL DEFAULT 'active' COMMENT '状态',
  `points` INT NOT NULL DEFAULT 0 COMMENT '积分余额',
  `hire_date` DATE NOT NULL COMMENT '入职时间',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_employee_id` (`employee_id`),
  KEY `idx_role` (`role`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表（员工+管理员）';
```

**字段说明**:
- `employee_id`: 工号，格式为 `YYYYMMDD-NNN` (如 20260114-001)，初始管理员为 `admin`
- `role`: 角色，`employee` 为员工，`admin` 为管理员
- `status`: 状态，`active` 在职，`resigned` 离职，`probation` 试用期
- `points`: 积分余额，新员工创建时自动设置为 1000

#### 3.2.2 产品表 (products)

```sql
CREATE TABLE `products` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '产品ID',
  `name` VARCHAR(100) NOT NULL COMMENT '产品名称',
  `description` TEXT NOT NULL COMMENT '产品描述',
  `image_url` VARCHAR(255) DEFAULT NULL COMMENT '产品图片URL',
  `points_required` INT NOT NULL COMMENT '所需积分',
  `status` ENUM('online', 'offline') NOT NULL DEFAULT 'offline' COMMENT '上架状态',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除（软删除）',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_is_deleted` (`is_deleted`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='产品表';
```

**字段说明**:
- `image_url`: 图片访问路径，格式为 `/images/{uuid}.{ext}`
- `status`: `online` 已上架，`offline` 已下架
- `is_deleted`: 软删除标记，0 未删除，1 已删除


#### 3.2.3 订单表 (orders)

```sql
CREATE TABLE `orders` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `order_no` VARCHAR(50) NOT NULL COMMENT '订单号',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `product_id` BIGINT UNSIGNED NOT NULL COMMENT '产品ID',
  `product_snapshot` JSON NOT NULL COMMENT '产品快照',
  `points_cost` INT NOT NULL COMMENT '消耗积分',
  `status` ENUM('pending', 'approved', 'rejected', 'cancelled') NOT NULL DEFAULT 'pending' COMMENT '订单状态',
  `applied_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
  `reviewed_at` TIMESTAMP NULL DEFAULT NULL COMMENT '审核时间',
  `reviewer_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '审核人ID',
  `reject_reason` TEXT NULL DEFAULT NULL COMMENT '拒绝原因',
  `approval_note` TEXT NULL DEFAULT NULL COMMENT '核销备注',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_status` (`status`),
  KEY `idx_applied_at` (`applied_at`),
  CONSTRAINT `fk_orders_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_orders_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT `fk_orders_reviewer` FOREIGN KEY (`reviewer_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';
```

**字段说明**:
- `order_no`: 订单号，格式为 `YYYYMMDDHHmmss-XXXXXX` (如 20260114153045-ABC123)
- `product_snapshot`: 产品快照，JSON 格式存储兑换时的产品信息
  ```json
  {
    "name": "产品名称",
    "description": "产品描述",
    "image_url": "/images/xxx.jpg",
    "points_required": 500
  }
  ```
- `status`: 订单状态
  - `pending`: 待核销
  - `approved`: 已核销
  - `rejected`: 已拒绝
  - `cancelled`: 已取消

#### 3.2.4 积分日志表 (point_logs)

```sql
CREATE TABLE `point_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `operator_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '操作人ID',
  `amount` INT NOT NULL COMMENT '积分变动数量（正数增加，负数减少）',
  `reason` TEXT NOT NULL COMMENT '变动原因',
  `type` ENUM('onboarding', 'admin_grant', 'admin_deduct', 'redeem', 'order_cancel_refund', 'order_reject_refund') NOT NULL COMMENT '变动类型',
  `order_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '关联订单ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_operator_id` (`operator_id`),
  KEY `idx_type` (`type`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_point_logs_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_point_logs_operator` FOREIGN KEY (`operator_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_point_logs_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='积分日志表';
```

**字段说明**:
- `type`: 变动类型
  - `onboarding`: 入职发放
  - `admin_grant`: 管理员发放
  - `admin_deduct`: 管理员扣除
  - `redeem`: 兑换消费
  - `order_cancel_refund`: 订单取消退回
  - `order_reject_refund`: 订单拒绝退回


#### 3.2.5 审计日志表 (audit_logs)

```sql
CREATE TABLE `audit_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `operator_id` BIGINT UNSIGNED NOT NULL COMMENT '操作人ID',
  `operation_type` VARCHAR(50) NOT NULL COMMENT '操作类型',
  `target_type` VARCHAR(50) NOT NULL COMMENT '操作对象类型',
  `target_id` BIGINT UNSIGNED NOT NULL COMMENT '操作对象ID',
  `before_data` JSON NULL DEFAULT NULL COMMENT '操作前数据',
  `after_data` JSON NULL DEFAULT NULL COMMENT '操作后数据',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `idx_operator_id` (`operator_id`),
  KEY `idx_operation_type` (`operation_type`),
  KEY `idx_target` (`target_type`, `target_id`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_audit_logs_operator` FOREIGN KEY (`operator_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='审计日志表';
```

**字段说明**:
- `operation_type`: 操作类型
  - `create_employee`: 创建员工
  - `update_employee`: 修改员工
  - `create_product`: 创建产品
  - `update_product`: 修改产品
  - `delete_product`: 删除产品
  - `grant_points`: 发放积分
  - `deduct_points`: 扣除积分
  - `approve_order`: 核销订单
  - `reject_order`: 拒绝订单
- `target_type`: 操作对象类型 (`user`, `product`, `order`, `point`)
- `before_data` / `after_data`: JSON 格式存储操作前后的完整数据

#### 3.2.6 应用日志表 (app_logs)

```sql
CREATE TABLE `app_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `level` ENUM('INFO', 'WARNING', 'ERROR') NOT NULL DEFAULT 'INFO' COMMENT '日志级别',
  `message` TEXT NOT NULL COMMENT '日志消息',
  `source` VARCHAR(100) NOT NULL COMMENT '日志来源',
  `user_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '关联用户ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_level` (`level`),
  KEY `idx_source` (`source`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用日志表';
```

**字段说明**:
- `level`: 日志级别 (INFO, WARNING, ERROR)
- `source`: 日志来源（如 `auth`, `product`, `order`, `point`）
- `user_id`: 关联用户ID（可选）

### 3.3 数据库索引策略

#### 主键索引
- 所有表使用自增 BIGINT 作为主键
- 保证查询性能和数据唯一性

#### 唯一索引
- `users.employee_id`: 工号唯一
- `orders.order_no`: 订单号唯一

#### 普通索引
- 外键字段建立索引（提升关联查询性能）
- 状态字段建立索引（提升筛选查询性能）
- 时间字段建立索引（提升排序和范围查询性能）

#### 复合索引
- `audit_logs(target_type, target_id)`: 提升按对象查询审计日志的性能


### 3.4 数据库 ER 图

```
┌─────────────────────────────────────────────────────────────────┐
│                            users                                 │
│─────────────────────────────────────────────────────────────────│
│ PK  id                BIGINT                                     │
│ UK  employee_id       VARCHAR(50)                                │
│     name              VARCHAR(100)                               │
│     department        VARCHAR(100)                               │
│     position          VARCHAR(100)                               │
│     email             VARCHAR(100)                               │
│     phone             VARCHAR(20)                                │
│     password          VARCHAR(255)                               │
│     role              ENUM('employee','admin')                   │
│     status            ENUM('active','resigned','probation')      │
│     points            INT                                        │
│     hire_date         DATE                                       │
│     created_at        TIMESTAMP                                  │
│     updated_at        TIMESTAMP                                  │
└──────────────┬──────────────────────────────────────────────────┘
               │
               │ 1:N (user_id)
               │
┌──────────────▼──────────────────────────────────────────────────┐
│                           orders                                 │
│─────────────────────────────────────────────────────────────────│
│ PK  id                BIGINT                                     │
│ UK  order_no          VARCHAR(50)                                │
│ FK  user_id           BIGINT                                     │
│ FK  product_id        BIGINT                                     │
│     product_snapshot  JSON                                       │
│     points_cost       INT                                        │
│     status            ENUM('pending','approved',...)             │
│     applied_at        TIMESTAMP                                  │
│     reviewed_at       TIMESTAMP                                  │
│ FK  reviewer_id       BIGINT                                     │
│     reject_reason     TEXT                                       │
│     approval_note     TEXT                                       │
│     created_at        TIMESTAMP                                  │
│     updated_at        TIMESTAMP                                  │
└──────────────┬──────────────────────────────────────────────────┘
               │
               │ N:1 (product_id)
               │
┌──────────────▼──────────────────────────────────────────────────┐
│                          products                                │
│─────────────────────────────────────────────────────────────────│
│ PK  id                BIGINT                                     │
│     name              VARCHAR(100)                               │
│     description       TEXT                                       │
│     image_url         VARCHAR(255)                               │
│     points_required   INT                                        │
│     status            ENUM('online','offline')                   │
│     is_deleted        TINYINT(1)                                 │
│     created_at        TIMESTAMP                                  │
│     updated_at        TIMESTAMP                                  │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                         point_logs                               │
│─────────────────────────────────────────────────────────────────│
│ PK  id                BIGINT                                     │
│ FK  user_id           BIGINT                                     │
│ FK  operator_id       BIGINT                                     │
│     amount            INT                                        │
│     reason            TEXT                                       │
│     type              ENUM('onboarding','admin_grant',...)       │
│ FK  order_id          BIGINT                                     │
│     created_at        TIMESTAMP                                  │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                         audit_logs                               │
│─────────────────────────────────────────────────────────────────│
│ PK  id                BIGINT                                     │
│ FK  operator_id       BIGINT                                     │
│     operation_type    VARCHAR(50)                                │
│     target_type       VARCHAR(50)                                │
│     target_id         BIGINT                                     │
│     before_data       JSON                                       │
│     after_data        JSON                                       │
│     created_at        TIMESTAMP                                  │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                          app_logs                                │
│─────────────────────────────────────────────────────────────────│
│ PK  id                BIGINT                                     │
│     level             ENUM('INFO','WARNING','ERROR')             │
│     message           TEXT                                       │
│     source            VARCHAR(100)                               │
│ FK  user_id           BIGINT                                     │
│     created_at        TIMESTAMP                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 3.5 事务管理

#### 兑换产品事务
兑换产品时需要保证数据一致性，使用数据库事务：

```go
// 伪代码示例
tx := db.Begin()

// 1. 扣除用户积分
err := tx.Model(&user).Update("points", gorm.Expr("points - ?", pointsCost)).Error
if err != nil {
    tx.Rollback()
    return err
}

// 2. 创建订单
err = tx.Create(&order).Error
if err != nil {
    tx.Rollback()
    return err
}

// 3. 记录积分日志
err = tx.Create(&pointLog).Error
if err != nil {
    tx.Rollback()
    return err
}

tx.Commit()
```

#### 订单拒绝/取消事务
订单拒绝或取消时需要退回积分：

```go
tx := db.Begin()

// 1. 更新订单状态
err := tx.Model(&order).Updates(map[string]interface{}{
    "status": newStatus,
    "reviewed_at": time.Now(),
}).Error
if err != nil {
    tx.Rollback()
    return err
}

// 2. 退回积分
err = tx.Model(&user).Update("points", gorm.Expr("points + ?", pointsCost)).Error
if err != nil {
    tx.Rollback()
    return err
}

// 3. 记录积分日志
err = tx.Create(&pointLog).Error
if err != nil {
    tx.Rollback()
    return err
}

tx.Commit()
```


## 4. API 接口设计

### 4.1 API 设计原则

- **RESTful 风格**: 使用标准 HTTP 方法（GET, POST, PUT, DELETE）
- **统一前缀**: 所有 API 以 `/api` 开头
- **版本控制**: 暂不使用版本号（MVP 版本）
- **响应格式**: JSON
- **分页参数**: `page` (页码，从1开始), `page_size` (每页数量，默认50)
- **认证方式**: JWT Token (存储在 Cookie 中)

### 4.2 通用响应格式

#### 成功响应
```json
{
  "data": {},
  "message": "操作成功"
}
```

#### 分页响应
```json
{
  "data": {
    "items": [],
    "total": 100,
    "page": 1,
    "page_size": 50,
    "total_pages": 2
  },
  "message": "查询成功"
}
```

#### 错误响应
```json
{
  "error": "错误描述",
  "code": "ERROR_CODE",
  "details": {}
}
```

### 4.3 认证相关 API

#### 4.3.1 用户登录
```
POST /api/auth/login
```

**请求体**:
```json
{
  "employee_id": "admin",
  "password": "123456"
}
```

**成功响应** (200):
```json
{
  "data": {
    "user": {
      "id": 1,
      "employee_id": "admin",
      "name": "管理员",
      "role": "admin",
      "points": 0
    }
  },
  "message": "登录成功"
}
```
**Set-Cookie**: `token=eyJhbGc...; HttpOnly; SameSite=Lax; Max-Age=86400; Path=/`

**错误响应**:
- 401: `USER_NOT_FOUND` 或 `INVALID_PASSWORD`
- 403: `ACCOUNT_DISABLED` (账号已禁用)

#### 4.3.2 用户登出
```
POST /api/auth/logout
```

**成功响应** (200):
```json
{
  "message": "登出成功"
}
```
**Set-Cookie**: `token=; Max-Age=0; Path=/`

#### 4.3.3 修改密码
```
PUT /api/auth/password
```

**请求体**:
```json
{
  "old_password": "123456",
  "new_password": "newpass"
}
```

**成功响应** (200):
```json
{
  "message": "密码修改成功"
}
```

**错误响应**:
- 401: `INVALID_PASSWORD` (旧密码错误)


### 4.4 员工管理 API (管理员)

#### 4.4.1 创建员工
```
POST /api/admin/employees
```

**权限**: 管理员

**请求体**:
```json
{
  "name": "张三",
  "department": "技术部",
  "position": "工程师",
  "email": "zhangsan@example.com",
  "phone": "13800138000",
  "hire_date": "2026-01-14",
  "password": "123456",
  "status": "active"
}
```

**成功响应** (201):
```json
{
  "data": {
    "id": 2,
    "employee_id": "20260114-001",
    "name": "张三",
    "role": "employee",
    "points": 1000
  },
  "message": "员工创建成功"
}
```

**错误响应**:
- 409: `EMPLOYEE_LIMIT_REACHED` (当天创建数量已达999)

#### 4.4.2 获取员工列表
```
GET /api/admin/employees?page=1&page_size=50&name=&department=&position=
```

**权限**: 管理员

**查询参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认50）
- `name`: 姓名筛选（可选）
- `department`: 部门筛选（可选）
- `position`: 岗位筛选（可选）

**成功响应** (200):
```json
{
  "data": {
    "items": [
      {
        "id": 2,
        "employee_id": "20260114-001",
        "name": "张三",
        "department": "技术部",
        "position": "工程师",
        "email": "zhangsan@example.com",
        "phone": "13800138000",
        "status": "active",
        "points": 1000,
        "hire_date": "2026-01-14",
        "created_at": "2026-01-14T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 50,
    "total_pages": 2
  },
  "message": "查询成功"
}
```

#### 4.4.3 获取员工详情
```
GET /api/admin/employees/:id
```

**权限**: 管理员

**成功响应** (200):
```json
{
  "data": {
    "id": 2,
    "employee_id": "20260114-001",
    "name": "张三",
    "department": "技术部",
    "position": "工程师",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "role": "employee",
    "status": "active",
    "points": 1000,
    "hire_date": "2026-01-14",
    "created_at": "2026-01-14T10:00:00Z",
    "updated_at": "2026-01-14T10:00:00Z"
  },
  "message": "查询成功"
}
```

#### 4.4.4 更新员工信息
```
PUT /api/admin/employees/:id
```

**权限**: 管理员

**请求体**:
```json
{
  "name": "张三",
  "department": "技术部",
  "position": "高级工程师",
  "email": "zhangsan@example.com",
  "phone": "13800138000",
  "status": "active",
  "hire_date": "2026-01-14"
}
```

**成功响应** (200):
```json
{
  "data": {
    "id": 2,
    "employee_id": "20260114-001",
    "name": "张三",
    "position": "高级工程师"
  },
  "message": "更新成功"
}
```

#### 4.4.5 重置员工密码
```
PUT /api/admin/employees/:id/password
```

**权限**: 管理员

**请求体**:
```json
{
  "new_password": "newpass123"
}
```

**成功响应** (200):
```json
{
  "message": "密码重置成功"
}
```


### 4.5 管理员管理 API (管理员)

#### 4.5.1 创建管理员
```
POST /api/admin/admins
```

**权限**: 管理员（仅 admin 工号可创建）

**请求体**:
```json
{
  "name": "李四",
  "department": "运营部",
  "position": "运营经理",
  "email": "lisi@example.com",
  "phone": "13900139000",
  "hire_date": "2026-01-14",
  "password": "123456",
  "status": "active"
}
```

**成功响应** (201):
```json
{
  "data": {
    "id": 3,
    "employee_id": "20260114-002",
    "name": "李四",
    "role": "admin"
  },
  "message": "管理员创建成功"
}
```

### 4.6 产品管理 API

#### 4.6.1 创建产品 (管理员)
```
POST /api/admin/products
```

**权限**: 管理员

**请求体** (multipart/form-data):
```
name: "AirPods Pro"
description: "苹果无线耳机"
points_required: 500
image: <file>
```

**成功响应** (201):
```json
{
  "data": {
    "id": 1,
    "name": "AirPods Pro",
    "description": "苹果无线耳机",
    "image_url": "/images/550e8400-e29b-41d4-a716-446655440000.jpg",
    "points_required": 500,
    "status": "offline"
  },
  "message": "产品创建成功"
}
```

**错误响应**:
- 400: `INVALID_INPUT` (参数错误)
- 400: `FILE_TOO_LARGE` (图片超过1MB)

#### 4.6.2 获取产品列表 (管理员)
```
GET /api/admin/products?page=1&page_size=50&include_deleted=false
```

**权限**: 管理员

**查询参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认50）
- `include_deleted`: 是否包含已删除产品（默认false）

**成功响应** (200):
```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "name": "AirPods Pro",
        "description": "苹果无线耳机",
        "image_url": "/images/xxx.jpg",
        "points_required": 500,
        "status": "online",
        "is_deleted": false,
        "created_at": "2026-01-14T10:00:00Z"
      }
    ],
    "total": 9,
    "page": 1,
    "page_size": 50,
    "total_pages": 1
  },
  "message": "查询成功"
}
```

#### 4.6.3 获取产品列表 (员工)
```
GET /api/products?page=1&page_size=50
```

**权限**: 员工

**成功响应** (200):
```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "name": "AirPods Pro",
        "description": "苹果无线耳机",
        "image_url": "/images/xxx.jpg",
        "points_required": 500,
        "status": "online",
        "created_at": "2026-01-14T10:00:00Z"
      }
    ],
    "total": 9,
    "page": 1,
    "page_size": 50,
    "total_pages": 1
  },
  "message": "查询成功"
}
```

#### 4.6.4 获取产品详情
```
GET /api/products/:id
```

**权限**: 员工/管理员

**成功响应** (200):
```json
{
  "data": {
    "id": 1,
    "name": "AirPods Pro",
    "description": "苹果无线耳机",
    "image_url": "/images/xxx.jpg",
    "points_required": 500,
    "status": "online",
    "created_at": "2026-01-14T10:00:00Z",
    "updated_at": "2026-01-14T10:00:00Z"
  },
  "message": "查询成功"
}
```

#### 4.6.5 更新产品 (管理员)
```
PUT /api/admin/products/:id
```

**权限**: 管理员

**请求体** (multipart/form-data):
```
name: "AirPods Pro 2"
description: "苹果无线耳机第二代"
points_required: 600
image: <file> (可选)
```

**成功响应** (200):
```json
{
  "data": {
    "id": 1,
    "name": "AirPods Pro 2",
    "points_required": 600
  },
  "message": "更新成功"
}
```

#### 4.6.6 删除产品 (管理员)
```
DELETE /api/admin/products/:id
```

**权限**: 管理员

**成功响应** (200):
```json
{
  "message": "删除成功"
}
```

#### 4.6.7 上架/下架产品 (管理员)
```
PUT /api/admin/products/:id/status
```

**权限**: 管理员

**请求体**:
```json
{
  "status": "online"
}
```

**成功响应** (200):
```json
{
  "message": "状态更新成功"
}
```


### 4.7 图片访问 API

#### 4.7.1 获取图片
```
GET /images/:filename
```

**权限**: 公开（需登录）

**示例**: `GET /images/550e8400-e29b-41d4-a716-446655440000.jpg`

**成功响应**: 返回图片文件（Content-Type: image/jpeg 或 image/png 等）

**错误响应**:
- 404: `FILE_NOT_FOUND`

### 4.8 积分管理 API

#### 4.8.1 获取个人积分余额 (员工)
```
GET /api/points/balance
```

**权限**: 员工

**成功响应** (200):
```json
{
  "data": {
    "points": 1500,
    "user_id": 2,
    "employee_id": "20260114-001"
  },
  "message": "查询成功"
}
```

#### 4.8.2 获取个人积分日志 (员工)
```
GET /api/points/logs?page=1&page_size=50
```

**权限**: 员工

**成功响应** (200):
```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "amount": 1000,
        "reason": "入职发放",
        "type": "onboarding",
        "created_at": "2026-01-14T10:00:00Z"
      },
      {
        "id": 2,
        "amount": -500,
        "reason": "兑换产品：AirPods Pro",
        "type": "redeem",
        "order_id": 1,
        "created_at": "2026-01-14T11:00:00Z"
      }
    ],
    "total": 2,
    "page": 1,
    "page_size": 50,
    "total_pages": 1
  },
  "message": "查询成功"
}
```

#### 4.8.3 获取所有积分日志 (管理员)
```
GET /api/admin/points/logs?page=1&page_size=50
```

**权限**: 管理员

**成功响应** (200):
```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "user_id": 2,
        "user_name": "张三",
        "employee_id": "20260114-001",
        "operator_id": null,
        "operator_name": null,
        "amount": 1000,
        "reason": "入职发放",
        "type": "onboarding",
        "order_id": null,
        "created_at": "2026-01-14T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 50,
    "total_pages": 2
  },
  "message": "查询成功"
}
```

#### 4.8.4 发放积分 (管理员)
```
POST /api/admin/points/grant
```

**权限**: 管理员

**请求体**:
```json
{
  "user_id": 2,
  "amount": 500,
  "reason": "月度绩效奖励"
}
```

**成功响应** (200):
```json
{
  "data": {
    "log_id": 3,
    "user_id": 2,
    "new_balance": 1500
  },
  "message": "积分发放成功"
}
```

#### 4.8.5 批量发放积分 (管理员)
```
POST /api/admin/points/grant-batch
```

**权限**: 管理员

**请求体**:
```json
{
  "user_ids": [2, 3, 4],
  "amount": 500,
  "reason": "月度绩效奖励"
}
```

或按条件批量发放:
```json
{
  "filter": {
    "department": "技术部"
  },
  "amount": 500,
  "reason": "部门奖励"
}
```

**成功响应** (200):
```json
{
  "data": {
    "success_count": 3,
    "failed_count": 0
  },
  "message": "批量发放成功"
}
```

#### 4.8.6 扣除积分 (管理员)
```
POST /api/admin/points/deduct
```

**权限**: 管理员

**请求体**:
```json
{
  "user_id": 2,
  "amount": 100,
  "reason": "违规扣除"
}
```

**成功响应** (200):
```json
{
  "data": {
    "log_id": 4,
    "user_id": 2,
    "new_balance": 1400
  },
  "message": "积分扣除成功"
}
```


### 4.9 订单管理 API

#### 4.9.1 兑换产品 (员工)
```
POST /api/orders/redeem
```

**权限**: 员工

**请求体**:
```json
{
  "product_id": 1
}
```

**成功响应** (201):
```json
{
  "data": {
    "order_id": 1,
    "order_no": "20260114153045-ABC123",
    "product_name": "AirPods Pro",
    "points_cost": 500,
    "status": "pending",
    "new_balance": 500
  },
  "message": "兑换成功"
}
```

**错误响应**:
- 400: `PRODUCT_NOT_FOUND`
- 400: `PRODUCT_OFFLINE` (产品已下架)
- 400: `INSUFFICIENT_POINTS` (积分不足)

#### 4.9.2 获取个人订单列表 (员工)
```
GET /api/orders?page=1&page_size=50
```

**权限**: 员工

**成功响应** (200):
```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "order_no": "20260114153045-ABC123",
        "product": {
          "name": "AirPods Pro",
          "image_url": "/images/xxx.jpg"
        },
        "points_cost": 500,
        "status": "pending",
        "applied_at": "2026-01-14T15:30:45Z",
        "reviewed_at": null,
        "reviewer_name": null,
        "reject_reason": null,
        "approval_note": null
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 50,
    "total_pages": 1
  },
  "message": "查询成功"
}
```

#### 4.9.3 获取订单详情 (员工)
```
GET /api/orders/:id
```

**权限**: 员工（只能查看自己的订单）

**成功响应** (200):
```json
{
  "data": {
    "id": 1,
    "order_no": "20260114153045-ABC123",
    "product_snapshot": {
      "name": "AirPods Pro",
      "description": "苹果无线耳机",
      "image_url": "/images/xxx.jpg",
      "points_required": 500
    },
    "points_cost": 500,
    "status": "pending",
    "applied_at": "2026-01-14T15:30:45Z",
    "reviewed_at": null,
    "reviewer_name": null,
    "reject_reason": null,
    "approval_note": null
  },
  "message": "查询成功"
}
```

#### 4.9.4 取消订单 (员工)
```
PUT /api/orders/:id/cancel
```

**权限**: 员工（只能取消自己的订单）

**成功响应** (200):
```json
{
  "data": {
    "order_id": 1,
    "status": "cancelled",
    "refunded_points": 500,
    "new_balance": 1000
  },
  "message": "订单已取消，积分已退回"
}
```

**错误响应**:
- 400: `ORDER_CANNOT_CANCEL` (订单状态不是 pending)

#### 4.9.5 获取所有订单列表 (管理员)
```
GET /api/admin/orders?page=1&page_size=50&status=
```

**权限**: 管理员

**查询参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认50）
- `status`: 状态筛选（可选：pending, approved, rejected, cancelled）

**成功响应** (200):
```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "order_no": "20260114153045-ABC123",
        "user": {
          "id": 2,
          "employee_id": "20260114-001",
          "name": "张三"
        },
        "product": {
          "name": "AirPods Pro",
          "image_url": "/images/xxx.jpg"
        },
        "points_cost": 500,
        "status": "pending",
        "applied_at": "2026-01-14T15:30:45Z",
        "reviewed_at": null,
        "reviewer_name": null
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 50,
    "total_pages": 1
  },
  "message": "查询成功"
}
```

#### 4.9.6 核销订单 (管理员)
```
PUT /api/admin/orders/:id/approve
```

**权限**: 管理员

**请求体**:
```json
{
  "approval_note": "已发货"
}
```

**成功响应** (200):
```json
{
  "data": {
    "order_id": 1,
    "status": "approved",
    "reviewed_at": "2026-01-14T16:00:00Z"
  },
  "message": "订单核销成功"
}
```

#### 4.9.7 拒绝订单 (管理员)
```
PUT /api/admin/orders/:id/reject
```

**权限**: 管理员

**请求体**:
```json
{
  "reject_reason": "产品缺货"
}
```

**成功响应** (200):
```json
{
  "data": {
    "order_id": 1,
    "status": "rejected",
    "refunded_points": 500,
    "reviewed_at": "2026-01-14T16:00:00Z"
  },
  "message": "订单已拒绝，积分已退回"
}
```

**错误响应**:
- 400: `REJECT_REASON_REQUIRED` (拒绝原因必填)


### 4.10 系统概况 API (管理员)

#### 4.10.1 获取系统概况
```
GET /api/admin/dashboard
```

**权限**: 管理员

**成功响应** (200):
```json
{
  "data": {
    "total_employees": 100,
    "total_products": 9,
    "pending_orders": 5,
    "today_orders": 12,
    "month_orders": 156,
    "today_points_granted": 5000,
    "month_points_granted": 50000
  },
  "message": "查询成功"
}
```

### 4.11 个人中心 API (员工)

#### 4.11.1 获取个人信息
```
GET /api/profile
```

**权限**: 员工

**成功响应** (200):
```json
{
  "data": {
    "id": 2,
    "employee_id": "20260114-001",
    "name": "张三",
    "department": "技术部",
    "position": "工程师",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "status": "active",
    "points": 1000,
    "hire_date": "2026-01-14",
    "created_at": "2026-01-14T10:00:00Z"
  },
  "message": "查询成功"
}
```

## 5. 前端设计

### 5.1 技术架构

```
┌─────────────────────────────────────────────────────────────┐
│                      前端应用架构                            │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌────────────────────────────────────────────────────┐    │
│  │  Express Server (端口 8000)                        │    │
│  │  - 托管静态文件                                     │    │
│  │  - History 路由支持                                │    │
│  └────────────────────────────────────────────────────┘    │
│                                                              │
│  ┌────────────────────────────────────────────────────┐    │
│  │  静态资源 (dist/)                                   │    │
│  │  ├── index.html                                     │    │
│  │  ├── js/                                            │    │
│  │  │   ├── app.js (主应用逻辑)                       │    │
│  │  │   ├── router.js (路由管理)                      │    │
│  │  │   ├── api.js (API 调用)                         │    │
│  │  │   ├── auth.js (认证逻辑)                        │    │
│  │  │   └── utils.js (工具函数)                       │    │
│  │  ├── css/                                            │    │
│  │  │   ├── antd.min.css (Ant Design 样式)            │    │
│  │  │   └── custom.css (自定义样式)                   │    │
│  │  └── images/                                         │    │
│  └────────────────────────────────────────────────────┘    │
│                                                              │
│  ┌────────────────────────────────────────────────────┐    │
│  │  状态管理                                           │    │
│  │  ├── Cookie: JWT Token                             │    │
│  │  └── localStorage: 用户信息                        │    │
│  └────────────────────────────────────────────────────┘    │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 5.2 页面结构

#### 5.2.1 员工端页面

```
员工端路由结构:
/
├── /login                    # 登录页
├── /products                 # 产品列表页（首页）
├── /products/:id             # 产品详情页
├── /profile                  # 个人中心
│   ├── /profile/info         # 个人信息
│   ├── /profile/points       # 积分余额和日志
│   ├── /profile/orders       # 兑换历史
│   └── /profile/password     # 修改密码
└── /orders/:id               # 订单详情页
```

#### 5.2.2 管理员端页面

```
管理员端路由结构:
/
├── /login                    # 登录页
├── /admin/dashboard          # 系统概况（首页）
├── /admin/employees          # 员工管理
│   ├── /admin/employees/list # 员工列表
│   ├── /admin/employees/new  # 创建员工
│   └── /admin/employees/:id  # 员工详情/编辑
├── /admin/products           # 产品管理
│   ├── /admin/products/list  # 产品列表
│   ├── /admin/products/new   # 创建产品
│   └── /admin/products/:id   # 产品详情/编辑
├── /admin/points             # 积分管理
│   ├── /admin/points/grant   # 发放积分
│   └── /admin/points/logs    # 积分日志
├── /admin/orders             # 订单管理
│   ├── /admin/orders/list    # 订单列表
│   └── /admin/orders/:id     # 订单详情
└── /admin/profile            # 个人设置
    └── /admin/profile/password # 修改密码
```


### 5.3 页面布局设计

#### 5.3.1 通用布局

```
┌─────────────────────────────────────────────────────────────┐
│  顶部导航栏                                                  │
│  Logo | 导航菜单 | 用户信息 | 登出                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│                                                              │
│                     主内容区域                               │
│                                                              │
│                                                              │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

#### 5.3.2 员工端导航菜单
- 产品列表
- 个人中心
  - 个人信息
  - 积分余额
  - 兑换历史
  - 修改密码

#### 5.3.3 管理员端导航菜单
- 系统概况
- 员工管理
- 产品管理
- 积分管理
- 订单管理
- 个人设置

### 5.4 关键页面设计

#### 5.4.1 登录页
- 工号输入框
- 密码输入框
- 登录按钮
- 错误提示区域

#### 5.4.2 产品列表页（员工）
- 产品卡片网格布局（每行3-4个）
- 每个卡片显示：
  - 产品图片
  - 产品名称
  - 所需积分
  - 上架状态标签
  - 查看详情按钮

#### 5.4.3 产品详情页（员工）
- 产品图片（大图）
- 产品名称
- 产品描述
- 所需积分
- 上架状态
- 当前积分余额
- 兑换按钮（根据状态禁用）

#### 5.4.4 系统概况页（管理员）
- 统计卡片布局
  - 员工总数
  - 产品总数
  - 待审核订单数
  - 今日兑换订单数
  - 本月兑换订单数
  - 今日积分发放总量
  - 本月积分发放总量

#### 5.4.5 员工列表页（管理员）
- 筛选区域（姓名、部门、岗位）
- 创建员工按钮
- 表格显示：
  - 工号
  - 姓名
  - 部门
  - 岗位
  - 状态
  - 积分余额
  - 操作（查看/编辑/重置密码）
- 分页组件

#### 5.4.6 产品列表页（管理员）
- 创建产品按钮
- 表格显示：
  - 产品图片（缩略图）
  - 产品名称
  - 所需积分
  - 上架状态
  - 创建时间
  - 操作（编辑/删除/上架/下架）
- 分页组件

#### 5.4.7 订单列表页（管理员）
- 状态筛选（待核销/已核销/已拒绝/已取消）
- 表格显示：
  - 订单号
  - 员工信息
  - 产品信息
  - 消耗积分
  - 订单状态
  - 申请时间
  - 操作（核销/拒绝/查看详情）
- 分页组件

### 5.5 前端路由实现

使用 History API 实现前端路由：

```javascript
// router.js
class Router {
  constructor() {
    this.routes = {};
    this.init();
  }

  init() {
    window.addEventListener('popstate', () => {
      this.loadRoute(window.location.pathname);
    });
  }

  register(path, handler) {
    this.routes[path] = handler;
  }

  navigate(path) {
    window.history.pushState({}, '', path);
    this.loadRoute(path);
  }

  loadRoute(path) {
    const handler = this.routes[path] || this.routes['/404'];
    if (handler) {
      handler();
    }
  }
}
```

### 5.6 状态管理

#### localStorage 存储内容
```javascript
{
  "user": {
    "id": 2,
    "employee_id": "20260114-001",
    "name": "张三",
    "role": "employee",
    "points": 1000
  }
}
```

#### Cookie 存储内容
- `token`: JWT Token (HttpOnly, 由后端设置)

### 5.7 API 调用封装

```javascript
// api.js
class API {
  constructor(baseURL) {
    this.baseURL = baseURL;
  }

  async request(method, path, data = null) {
    const options = {
      method,
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include', // 携带 Cookie
    };

    if (data) {
      options.body = JSON.stringify(data);
    }

    const response = await fetch(`${this.baseURL}${path}`, options);
    const result = await response.json();

    if (!response.ok) {
      throw new Error(result.error || '请求失败');
    }

    return result;
  }

  get(path) {
    return this.request('GET', path);
  }

  post(path, data) {
    return this.request('POST', path, data);
  }

  put(path, data) {
    return this.request('PUT', path, data);
  }

  delete(path) {
    return this.request('DELETE', path);
  }
}

const api = new API('http://localhost:8080/api');
```

### 5.8 错误处理

使用 Toast 提示显示错误信息：

```javascript
// utils.js
function showToast(message, type = 'info') {
  // 使用 Ant Design 的 message 组件样式
  const toast = document.createElement('div');
  toast.className = `ant-message-${type}`;
  toast.textContent = message;
  document.body.appendChild(toast);
  
  setTimeout(() => {
    toast.remove();
  }, 3000);
}

// 使用示例
try {
  await api.post('/orders/redeem', { product_id: 1 });
  showToast('兑换成功', 'success');
} catch (error) {
  showToast(error.message, 'error');
}
```


## 6. 文件存储设计

### 6.1 图片存储

#### 存储位置
- 本地文件系统：`/uploads/image/`
- 相对于后端服务根目录

#### 文件命名规则
- 使用 UUID v4 生成唯一文件名
- 保留原始文件扩展名
- 示例：`550e8400-e29b-41d4-a716-446655440000.jpg`

#### 上传流程
```
1. 前端选择图片文件
2. 通过 multipart/form-data 上传到后端
3. 后端验证文件大小（≤1MB）
4. 生成 UUID 文件名
5. 保存到 /uploads/image/ 目录
6. 返回图片访问路径：/images/{uuid}.{ext}
7. 存储路径到数据库
```

#### 访问方式
- API 代理访问：`GET /images/{filename}`
- 后端读取文件并返回
- 设置正确的 Content-Type

#### 删除策略
- 产品删除时不删除图片文件（保留历史数据）
- 产品更新图片时可选择删除旧图片

### 6.2 目录结构

```
backend/
├── uploads/
│   └── image/
│       ├── 550e8400-e29b-41d4-a716-446655440000.jpg
│       ├── 660e8400-e29b-41d4-a716-446655440001.png
│       └── ...
```

## 7. 日志设计

### 7.1 应用日志

#### 日志级别
- `INFO`: 正常操作信息
- `WARNING`: 警告信息
- `ERROR`: 错误信息

#### 日志内容
- 用户登录/登出
- API 请求（路径、方法、用户、响应时间）
- 数据库操作错误
- 文件操作错误
- 业务逻辑错误

#### 存储方式
- 存储到 `app_logs` 表
- 永久保留

### 7.2 审计日志

#### 记录操作
- 创建员工
- 修改员工信息
- 创建产品
- 修改产品信息
- 删除产品
- 发放积分
- 扣除积分
- 核销订单
- 拒绝订单

#### 日志内容
- 操作人
- 操作时间
- 操作类型
- 操作对象
- 操作前数据（JSON）
- 操作后数据（JSON）

#### 存储方式
- 存储到 `audit_logs` 表
- 永久保留

## 8. 安全设计

### 8.1 密码安全

#### 密码哈希
- 使用 bcrypt 算法
- Cost factor: 10
- 不存储明文密码

#### 密码策略
- 无复杂度要求（MVP 版本）
- 无长度限制
- 无过期策略

### 8.2 会话安全

#### JWT 配置
- 签名算法：HS256
- 密钥：通过环境变量配置
- 有效期：24 小时
- 无 Refresh Token

#### Cookie 安全
- HttpOnly: true（防止 XSS）
- Secure: false（本地开发）
- SameSite: Lax（防止 CSRF）

### 8.3 API 安全

#### 认证
- 所有 API（除登录外）需要 JWT 认证
- 通过中间件验证 Token

#### 授权
- 员工只能访问员工端 API
- 管理员可以访问所有 API
- 员工只能查看/操作自己的数据

#### 输入验证
- 后端进行基本的手动验证
- 验证必填字段
- 验证数据类型
- 验证数据范围

### 8.4 文件上传安全

#### 限制
- 文件大小：≤1MB
- 文件类型：不限制（MVP 版本）
- 文件名：使用 UUID，避免路径遍历

## 9. 性能优化

### 9.1 数据库优化

#### 索引策略
- 主键索引
- 唯一索引（工号、订单号）
- 外键索引
- 状态字段索引
- 时间字段索引

#### 查询优化
- 使用分页查询（默认50条/页）
- 避免 SELECT *
- 使用 JOIN 减少查询次数

### 9.2 缓存策略

MVP 版本不使用缓存，直接查询数据库。

### 9.3 前端优化

#### 资源优化
- Webpack 打包压缩（生产环境）
- 不使用代码分割（MVP 版本）
- 不使用懒加载（MVP 版本）

#### 请求优化
- 使用分页加载
- 避免重复请求


## 10. 部署设计

### 10.1 部署架构

```
┌─────────────────────────────────────────────────────────────┐
│                      本地服务器                              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌────────────────────────────────────────────────────┐    │
│  │  前端服务 (Node.js + Express)                      │    │
│  │  端口: 8000                                         │    │
│  │  进程: node server.js                              │    │
│  └────────────────────────────────────────────────────┘    │
│                                                              │
│  ┌────────────────────────────────────────────────────┐    │
│  │  后端服务 (Go + Gin)                               │    │
│  │  端口: 8080                                         │    │
│  │  进程: ./awsomeshop                                │    │
│  └────────────────────────────────────────────────────┘    │
│                                                              │
│  ┌────────────────────────────────────────────────────┐    │
│  │  MySQL 数据库                                       │    │
│  │  端口: 3306                                         │    │
│  │  数据库: awsomeshop                                │    │
│  └────────────────────────────────────────────────────┘    │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 10.2 环境配置

#### 后端环境变量 (.env)
```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=awsomeshop

# JWT 配置
JWT_SECRET=your-secret-key-here

# 服务器配置
SERVER_PORT=8080

# 文件上传配置
UPLOAD_DIR=/uploads/image/
MAX_UPLOAD_SIZE=1048576  # 1MB in bytes

# 日志配置
LOG_LEVEL=INFO
```

#### 前端配置
API 地址写死在代码中：
```javascript
const API_BASE_URL = 'http://localhost:8080/api';
```

### 10.3 Docker Compose 配置（本地开发）

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: awsomeshop-mysql
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: awsomeshop
    ports:
      - "3306:3306"
    volumes:
      - mysql-data:/var/lib/mysql
    networks:
      - awsomeshop-network

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: awsomeshop-backend
    environment:
      DB_HOST: mysql
      DB_PORT: 3306
      DB_USER: root
      DB_PASSWORD: password
      DB_NAME: awsomeshop
      JWT_SECRET: dev-secret-key
      SERVER_PORT: 8080
    ports:
      - "8080:8080"
    depends_on:
      - mysql
    volumes:
      - ./backend/uploads:/app/uploads
    networks:
      - awsomeshop-network

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: awsomeshop-frontend
    ports:
      - "8000:8000"
    networks:
      - awsomeshop-network

volumes:
  mysql-data:

networks:
  awsomeshop-network:
    driver: bridge
```

### 10.4 生产环境部署

#### 部署步骤

1. **准备服务器**
   - 安装 MySQL 8.0
   - 安装 Node.js 18+
   - 安装 Go 1.21+

2. **部署数据库**
   ```bash
   # 创建数据库
   mysql -u root -p
   CREATE DATABASE awsomeshop CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   
   # 执行数据库迁移脚本
   mysql -u root -p awsomeshop < schema.sql
   
   # 执行种子数据脚本
   mysql -u root -p awsomeshop < seed.sql
   ```

3. **部署后端服务**
   ```bash
   # 编译 Go 应用
   cd backend
   go build -o awsomeshop main.go
   
   # 配置环境变量
   cp .env.example .env
   vim .env
   
   # 创建上传目录
   mkdir -p uploads/image
   
   # 启动服务
   ./awsomeshop
   ```

4. **部署前端服务**
   ```bash
   # 构建前端
   cd frontend
   npm install
   npm run build
   
   # 启动 Express 服务器
   node server.js
   ```

5. **验证部署**
   - 访问 http://localhost:8000
   - 使用 admin/123456 登录
   - 检查各功能是否正常

### 10.5 数据库初始化

#### 种子数据脚本 (seed.sql)

```sql
-- 创建初始管理员账号
INSERT INTO `users` (
  `employee_id`, `name`, `department`, `position`, 
  `email`, `phone`, `password`, `role`, `status`, 
  `points`, `hire_date`
) VALUES (
  'admin', '系统管理员', '管理部', '系统管理员',
  'admin@example.com', '13800000000', 
  '$2a$10$...', -- bcrypt hash of '123456'
  'admin', 'active', 0, '2026-01-01'
);

-- 创建示例产品
INSERT INTO `products` (`name`, `description`, `points_required`, `status`) VALUES
('AirPods Pro', '苹果无线降噪耳机，支持主动降噪和空间音频', 500, 'online'),
('小米手环7', '智能运动手环，支持心率监测和睡眠分析', 100, 'online'),
('星巴克咖啡券', '星巴克中杯咖啡兑换券，全国门店通用', 50, 'online'),
('京东购物卡 100元', '京东电子购物卡，可在京东商城使用', 100, 'online'),
('罗技无线鼠标', '罗技 MX Master 3 无线鼠标，办公利器', 300, 'online'),
('膳魔师保温杯', '膳魔师不锈钢保温杯 500ml，保温24小时', 150, 'online'),
('小米充电宝', '小米移动电源3 20000mAh，双向快充', 200, 'online'),
('网易云音乐年卡', '网易云音乐黑胶VIP年卡，畅享无损音质', 180, 'online'),
('kindle电子书阅读器', 'Kindle Paperwhite 电子书阅读器 8GB', 800, 'online');
```

### 10.6 监控和维护

#### 日志查看
- 应用日志：查询 `app_logs` 表
- 审计日志：查询 `audit_logs` 表

#### 数据备份
MVP 版本不需要自动备份，可手动执行：
```bash
mysqldump -u root -p awsomeshop > backup_$(date +%Y%m%d).sql
```

#### 进程管理
不使用进程管理工具，手动启动服务：
```bash
# 后端
cd backend && ./awsomeshop &

# 前端
cd frontend && node server.js &
```


## 11. 测试设计

### 11.1 测试策略

#### 测试范围
- 单元测试：后端业务逻辑
- 不包含：集成测试、API 测试、前端测试

#### 测试框架
- Go 标准 testing 包
- testify 断言库
- sqlmock 数据库 Mock

#### 测试覆盖率目标
- 60%

### 11.2 单元测试设计

#### 测试目录结构
```
backend/
├── internal/
│   ├── auth/
│   │   ├── auth.go
│   │   └── auth_test.go
│   ├── user/
│   │   ├── service.go
│   │   └── service_test.go
│   ├── product/
│   │   ├── service.go
│   │   └── service_test.go
│   ├── point/
│   │   ├── service.go
│   │   └── service_test.go
│   └── order/
│       ├── service.go
│       └── service_test.go
```

#### 测试用例示例

```go
// user/service_test.go
package user

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/DATA-DOG/go-sqlmock"
)

func TestCreateEmployee(t *testing.T) {
    // 创建 mock 数据库
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    // 设置期望的数据库操作
    mock.ExpectBegin()
    mock.ExpectExec("INSERT INTO users").
        WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectExec("INSERT INTO point_logs").
        WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectCommit()

    // 执行测试
    service := NewUserService(db)
    user, err := service.CreateEmployee(&CreateEmployeeRequest{
        Name:       "张三",
        Department: "技术部",
        Position:   "工程师",
        Email:      "zhangsan@example.com",
        Phone:      "13800138000",
        HireDate:   "2026-01-14",
        Password:   "123456",
    })

    // 断言结果
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "张三", user.Name)
    assert.Equal(t, 1000, user.Points)
    
    // 验证所有期望的数据库操作都被执行
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateEmployee_LimitReached(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    // 模拟当天已创建999个员工
    mock.ExpectQuery("SELECT COUNT").
        WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(999))

    service := NewUserService(db)
    user, err := service.CreateEmployee(&CreateEmployeeRequest{
        Name: "张三",
    })

    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "EMPLOYEE_LIMIT_REACHED")
}
```

#### 需要测试的关键功能

**用户模块**:
- 创建员工（正常、达到上限）
- 生成工号（正常、跨天重置）
- 密码哈希和验证
- 用户认证

**产品模块**:
- 创建产品
- 更新产品
- 软删除产品
- 上架/下架产品

**积分模块**:
- 发放积分
- 扣除积分
- 批量发放积分
- 记录积分日志

**订单模块**:
- 创建订单（正常、积分不足、产品下架）
- 取消订单（正常、状态不允许）
- 核销订单
- 拒绝订单
- 事务回滚测试

### 11.3 测试执行

#### 运行测试
```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/user

# 查看测试覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 12. 工号生成算法

### 12.1 工号格式
- 格式：`YYYYMMDD-NNN`
- 示例：`20260114-001`
- 说明：
  - `YYYYMMDD`: 创建日期
  - `NNN`: 当天的序号（001-999）

### 12.2 生成逻辑

```go
func GenerateEmployeeID(db *gorm.DB) (string, error) {
    today := time.Now().Format("20060102")
    
    // 查询当天已创建的员工数量
    var count int64
    err := db.Model(&User{}).
        Where("employee_id LIKE ?", today+"-%").
        Where("role = ?", "employee").
        Count(&count).Error
    
    if err != nil {
        return "", err
    }
    
    // 检查是否达到上限
    if count >= 999 {
        return "", errors.New("EMPLOYEE_LIMIT_REACHED")
    }
    
    // 生成新工号
    sequence := count + 1
    employeeID := fmt.Sprintf("%s-%03d", today, sequence)
    
    return employeeID, nil
}
```

### 12.3 特殊情况处理

#### 初始管理员
- 工号：`admin`
- 不遵循日期格式规则

#### 新创建的管理员
- 遵循员工工号格式：`YYYYMMDD-NNN`
- 与员工共享序号池

#### 并发创建
- 不考虑并发问题（MVP 版本）
- 可能出现工号冲突，由数据库唯一索引保证

## 13. 订单号生成算法

### 13.1 订单号格式
- 格式：`YYYYMMDDHHmmss-XXXXXX`
- 示例：`20260114153045-ABC123`
- 说明：
  - `YYYYMMDDHHmmss`: 创建时间戳
  - `XXXXXX`: 6位随机字符（大写字母+数字）

### 13.2 生成逻辑

```go
func GenerateOrderNo() string {
    timestamp := time.Now().Format("20060102150405")
    random := generateRandomString(6)
    return fmt.Sprintf("%s-%s", timestamp, random)
}

func generateRandomString(length int) string {
    const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}
```

## 14. 数据迁移

### 14.1 迁移策略
- 使用 SQL 脚本进行数据库迁移
- 不使用 ORM 自动迁移

### 14.2 迁移脚本

#### schema.sql
包含所有表结构定义（见 3.2 节）

#### seed.sql
包含初始数据（见 10.5 节）

### 14.3 迁移执行
```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE awsomeshop CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 执行表结构迁移
mysql -u root -p awsomeshop < schema.sql

# 执行种子数据迁移
mysql -u root -p awsomeshop < seed.sql
```


## 15. 项目结构

### 15.1 后端项目结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # 应用入口
├── internal/
│   ├── auth/                       # 认证模块
│   │   ├── handler.go              # HTTP 处理器
│   │   ├── service.go              # 业务逻辑
│   │   ├── middleware.go           # JWT 中间件
│   │   └── auth_test.go            # 单元测试
│   ├── user/                       # 用户模块
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go           # 数据访问层
│   │   ├── model.go                # 数据模型
│   │   └── service_test.go
│   ├── product/                    # 产品模块
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── model.go
│   │   └── service_test.go
│   ├── point/                      # 积分模块
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── model.go
│   │   └── service_test.go
│   ├── order/                      # 订单模块
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── model.go
│   │   └── service_test.go
│   ├── audit/                      # 审计日志模块
│   │   ├── service.go
│   │   └── model.go
│   ├── logger/                     # 应用日志模块
│   │   ├── service.go
│   │   └── model.go
│   └── common/                     # 公共模块
│       ├── response.go             # 统一响应格式
│       ├── errors.go               # 错误码定义
│       └── utils.go                # 工具函数
├── pkg/
│   ├── database/                   # 数据库连接
│   │   └── mysql.go
│   └── config/                     # 配置管理
│       └── config.go
├── uploads/                        # 文件上传目录
│   └── image/
├── migrations/                     # 数据库迁移脚本
│   ├── schema.sql
│   └── seed.sql
├── .env.example                    # 环境变量示例
├── .env                            # 环境变量（不提交）
├── go.mod                          # Go 模块定义
├── go.sum                          # Go 依赖锁定
├── Dockerfile                      # Docker 镜像构建
└── README.md                       # 项目说明
```

### 15.2 前端项目结构

```
frontend/
├── src/
│   ├── js/
│   │   ├── app.js                  # 应用入口
│   │   ├── router.js               # 路由管理
│   │   ├── api.js                  # API 调用封装
│   │   ├── auth.js                 # 认证逻辑
│   │   ├── utils.js                # 工具函数
│   │   ├── pages/                  # 页面组件
│   │   │   ├── login.js
│   │   │   ├── employee/
│   │   │   │   ├── products.js
│   │   │   │   ├── product-detail.js
│   │   │   │   ├── profile.js
│   │   │   │   └── orders.js
│   │   │   └── admin/
│   │   │       ├── dashboard.js
│   │   │       ├── employees.js
│   │   │       ├── products.js
│   │   │       ├── points.js
│   │   │       └── orders.js
│   │   └── components/             # 可复用组件
│   │       ├── navbar.js
│   │       ├── toast.js
│   │       └── pagination.js
│   ├── css/
│   │   ├── antd.min.css            # Ant Design 样式
│   │   └── custom.css              # 自定义样式
│   └── index.html                  # HTML 模板
├── dist/                           # 构建产物（不提交）
│   ├── index.html
│   ├── js/
│   │   └── bundle.js
│   └── css/
│       └── style.css
├── server.js                       # Express 服务器
├── webpack.config.js               # Webpack 配置
├── package.json                    # NPM 依赖
├── package-lock.json               # NPM 依赖锁定
├── Dockerfile                      # Docker 镜像构建
└── README.md                       # 项目说明
```

### 15.3 根目录结构

```
awsomeshop/
├── backend/                        # 后端项目
├── frontend/                       # 前端项目
├── docker-compose.yml              # Docker Compose 配置
└── README.md                       # 项目总体说明
```

## 16. 开发规范

### 16.1 代码规范

#### Go 代码规范
- 遵循 Go 官方代码规范
- 使用 gofmt 格式化代码
- 使用 golint 检查代码质量
- 函数命名：驼峰命名法
- 常量命名：大写字母+下划线

#### JavaScript 代码规范
- 使用 ES6+ 语法
- 使用 const/let，避免 var
- 函数命名：驼峰命名法
- 常量命名：大写字母+下划线

### 16.2 Git 规范

#### 分支策略
- `main`: 主分支，稳定版本
- `develop`: 开发分支
- `feature/*`: 功能分支
- `bugfix/*`: 修复分支

#### Commit 规范
```
<type>: <subject>

<body>

<footer>
```

类型（type）:
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建/工具相关

示例:
```
feat: 实现员工创建功能

- 添加员工创建 API
- 实现工号自动生成
- 添加入职积分发放逻辑

Closes #123
```

### 16.3 API 文档规范

使用 Swagger/OpenAPI 生成 API 文档：

```go
// @Summary 创建员工
// @Description 管理员创建新员工账号，自动发放1000积分
// @Tags 员工管理
// @Accept json
// @Produce json
// @Param request body CreateEmployeeRequest true "员工信息"
// @Success 201 {object} Response{data=User}
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/admin/employees [post]
// @Security BearerAuth
func (h *UserHandler) CreateEmployee(c *gin.Context) {
    // ...
}
```

## 17. 设计决策记录

### 17.1 架构决策

| 决策 | 原因 | 影响 |
|------|------|------|
| 前后端分离 | 职责清晰，便于维护 | 需要处理跨域问题 |
| 不使用 API 网关 | MVP 版本，架构简单 | 后续扩展需要重构 |
| 统一用户表 | 简化数据模型 | 通过 role 字段区分角色 |
| 不使用缓存 | MVP 版本，用户量小 | 性能可能受限 |
| 不考虑并发 | MVP 版本，简化实现 | 可能出现数据不一致 |

### 17.2 技术选型决策

| 技术 | 选择 | 原因 |
|------|------|------|
| 后端语言 | Go | 性能好，并发支持强 |
| 后端框架 | Gin | 轻量级，性能优秀 |
| ORM | GORM | Go 生态最成熟的 ORM |
| 前端 | 原生 JS | MVP 版本，避免框架复杂度 |
| UI | Ant Design CSS | 成熟的设计系统 |
| 数据库 | MySQL | 成熟稳定，易于维护 |
| 认证 | JWT | 无状态，易于扩展 |

## 18. 风险和限制

### 18.1 技术风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 并发问题 | 数据不一致 | 后续版本添加锁机制 |
| 工号冲突 | 创建失败 | 数据库唯一索引保证 |
| 文件存储 | 磁盘空间不足 | 定期清理无用文件 |
| 无备份 | 数据丢失 | 手动定期备份 |

### 18.2 功能限制

| 限制 | 说明 | 后续计划 |
|------|------|----------|
| 每天最多创建999个员工 | 工号格式限制 | 扩展为4位序号 |
| 无缓存 | 性能受限 | 添加 Redis 缓存 |
| 无搜索功能 | 产品列表无法搜索 | 添加搜索功能 |
| 无数据导出 | 无法导出报表 | 添加导出功能 |

## 19. 后续优化建议

### 19.1 性能优化
- 添加 Redis 缓存（产品列表、用户信息）
- 数据库读写分离
- 图片使用 CDN 加速
- 前端资源压缩和懒加载

### 19.2 功能增强
- 产品搜索和筛选
- 数据导出（Excel）
- 邮件通知（订单审核结果）
- 移动端适配
- 产品分类管理
- 订单批量审核

### 19.3 安全增强
- 密码复杂度要求
- 登录失败次数限制
- 操作日志审计查询界面
- 文件类型白名单
- HTTPS 支持

### 19.4 运维增强
- 自动化部署脚本
- 数据库自动备份
- 监控告警系统
- 日志聚合分析
- 性能监控

## 20. 审批记录

| 角色 | 姓名 | 审批状态 | 审批日期 | 备注 |
|------|------|----------|----------|------|
| 技术负责人 | 用户 | 已批准 | 2026-01-14 | |
| 架构师 | 用户 | 已批准 | 2026-01-14 | |
| 产品经理 | 用户 | 已批准 | 2026-01-14 | |
| 项目经理 | 用户 | 已批准 | 2026-01-14 | |

---

**文档结束**

