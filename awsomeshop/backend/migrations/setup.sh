#!/bin/bash

# AWSomeShop 数据库设置脚本
# 用于自动化数据库初始化过程

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# 检查 MySQL 是否安装
check_mysql() {
    if ! command -v mysql &> /dev/null; then
        print_error "MySQL/MariaDB is not installed"
        exit 1
    fi
    print_info "MySQL/MariaDB found: $(mysql --version)"
}

# 获取数据库连接信息
get_db_credentials() {
    echo ""
    print_info "Please enter MySQL connection details:"
    
    read -p "MySQL Host (default: localhost): " DB_HOST
    DB_HOST=${DB_HOST:-localhost}
    
    read -p "MySQL Port (default: 3306): " DB_PORT
    DB_PORT=${DB_PORT:-3306}
    
    read -p "MySQL User (default: root): " DB_USER
    DB_USER=${DB_USER:-root}
    
    read -sp "MySQL Password: " DB_PASSWORD
    echo ""
}

# 测试数据库连接
test_connection() {
    print_info "Testing database connection..."
    
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1" &> /dev/null; then
        print_info "Database connection successful"
        return 0
    else
        print_error "Failed to connect to database"
        return 1
    fi
}

# 执行 SQL 脚本
execute_sql() {
    local sql_file=$1
    local description=$2
    
    print_info "$description"
    
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" < "$sql_file"; then
        print_info "$description - Success"
        return 0
    else
        print_error "$description - Failed"
        return 1
    fi
}

# 主函数
main() {
    echo "========================================"
    echo "  AWSomeShop Database Setup"
    echo "========================================"
    echo ""
    
    # 检查 MySQL
    check_mysql
    
    # 获取连接信息
    get_db_credentials
    
    # 测试连接
    if ! test_connection; then
        exit 1
    fi
    
    echo ""
    print_warning "This will create/reset the 'awsomeshop' database"
    read -p "Continue? (y/N): " confirm
    
    if [[ ! $confirm =~ ^[Yy]$ ]]; then
        print_info "Setup cancelled"
        exit 0
    fi
    
    echo ""
    
    # 执行架构脚本
    if ! execute_sql "schema.sql" "Creating database schema..."; then
        exit 1
    fi
    
    # 执行种子数据脚本
    if ! execute_sql "seed.sql" "Inserting seed data..."; then
        exit 1
    fi
    
    # 执行验证脚本
    echo ""
    print_info "Validating database setup..."
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" < validate.sql
    
    echo ""
    echo "========================================"
    print_info "Database setup completed successfully!"
    echo "========================================"
    echo ""
    print_info "Initial admin credentials:"
    echo "  Employee ID: admin"
    echo "  Password: admin123"
    echo ""
    print_warning "Please change the admin password after first login!"
}

# 运行主函数
main
