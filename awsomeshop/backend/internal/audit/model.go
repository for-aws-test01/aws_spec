package audit

import (
	"time"
)

// OperationType 操作类型
type OperationType string

const (
	// OperationCreateEmployee 创建员工
	OperationCreateEmployee OperationType = "create_employee"
	// OperationUpdateEmployee 修改员工
	OperationUpdateEmployee OperationType = "update_employee"
	// OperationResetPassword 重置密码
	OperationResetPassword OperationType = "reset_password"
	// OperationCreateProduct 创建产品
	OperationCreateProduct OperationType = "create_product"
	// OperationUpdateProduct 修改产品
	OperationUpdateProduct OperationType = "update_product"
	// OperationDeleteProduct 删除产品
	OperationDeleteProduct OperationType = "delete_product"
	// OperationGrantPoints 发放积分
	OperationGrantPoints OperationType = "grant_points"
	// OperationDeductPoints 扣除积分
	OperationDeductPoints OperationType = "deduct_points"
	// OperationApproveOrder 核销订单
	OperationApproveOrder OperationType = "approve_order"
	// OperationRejectOrder 拒绝订单
	OperationRejectOrder OperationType = "reject_order"
)

// TargetType 操作对象类型
type TargetType string

const (
	// TargetTypeUser 用户
	TargetTypeUser TargetType = "user"
	// TargetTypeProduct 产品
	TargetTypeProduct TargetType = "product"
	// TargetTypeOrder 订单
	TargetTypeOrder TargetType = "order"
	// TargetTypePoint 积分
	TargetTypePoint TargetType = "point"
)

// AuditLog 审计日志模型
type AuditLog struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	OperatorID    uint           `gorm:"type:bigint unsigned;not null;index" json:"operator_id"`
	OperationType OperationType  `gorm:"type:varchar(50);not null;index" json:"operation_type"`
	TargetType    TargetType     `gorm:"type:varchar(50);not null;index:idx_target" json:"target_type"`
	TargetID      uint           `gorm:"type:bigint unsigned;not null;index:idx_target" json:"target_id"`
	BeforeData    *string        `gorm:"type:json" json:"before_data,omitempty"`
	AfterData     *string        `gorm:"type:json" json:"after_data,omitempty"`
	CreatedAt     time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}
