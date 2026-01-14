package audit

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestNewService(t *testing.T) {
	db, _ := setupTestDB(t)
	service := NewService(db)
	assert.NotNil(t, service)
	assert.NotNil(t, service.db)
}

func TestService_Log(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewService(db)

	operatorID := uint(1)
	targetID := uint(100)
	beforeData := map[string]interface{}{"name": "old name"}
	afterData := map[string]interface{}{"name": "new name"}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `audit_logs`").
		WithArgs(operatorID, OperationUpdateEmployee, TargetTypeUser, targetID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.Log(operatorID, OperationUpdateEmployee, TargetTypeUser, targetID, beforeData, afterData)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_LogWithNilData(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewService(db)

	operatorID := uint(1)
	targetID := uint(100)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `audit_logs`").
		WithArgs(operatorID, OperationCreateEmployee, TargetTypeUser, targetID, nil, nil, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.Log(operatorID, OperationCreateEmployee, TargetTypeUser, targetID, nil, nil)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_LogCreateEmployee(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewService(db)

	operatorID := uint(1)
	targetID := uint(100)
	afterData := map[string]interface{}{
		"employee_id": "20260114-001",
		"name":        "张三",
		"role":        "employee",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `audit_logs`").
		WithArgs(operatorID, OperationCreateEmployee, TargetTypeUser, targetID, nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.LogCreateEmployee(operatorID, targetID, afterData)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_LogUpdateEmployee(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewService(db)

	operatorID := uint(1)
	targetID := uint(100)
	beforeData := map[string]interface{}{"position": "工程师"}
	afterData := map[string]interface{}{"position": "高级工程师"}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `audit_logs`").
		WithArgs(operatorID, OperationUpdateEmployee, TargetTypeUser, targetID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.LogUpdateEmployee(operatorID, targetID, beforeData, afterData)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_LogCreateProduct(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewService(db)

	operatorID := uint(1)
	targetID := uint(200)
	afterData := map[string]interface{}{
		"name":            "AirPods Pro",
		"points_required": 500,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `audit_logs`").
		WithArgs(operatorID, OperationCreateProduct, TargetTypeProduct, targetID, nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.LogCreateProduct(operatorID, targetID, afterData)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_LogUpdateProduct(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewService(db)

	operatorID := uint(1)
	targetID := uint(200)
	beforeData := map[string]interface{}{"points_required": 500}
	afterData := map[string]interface{}{"points_required": 600}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `audit_logs`").
		WithArgs(operatorID, OperationUpdateProduct, TargetTypeProduct, targetID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.LogUpdateProduct(operatorID, targetID, beforeData, afterData)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_LogDeleteProduct(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewService(db)

	operatorID := uint(1)
	targetID := uint(200)
	beforeData := map[string]interface{}{
		"name":            "AirPods Pro",
		"points_required": 500,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `audit_logs`").
		WithArgs(operatorID, OperationDeleteProduct, TargetTypeProduct, targetID, sqlmock.AnyArg(), nil, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.LogDeleteProduct(operatorID, targetID, beforeData)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_LogGrantPoints(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewService(db)

	operatorID := uint(1)
	targetID := uint(100)
	afterData := map[string]interface{}{
		"amount": 1000,
		"reason": "入职发放",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `audit_logs`").
		WithArgs(operatorID, OperationGrantPoints, TargetTypePoint, targetID, nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.LogGrantPoints(operatorID, targetID, afterData)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_LogDeductPoints(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewService(db)

	operatorID := uint(1)
	targetID := uint(100)
	afterData := map[string]interface{}{
		"amount": -500,
		"reason": "违规扣除",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `audit_logs`").
		WithArgs(operatorID, OperationDeductPoints, TargetTypePoint, targetID, nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.LogDeductPoints(operatorID, targetID, afterData)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_LogApproveOrder(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewService(db)

	operatorID := uint(1)
	targetID := uint(300)
	beforeData := map[string]interface{}{"status": "pending"}
	afterData := map[string]interface{}{"status": "approved"}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `audit_logs`").
		WithArgs(operatorID, OperationApproveOrder, TargetTypeOrder, targetID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.LogApproveOrder(operatorID, targetID, beforeData, afterData)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_LogRejectOrder(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewService(db)

	operatorID := uint(1)
	targetID := uint(300)
	beforeData := map[string]interface{}{"status": "pending"}
	afterData := map[string]interface{}{"status": "rejected", "reject_reason": "产品缺货"}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `audit_logs`").
		WithArgs(operatorID, OperationRejectOrder, TargetTypeOrder, targetID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.LogRejectOrder(operatorID, targetID, beforeData, afterData)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
