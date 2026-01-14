-- 数据库验证脚本
-- 用于验证数据库和表是否正确创建

USE `awsomeshop`;

-- 验证所有表是否存在
SELECT 'Checking tables...' AS status;

SELECT 
  TABLE_NAME,
  TABLE_ROWS,
  TABLE_COMMENT
FROM 
  INFORMATION_SCHEMA.TABLES
WHERE 
  TABLE_SCHEMA = 'awsomeshop'
ORDER BY 
  TABLE_NAME;

-- 验证管理员账号是否创建
SELECT 'Checking admin user...' AS status;

SELECT 
  id,
  employee_id,
  name,
  role,
  status,
  points,
  hire_date
FROM 
  users
WHERE 
  employee_id = 'admin';

-- 验证产品数量
SELECT 'Checking products...' AS status;

SELECT 
  COUNT(*) AS total_products,
  SUM(CASE WHEN status = 'online' THEN 1 ELSE 0 END) AS online_products,
  SUM(CASE WHEN status = 'offline' THEN 1 ELSE 0 END) AS offline_products
FROM 
  products
WHERE 
  is_deleted = 0;

-- 显示所有产品
SELECT 
  id,
  name,
  points_required,
  status
FROM 
  products
WHERE 
  is_deleted = 0
ORDER BY 
  id;

SELECT 'Database validation complete!' AS status;
