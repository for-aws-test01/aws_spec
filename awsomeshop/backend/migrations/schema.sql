-- AWSomeShop 数据库架构
-- 版本: 1.0
-- 创建日期: 2026-01-14

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

-- ============================================================================
-- 产品表 (products)
-- 存储产品信息
-- ============================================================================
CREATE TABLE IF NOT EXISTS `products` (
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

-- ============================================================================
-- 审计日志表 (audit_logs)
-- 记录管理员操作的审计日志
-- ============================================================================
CREATE TABLE IF NOT EXISTS `audit_logs` (
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

-- ============================================================================
-- 应用日志表 (app_logs)
-- 记录应用运行日志
-- ============================================================================
CREATE TABLE IF NOT EXISTS `app_logs` (
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
