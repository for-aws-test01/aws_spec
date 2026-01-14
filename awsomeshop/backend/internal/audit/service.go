package audit

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

// Service 审计日志服务
type Service struct {
	db *gorm.DB
}

// NewService 创建审计日志服务实例
func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

// Log 记录审计日志
// operatorID: 操作人ID
// operationType: 操作类型
// targetType: 操作对象类型
// targetID: 操作对象ID
// beforeData: 操作前数据（可选）
// afterData: 操作后数据（可选）
func (s *Service) Log(
	operatorID uint,
	operationType OperationType,
	targetType TargetType,
	targetID uint,
	beforeData interface{},
	afterData interface{},
) error {
	auditLog := &AuditLog{
		OperatorID:    operatorID,
		OperationType: operationType,
		TargetType:    targetType,
		TargetID:      targetID,
	}

	// 序列化操作前数据
	if beforeData != nil {
		beforeJSON, err := json.Marshal(beforeData)
		if err != nil {
			return fmt.Errorf("failed to marshal before_data: %w", err)
		}
		beforeStr := string(beforeJSON)
		auditLog.BeforeData = &beforeStr
	}

	// 序列化操作后数据
	if afterData != nil {
		afterJSON, err := json.Marshal(afterData)
		if err != nil {
			return fmt.Errorf("failed to marshal after_data: %w", err)
		}
		afterStr := string(afterJSON)
		auditLog.AfterData = &afterStr
	}

	// 写入数据库
	if err := s.db.Create(auditLog).Error; err != nil {
		// 如果审计日志记录失败，打印到标准错误输出，但不返回错误
		// 避免审计日志记录失败影响主业务流程
		fmt.Printf("Failed to write audit log to database: %v\n", err)
		return err
	}

	return nil
}

// LogCreateEmployee 记录创建员工操作
func (s *Service) LogCreateEmployee(operatorID uint, targetID uint, afterData interface{}) error {
	return s.Log(operatorID, OperationCreateEmployee, TargetTypeUser, targetID, nil, afterData)
}

// LogUpdateEmployee 记录修改员工操作
func (s *Service) LogUpdateEmployee(operatorID uint, targetID uint, beforeData interface{}, afterData interface{}) error {
	return s.Log(operatorID, OperationUpdateEmployee, TargetTypeUser, targetID, beforeData, afterData)
}

// LogCreateProduct 记录创建产品操作
func (s *Service) LogCreateProduct(operatorID uint, targetID uint, afterData interface{}) error {
	return s.Log(operatorID, OperationCreateProduct, TargetTypeProduct, targetID, nil, afterData)
}

// LogUpdateProduct 记录修改产品操作
func (s *Service) LogUpdateProduct(operatorID uint, targetID uint, beforeData interface{}, afterData interface{}) error {
	return s.Log(operatorID, OperationUpdateProduct, TargetTypeProduct, targetID, beforeData, afterData)
}

// LogDeleteProduct 记录删除产品操作
func (s *Service) LogDeleteProduct(operatorID uint, targetID uint, beforeData interface{}) error {
	return s.Log(operatorID, OperationDeleteProduct, TargetTypeProduct, targetID, beforeData, nil)
}

// LogGrantPoints 记录发放积分操作
func (s *Service) LogGrantPoints(operatorID uint, targetID uint, afterData interface{}) error {
	return s.Log(operatorID, OperationGrantPoints, TargetTypePoint, targetID, nil, afterData)
}

// LogDeductPoints 记录扣除积分操作
func (s *Service) LogDeductPoints(operatorID uint, targetID uint, afterData interface{}) error {
	return s.Log(operatorID, OperationDeductPoints, TargetTypePoint, targetID, nil, afterData)
}

// LogApproveOrder 记录核销订单操作
func (s *Service) LogApproveOrder(operatorID uint, targetID uint, beforeData interface{}, afterData interface{}) error {
	return s.Log(operatorID, OperationApproveOrder, TargetTypeOrder, targetID, beforeData, afterData)
}

// LogRejectOrder 记录拒绝订单操作
func (s *Service) LogRejectOrder(operatorID uint, targetID uint, beforeData interface{}, afterData interface{}) error {
	return s.Log(operatorID, OperationRejectOrder, TargetTypeOrder, targetID, beforeData, afterData)
}
