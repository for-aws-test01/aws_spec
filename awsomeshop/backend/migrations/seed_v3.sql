-- AWSomeShop 种子数据 v3.0
-- 版本: 3.0
-- 创建日期: 2026-01-14

USE `awsomeshop`;

-- ============================================================================
-- 插入初始管理员账号
-- ============================================================================
-- 工号: admin
-- 密码: 123456 (bcrypt hash)
INSERT INTO `users` (
  `employee_id`,
  `name`,
  `department`,
  `position`,
  `email`,
  `phone`,
  `password`,
  `role`,
  `points`,
  `hire_date`
) VALUES (
  'admin',
  '系统管理员',
  '技术部',
  '系统管理员',
  'admin@awsomeshop.com',
  '13800000000',
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
  'admin',
  0,
  '2026-01-01'
);

-- ============================================================================
-- 插入示例产品（9个）
-- ============================================================================
INSERT INTO `products` (`name`, `description`, `image_url`, `points_required`) VALUES
('AirPods Pro', '苹果无线降噪耳机，主动降噪，通透模式，空间音频', 'https://store.storeimages.cdn-apple.com/8756/as-images.apple.com/is/MQD83?wid=1144&hei=1144&fmt=jpeg&qlt=90', 500),
('Kindle Paperwhite', '亚马逊电子书阅读器，6.8英寸屏幕，防水设计', 'https://m.media-amazon.com/images/I/51QCk82iGcL._AC_SL1000_.jpg', 300),
('小米手环8', '智能运动手环，心率监测，睡眠追踪，50米防水', 'https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/7f3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c.jpg', 150),
('星巴克咖啡券', '星巴克100元电子礼品卡，全国门店通用', 'https://www.starbucks.com.cn/images/gift-card.jpg', 100),
('优衣库购物券', '优衣库200元购物券，线上线下通用', 'https://www.uniqlo.cn/images/gift-card.jpg', 180),
('京东E卡', '京东500元电子卡，可购买自营商品', 'https://img14.360buyimg.com/n1/jfs/t1/gift-card.jpg', 450),
('罗技无线鼠标', '罗技MX Master 3无线鼠标，人体工学设计', 'https://resource.logitech.com/content/dam/logitech/en/products/mice/mx-master-3.png', 250),
('膳魔师保温杯', '膳魔师不锈钢保温杯，500ml，24小时保温', 'https://www.thermos.com/images/products/thermos-bottle.jpg', 120),
('网易云音乐年卡', '网易云音乐黑胶VIP年卡，畅听无损音乐', 'https://p5.music.126.net/obj/wo3DlcOGw6DClTvDisK1/vip-card.jpg', 200);
