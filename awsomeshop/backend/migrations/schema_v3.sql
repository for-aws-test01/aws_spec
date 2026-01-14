-- AWSomeShop 数据库架构 v3.0 (超级简化版)
-- 版本: 3.0
-- 创建日期: 2026-01-14
-- 说明: 删除 audit_logs 和 app_logs 表，简化字段

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `awsomeshop` 
  DEFAULT CHARACTER SET utf8mb4 
  COLLATE utf8mb4_unicode_ci;

USE `awsomeshop`;

-- ============================================================================
-- 用户表 (users)
-- 存储员工和管理员账号信息（统一表）
-- ============================================================================
CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `employee_id` VARCHAR(50) NOT NULL COMMENT '工号',
  `name` VARCHAR(100) NOT NULL COMMENT '姓名',
  `department` VARCHAR(100) NOT NULL COMMENT '部门',
  `position` VARCHAR(100) NOT NULL COMMENT '岗位',
  `email` VARCHAR(100) NOT NULL COMMENT '邮箱',
  `phone` VARCHAR(20) NOT NULL COMMENT '手机号',
  `password` VARCHAR(255) NOT NULL COMMENT '密码哈希',
  `role` ENUM('employee', 'admin') NOT NULL DEFAULT 'employee' COMMENT '角色',
  `points` INT NOT NULL DEFAULT 0 COMMENT '积分余额',
  `hire_date` DATE NOT NULL COMMENT '入职时间',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_employee_id` (`employee_id`),
  KEY `idx_role` (`role`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表（员工+管理员）';

-- ============================================================================
-- 产品表 (products)
-- 存储产品信息
-- ============================================================================
CREATE TABLE IF NOT EXISTS `products` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '产品ID',
  `name` VARCHAR(100) NOT NULL COMMENT '产品名称',
  `description` TEXT NOT NULL COMMENT '产品描述',
  `image_url` VARCHAR(500) DEFAULT NULL COMMENT '产品图片URL（外部链接）',
  `points_required` INT NOT NULL COMMENT '所需积分',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='产品表';

-- ============================================================================
-- 订单表 (orders)
-- 存储兑换订单信息
-- ============================================================================
CREATE TABLE IF NOT EXISTS `orders` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `order_no` VARCHAR(50) NOT NULL COMMENT '订单号',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `product_id` BIGINT UNSIGNED NOT NULL COMMENT '产品ID',
  `product_snapshot` JSON NOT NULL COMMENT '产品快照',
  `points_cost` INT NOT NULL COMMENT '消耗积分',
  `status` ENUM('pending', 'approved') NOT NULL DEFAULT 'pending' COMMENT '订单状态',
  `applied_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
  `reviewed_at` TIMESTAMP NULL DEFAULT NULL COMMENT '审核时间',
  `reviewer_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '审核人ID',
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

-- ============================================================================
-- 积分日志表 (point_logs)
-- 记录所有积分变动操作
-- ============================================================================
CREATE TABLE IF NOT EXISTS `point_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `operator_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '操作人ID',
  `amount` INT NOT NULL COMMENT '积分变动数量（正数增加，负数减少）',
  `reason` TEXT NOT NULL COMMENT '变动原因',
  `type` ENUM('onboarding', 'admin_grant', 'redeem') NOT NULL COMMENT '变动类型',
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

-- ============================================================================
-- 变更说明 v3.0
-- ============================================================================
-- 1. 删除 audit_logs 表（审计日志）
-- 2. 删除 app_logs 表（应用日志）
-- 3. users 表删除 status 字段（不需要状态管理）
-- 4. products 表删除 status 和 is_deleted 字段（产品创建后立即可用，不删除）
-- 5. orders 表只保留 pending 和 approved 两种状态
-- 6. orders 表删除 reject_reason 和 approval_note 字段
-- 7. point_logs 表只保留 3 种类型：onboarding, admin_grant, redeem
