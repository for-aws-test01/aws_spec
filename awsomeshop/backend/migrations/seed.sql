-- AWSomeShop 种子数据
-- 版本: 1.0
-- 创建日期: 2026-01-14
-- 说明: 插入初始管理员账号和示例产品数据

USE `awsomeshop`;

-- ============================================================================
-- 插入初始管理员账号
-- 工号: admin
-- 密码: admin123 (bcrypt 哈希值)
-- ============================================================================
INSERT INTO `users` (
  `employee_id`,
  `name`,
  `department`,
  `position`,
  `email`,
  `phone`,
  `password`,
  `role`,
  `status`,
  `points`,
  `hire_date`
) VALUES (
  'admin',
  '系统管理员',
  '技术部',
  '系统管理员',
  'admin@awsomeshop.com',
  '13800000000',
  -- bcrypt hash for 'admin123'
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
  'admin',
  'active',
  0,
  '2026-01-01'
);

-- ============================================================================
-- 插入示例产品数据（9个产品）
-- 所有产品默认状态为 offline（下架）
-- ============================================================================

-- 产品 1: AirPods Pro
INSERT INTO `products` (
  `name`,
  `description`,
  `image_url`,
  `points_required`,
  `status`,
  `is_deleted`
) VALUES (
  'AirPods Pro',
  'Apple AirPods Pro 第二代，主动降噪无线耳机，支持空间音频，USB-C充电盒',
  NULL,
  1500,
  'offline',
  0
);

-- 产品 2: Kindle Paperwhite
INSERT INTO `products` (
  `name`,
  `description`,
  `image_url`,
  `points_required`,
  `status`,
  `is_deleted`
) VALUES (
  'Kindle Paperwhite',
  'Amazon Kindle Paperwhite 电子书阅读器，6.8英寸显示屏，防水设计，16GB存储',
  NULL,
  800,
  'offline',
  0
);

-- 产品 3: 星巴克咖啡券
INSERT INTO `products` (
  `name`,
  `description`,
  `image_url`,
  `points_required`,
  `status`,
  `is_deleted`
) VALUES (
  '星巴克咖啡券',
  '星巴克100元电子礼品卡，可在全国门店使用，有效期一年',
  NULL,
  100,
  'offline',
  0
);

-- 产品 4: 小米手环8
INSERT INTO `products` (
  `name`,
  `description`,
  `image_url`,
  `points_required`,
  `status`,
  `is_deleted`
) VALUES (
  '小米手环8',
  '小米手环8 智能运动手环，心率监测，睡眠监测，50米防水，续航16天',
  NULL,
  200,
  'offline',
  0
);

-- 产品 5: 罗技无线鼠标
INSERT INTO `products` (
  `name`,
  `description`,
  `image_url`,
  `points_required`,
  `status`,
  `is_deleted`
) VALUES (
  '罗技无线鼠标',
  '罗技 MX Master 3S 无线鼠标，人体工学设计，8000DPI，支持多设备切换',
  NULL,
  600,
  'offline',
  0
);

-- 产品 6: 膳魔师保温杯
INSERT INTO `products` (
  `name`,
  `description`,
  `image_url`,
  `points_required`,
  `status`,
  `is_deleted`
) VALUES (
  '膳魔师保温杯',
  '膳魔师不锈钢保温杯，500ml容量，24小时保温，防漏设计',
  NULL,
  150,
  'offline',
  0
);

-- 产品 7: 京东购物卡
INSERT INTO `products` (
  `name`,
  `description`,
  `image_url`,
  `points_required`,
  `status`,
  `is_deleted`
) VALUES (
  '京东购物卡',
  '京东500元电子购物卡，可购买京东自营商品，有效期三年',
  NULL,
  500,
  'offline',
  0
);

-- 产品 8: 网易云音乐年卡
INSERT INTO `products` (
  `name`,
  `description`,
  `image_url`,
  `points_required`,
  `status`,
  `is_deleted`
) VALUES (
  '网易云音乐年卡',
  '网易云音乐黑胶VIP年卡，畅听千万曲库，无损音质，下载歌曲',
  NULL,
  300,
  'offline',
  0
);

-- 产品 9: 优衣库购物券
INSERT INTO `products` (
  `name`,
  `description`,
  `image_url`,
  `points_required`,
  `status`,
  `is_deleted`
) VALUES (
  '优衣库购物券',
  '优衣库200元电子购物券，可在全国门店及官网使用，有效期一年',
  NULL,
  200,
  'offline',
  0
);

-- ============================================================================
-- 数据插入完成
-- ============================================================================
