package logger

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

func TestService_Info(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewService(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `app_logs`").
		WithArgs(LogLevelInfo, "test message", "auth", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.Info("auth", "test message", nil)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_Warning(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewService(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `app_logs`").
		WithArgs(LogLevelWarning, "warning message", "product", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.Warning("product", "warning message", nil)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewService(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `app_logs`").
		WithArgs(LogLevelError, "error message", "order", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.Error("order", "error message", nil)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_InfoWithUser(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewService(db)

	userID := uint(123)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `app_logs`").
		WithArgs(LogLevelInfo, "user action", "auth", userID, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.InfoWithUser("auth", "user action", userID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_Infof(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewService(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `app_logs`").
		WithArgs(LogLevelInfo, "User 123 logged in", "auth", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.Infof("auth", "User %d logged in", 123)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_Warningf(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewService(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `app_logs`").
		WithArgs(LogLevelWarning, "Product 456 stock low", "product", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.Warningf("product", "Product %d stock low", 456)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_Errorf(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewService(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `app_logs`").
		WithArgs(LogLevelError, "Failed to process order 789", "order", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.Errorf("order", "Failed to process order %d", 789)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
