# AWSomeShop v3.0 变更日志

## 版本信息
- **版本**: v3.0 超级简化版
- **发布日期**: 2026-01-14
- **变更类型**: 重大简化

## 变更概述

v3.0 版本在 v2.0 基础上进一步简化，删除 35 个任务，工作量从 50 天减少到 49 天，相比 v1.0 减少 70%。

## 工作量对比

| 版本 | 任务数 | 工作量（天） | 相比v1.0减少 |
|------|--------|--------------|--------------|
| v1.0 原始版 | 85 | 141 | - |
| v2.0 激进简化 | 69 | 50 | 65% |
| v3.0 超级简化 | 34 | 49 | 70% |

## 删除的功能（35项）

### 1. 测试相关（5个任务，10天）
- ❌ TASK-071: 编写用户模块单元测试
- ❌ TASK-072: 编写产品模块单元测试
- ❌ TASK-073: 编写积分模块单元测试
- ❌ TASK-074: 编写订单模块单元测试
- ❌ TASK-075: 执行测试并达到覆盖率目标
- ❌ TASK-083: 性能测试

**理由**: MVP 阶段使用手动测试即可

### 2. 日志系统（2个任务，2天）
- ❌ TASK-006: 实现日志模块（应用日志 + 审计日志）
- ❌ 删除 audit_logs 表
- ❌ 删除 app_logs 表

**理由**: MVP 阶段使用简单的 console 输出

### 3. 图片管理（2个任务，3天）
- ❌ TASK-019: 实现图片上传功能
- ❌ TASK-020: 实现图片访问 API

**理由**: 改用图片 URL 输入，简化实现

### 4. 系统概况（2个任务，4天）
- ❌ TASK-042: 实现系统概况 API
- ❌ TASK-060: 实现系统概况页

**理由**: 管理员直接进入订单列表页

### 5. 个人信息（2个任务，1.5天）
- ❌ TASK-043: 实现个人信息查询 API
- ❌ TASK-055: 实现个人信息页

**理由**: 员工只需要看积分和订单

### 6. 密码管理（3个任务，4天）
- ❌ TASK-011: 实现修改密码 API
- ❌ TASK-017: 实现重置员工密码 API
- ❌ TASK-059: 实现修改密码页

**理由**: MVP 阶段使用固定密码

### 7. 积分管理（5个任务，7天）
- ❌ TASK-029: 实现积分日志查询 API（员工）
- ❌ TASK-030: 实现积分日志查询 API（管理员）
- ❌ TASK-032: 实现批量发放积分 API
- ❌ TASK-033: 实现扣除积分 API
- ❌ TASK-068: 实现积分日志页（管理员）

**理由**: 只保留单个发放积分功能

### 8. 员工管理（3个任务，4.5天）
- ❌ TASK-015: 实现员工详情查询 API
- ❌ TASK-016: 实现更新员工信息 API
- ❌ TASK-063: 实现员工详情/编辑页

**理由**: 只保留创建和列表功能

### 9. 产品管理（4个任务，6天）
- ❌ TASK-024: 实现产品详情查询 API（管理员）
- ❌ TASK-025: 实现更新产品 API
- ❌ TASK-026: 实现删除产品 API
- ❌ TASK-027: 实现产品上架/下架 API
- ❌ TASK-066: 实现产品详情/编辑页（管理员）

**理由**: 只保留创建和列表功能，产品创建后立即可用

### 10. 订单管理（4个任务，6.5天）
- ❌ TASK-037: 实现订单详情查询 API（员工）
- ❌ TASK-038: 实现取消订单 API
- ❌ TASK-041: 实现拒绝订单 API
- ❌ TASK-058: 实现订单详情页（员工）
- ❌ TASK-070: 实现订单详情页（管理员）

**理由**: 只保留核销功能，订单提交后不可取消

### 11. 管理员管理（1个任务，1天）
- ❌ TASK-018: 实现创建管理员 API

**理由**: 只使用初始 admin 账号

### 12. 部署和文档（3个任务，6天）
- ❌ TASK-078: 配置 Docker Compose
- ❌ TASK-080: 编写 API 文档（Swagger）

**理由**: 简化部署流程，使用简单的 README

## 数据库变更

### 删除的表（2个）
1. ❌ `audit_logs` - 审计日志表
2. ❌ `app_logs` - 应用日志表

### 简化的表结构

#### users 表
- ❌ 删除 `status` 字段（不需要状态管理）
- ✅ 保留 `role` 字段（employee/admin）

#### products 表
- ❌ 删除 `status` 字段（不需要上架/下架）
- ❌ 删除 `is_deleted` 字段（不需要软删除）
- ✅ `image_url` 改为外部链接（VARCHAR(500)）

#### orders 表
- ❌ 删除 `cancelled` 和 `rejected` 状态
- ❌ 删除 `reject_reason` 字段
- ❌ 删除 `approval_note` 字段
- ✅ 只保留 `pending` 和 `approved` 状态

#### point_logs 表
- ❌ 删除 `admin_deduct` 类型
- ❌ 删除 `order_cancel_refund` 类型
- ❌ 删除 `order_reject_refund` 类型
- ✅ 只保留 3 种类型：`onboarding`, `admin_grant`, `redeem`

## 保留的核心功能

### 员工端（4个页面）
1. ✅ 登录页
2. ✅ 产品列表页（首页）
3. ✅ 产品详情页（兑换）
4. ✅ 订单列表页
5. ✅ 积分余额（导航栏显示）

### 管理员端（5个页面）
1. ✅ 登录页
2. ✅ 员工列表页 + 创建员工
3. ✅ 产品列表页 + 创建产品
4. ✅ 订单列表页（直接核销）
5. ✅ 发放积分页

### 后端 API（14个接口）
1. ✅ POST /api/auth/login - 登录
2. ✅ POST /api/auth/logout - 登出
3. ✅ POST /api/admin/employees - 创建员工
4. ✅ GET /api/admin/employees - 员工列表
5. ✅ POST /api/admin/products - 创建产品
6. ✅ GET /api/admin/products - 产品列表（管理员）
7. ✅ GET /api/products - 产品列表（员工）
8. ✅ GET /api/products/:id - 产品详情
9. ✅ GET /api/points/balance - 积分余额
10. ✅ POST /api/admin/points/grant - 发放积分
11. ✅ POST /api/orders/redeem - 兑换产品
12. ✅ GET /api/orders - 订单列表（员工）
13. ✅ GET /api/admin/orders - 订单列表（管理员）
14. ✅ PUT /api/admin/orders/:id/approve - 核销订单

### 数据库（4个表）
1. ✅ users - 用户表（员工+管理员）
2. ✅ products - 产品表
3. ✅ orders - 订单表
4. ✅ point_logs - 积分日志表

## 业务流程简化

### 员工兑换流程
1. 员工登录 → 产品列表
2. 点击产品 → 产品详情
3. 点击兑换 → 扣除积分 → 创建订单
4. 查看订单列表

### 管理员核销流程
1. 管理员登录 → 订单列表
2. 点击核销按钮 → 订单状态变为已核销

### 积分发放流程
1. 管理员 → 发放积分页
2. 选择员工 → 输入积分和原因 → 提交

## 技术栈变更

### 前端
- ❌ 删除 Ant Design CSS
- ✅ 使用原生 CSS

### 后端
- ❌ 删除日志模块
- ❌ 删除审计日志模块
- ❌ 删除图片上传模块

## 迁移指南

### 从 v2.0 升级到 v3.0

1. **数据库迁移**
```sql
-- 删除不需要的表
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS app_logs;

-- 修改 users 表
ALTER TABLE users DROP COLUMN status;

-- 修改 products 表
ALTER TABLE products DROP COLUMN status;
ALTER TABLE products DROP COLUMN is_deleted;
ALTER TABLE products MODIFY image_url VARCHAR(500);

-- 修改 orders 表
ALTER TABLE orders MODIFY status ENUM('pending', 'approved');
ALTER TABLE orders DROP COLUMN reject_reason;
ALTER TABLE orders DROP COLUMN approval_note;

-- 修改 point_logs 表
ALTER TABLE point_logs MODIFY type ENUM('onboarding', 'admin_grant', 'redeem');
```

2. **代码迁移**
- 删除所有日志相关代码
- 删除图片上传相关代码
- 删除密码修改相关代码
- 删除订单取消/拒绝相关代码
- 简化产品管理代码

3. **前端迁移**
- 删除所有详情页
- 删除个人信息页
- 删除系统概况页
- 删除修改密码页
- 简化导航栏

## 注意事项

1. **不可逆变更**: 删除的表和字段无法恢复
2. **数据丢失**: 审计日志和应用日志数据将丢失
3. **功能限制**: 删除的功能无法使用
4. **测试**: 需要手动测试所有功能

## 后续计划

v3.0 是 MVP 的最终版本，后续版本可能会：
- 添加单元测试
- 恢复部分管理功能
- 添加数据统计
- 优化用户体验

---

**变更完成日期**: 2026-01-14  
**审批状态**: 待审批
