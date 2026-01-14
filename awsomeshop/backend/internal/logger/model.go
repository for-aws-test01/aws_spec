package logger

import (
	"time"
)

// LogLevel 日志级别
type LogLevel string

const (
	// LogLevelInfo 信息级别
	LogLevelInfo LogLevel = "INFO"
	// LogLevelWarning 警告级别
	LogLevelWarning LogLevel = "WARNING"
	// LogLevelError 错误级别
	LogLevelError LogLevel = "ERROR"
)

// AppLog 应用日志模型
type AppLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Level     LogLevel  `gorm:"type:enum('INFO','WARNING','ERROR');not null;default:'INFO'" json:"level"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Source    string    `gorm:"type:varchar(100);not null" json:"source"`
	UserID    *uint     `gorm:"type:bigint unsigned;index" json:"user_id,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (AppLog) TableName() string {
	return "app_logs"
}
