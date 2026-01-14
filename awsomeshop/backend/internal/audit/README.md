# Audit Service

## Overview

The audit service provides functionality to record user operations in the system. It tracks who performed what action, when, and what data changed.

## Features

- Records all administrative operations (create, update, delete)
- Stores before and after data snapshots in JSON format
- Supports multiple operation types and target types
- Provides convenient helper methods for common operations

## Usage

### Initialize the Service

```go
import (
    "awsomeshop/backend/internal/audit"
    "gorm.io/gorm"
)

// Create audit service instance
auditService := audit.NewService(db)
```

### Log Operations

#### Generic Log Method

```go
// Log any operation with before/after data
err := auditService.Log(
    operatorID,              // Who performed the operation
    audit.OperationUpdateEmployee,  // What operation
    audit.TargetTypeUser,    // What type of object
    targetID,                // Which object
    beforeData,              // Data before operation (optional)
    afterData,               // Data after operation (optional)
)
```

#### Convenience Methods

```go
// Create employee
auditService.LogCreateEmployee(operatorID, employeeID, employeeData)

// Update employee
auditService.LogUpdateEmployee(operatorID, employeeID, oldData, newData)

// Create product
auditService.LogCreateProduct(operatorID, productID, productData)

// Update product
auditService.LogUpdateProduct(operatorID, productID, oldData, newData)

// Delete product
auditService.LogDeleteProduct(operatorID, productID, productData)

// Grant points
auditService.LogGrantPoints(operatorID, userID, pointsData)

// Deduct points
auditService.LogDeductPoints(operatorID, userID, pointsData)

// Approve order
auditService.LogApproveOrder(operatorID, orderID, oldData, newData)

// Reject order
auditService.LogRejectOrder(operatorID, orderID, oldData, newData)
```

## Operation Types

- `OperationCreateEmployee` - 创建员工
- `OperationUpdateEmployee` - 修改员工
- `OperationCreateProduct` - 创建产品
- `OperationUpdateProduct` - 修改产品
- `OperationDeleteProduct` - 删除产品
- `OperationGrantPoints` - 发放积分
- `OperationDeductPoints` - 扣除积分
- `OperationApproveOrder` - 核销订单
- `OperationRejectOrder` - 拒绝订单

## Target Types

- `TargetTypeUser` - 用户
- `TargetTypeProduct` - 产品
- `TargetTypeOrder` - 订单
- `TargetTypePoint` - 积分

## Data Model

```go
type AuditLog struct {
    ID            uint          // 日志ID
    OperatorID    uint          // 操作人ID
    OperationType OperationType // 操作类型
    TargetType    TargetType    // 操作对象类型
    TargetID      uint          // 操作对象ID
    BeforeData    *string       // 操作前数据 (JSON)
    AfterData     *string       // 操作后数据 (JSON)
    CreatedAt     time.Time     // 操作时间
}
```

## Database Table

The audit logs are stored in the `audit_logs` table with the following structure:

- Primary key: `id`
- Indexes on: `operator_id`, `operation_type`, `(target_type, target_id)`, `created_at`
- Foreign key: `operator_id` references `users(id)`

## Error Handling

If audit log recording fails, the error is logged to stderr but does not affect the main business flow. This ensures that audit logging failures don't block critical operations.

## Testing

The service includes comprehensive unit tests using sqlmock. Run tests with:

```bash
go test ./internal/audit/...
```

## Implementation Notes

1. **JSON Serialization**: Before and after data are automatically serialized to JSON format
2. **Optional Data**: Both beforeData and afterData are optional (can be nil)
3. **Non-blocking**: Audit log failures are logged but don't return errors to prevent blocking business operations
4. **Immutable**: Audit logs are write-only and should never be modified or deleted
