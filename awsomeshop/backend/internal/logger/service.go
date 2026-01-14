package logger

import (
	"fmt"

	"gorm.io/gorm"
)

// Service 应用日志服务
type Service struct {
	db       *gorm.DB
	minLevel LogLevel // 最小日志级别，低于此级别的日志不会被记录
}

// NewService 创建应用日志服务实例
// 默认最小日志级别为 INFO（记录所有日志）
func NewService(db *gorm.DB) *Service {
	return &Service{
		db:       db,
		minLevel: LogLevelInfo,
	}
}

// NewServiceWithLevel 创建应用日志服务实例并指定最小日志级别
func NewServiceWithLevel(db *gorm.DB, minLevel LogLevel) *Service {
	return &Service{
		db:       db,
		minLevel: minLevel,
	}
}

// SetMinLevel 设置最小日志级别
func (s *Service) SetMinLevel(level LogLevel) {
	s.minLevel = level
}

// GetMinLevel 获取当前最小日志级别
func (s *Service) GetMinLevel() LogLevel {
	return s.minLevel
}

// shouldLog 判断是否应该记录该级别的日志
func (s *Service) shouldLog(level LogLevel) bool {
	// 日志级别优先级：ERROR > WARNING > INFO
	levelPriority := map[LogLevel]int{
		LogLevelInfo:    1,
		LogLevelWarning: 2,
		LogLevelError:   3,
	}

	currentPriority, ok := levelPriority[level]
	if !ok {
		// 未知级别，默认记录
		return true
	}

	minPriority, ok := levelPriority[s.minLevel]
	if !ok {
		// 未知最小级别，默认记录
		return true
	}

	return currentPriority >= minPriority
}

// Info 记录信息级别日志
func (s *Service) Info(source string, message string, userID *uint) error {
	return s.log(LogLevelInfo, source, message, userID)
}

// Warning 记录警告级别日志
func (s *Service) Warning(source string, message string, userID *uint) error {
	return s.log(LogLevelWarning, source, message, userID)
}

// Error 记录错误级别日志
func (s *Service) Error(source string, message string, userID *uint) error {
	return s.log(LogLevelError, source, message, userID)
}

// log 内部方法：记录日志到数据库
func (s *Service) log(level LogLevel, source string, message string, userID *uint) error {
	// 检查是否应该记录该级别的日志
	if !s.shouldLog(level) {
		// 日志级别低于最小级别，跳过记录
		return nil
	}

	appLog := &AppLog{
		Level:   level,
		Message: message,
		Source:  source,
		UserID:  userID,
	}

	if err := s.db.Create(appLog).Error; err != nil {
		// 如果日志记录失败，打印到标准错误输出，但不返回错误
		// 避免日志记录失败影响主业务流程
		fmt.Printf("Failed to write log to database: %v\n", err)
		return err
	}

	return nil
}

// InfoWithUser 记录带用户ID的信息日志
func (s *Service) InfoWithUser(source string, message string, userID uint) error {
	return s.Info(source, message, &userID)
}

// WarningWithUser 记录带用户ID的警告日志
func (s *Service) WarningWithUser(source string, message string, userID uint) error {
	return s.Warning(source, message, &userID)
}

// ErrorWithUser 记录带用户ID的错误日志
func (s *Service) ErrorWithUser(source string, message string, userID uint) error {
	return s.Error(source, message, &userID)
}

// Infof 格式化记录信息日志
func (s *Service) Infof(source string, format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return s.Info(source, message, nil)
}

// Warningf 格式化记录警告日志
func (s *Service) Warningf(source string, format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return s.Warning(source, message, nil)
}

// Errorf 格式化记录错误日志
func (s *Service) Errorf(source string, format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return s.Error(source, message, nil)
}
