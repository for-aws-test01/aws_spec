# 数据库迁移脚本

本目录包含 AWSomeShop 项目的数据库迁移脚本。

## 文件说明

- `schema.sql`: 数据库架构脚本，创建所有表结构
- `seed.sql`: 种子数据脚本，插入初始数据

## 执行顺序

1. 首先执行 `schema.sql` 创建数据库和表结构
2. 然后执行 `seed.sql` 插入初始数据

## 使用方法

### 方法一：使用自动化设置脚本（推荐）

```bash
# 进入 migrations 目录
cd awsomeshop/backend/migrations

# 运行设置脚本
./setup.sh
```

脚本会引导你输入数据库连接信息，并自动执行以下操作：
1. 测试数据库连接
2. 创建数据库和表结构
3. 插入种子数据
4. 验证数据库设置

### 方法二：使用 MySQL 命令行

```bash
# 1. 创建数据库和表结构
mysql -u root -p < schema.sql

# 2. 插入种子数据
mysql -u root -p < seed.sql

# 3. 验证数据库（可选）
mysql -u root -p < validate.sql
```

### 方法三：使用 MySQL 客户端

```bash
# 连接到 MySQL
mysql -u root -p

# 在 MySQL 命令行中执行
source /path/to/schema.sql;
source /path/to/seed.sql;
source /path/to/validate.sql;
```

### 方法四：使用 Docker（如果使用 Docker Compose）

```bash
# 假设 MySQL 容器名为 mysql
docker exec -i mysql mysql -u root -p<password> < schema.sql
docker exec -i mysql mysql -u root -p<password> < seed.sql
docker exec -i mysql mysql -u root -p<password> < validate.sql
```

## 初始数据

### 管理员账号

执行 `seed.sql` 后，系统会创建一个初始管理员账号：

- **工号**: `admin`
- **密码**: `admin123`
- **角色**: 管理员
- **邮箱**: admin@awsomeshop.com

### 示例产品

系统会创建 9 个示例产品：

1. AirPods Pro (1500 积分)
2. Kindle Paperwhite (800 积分)
3. 星巴克咖啡券 (100 积分)
4. 小米手环8 (200 积分)
5. 罗技无线鼠标 (600 积分)
6. 膳魔师保温杯 (150 积分)
7. 京东购物卡 (500 积分)
8. 网易云音乐年卡 (300 积分)
9. 优衣库购物券 (200 积分)

**注意**: 所有产品默认状态为"下架"（offline），需要管理员手动上架。

## 数据库配置

- **数据库名**: `awsomeshop`
- **字符集**: `utf8mb4`
- **排序规则**: `utf8mb4_unicode_ci`
- **存储引擎**: `InnoDB`

## 表结构

1. `users` - 用户表（员工+管理员）
2. `products` - 产品表
3. `orders` - 订单表
4. `point_logs` - 积分日志表
5. `audit_logs` - 审计日志表
6. `app_logs` - 应用日志表

## 注意事项

1. 执行脚本前请确保 MySQL 服务已启动
2. 确保有足够的权限创建数据库和表
3. 如果数据库已存在，脚本会跳过创建步骤
4. 管理员密码使用 bcrypt 加密存储
5. 建议在生产环境中修改默认管理员密码

## 重置数据库

如果需要重置数据库，可以执行以下命令：

```sql
DROP DATABASE IF EXISTS awsomeshop;
```

然后重新执行 `schema.sql` 和 `seed.sql`。
