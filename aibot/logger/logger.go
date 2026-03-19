// Package logger 提供企业微信智能机器人 SDK 的日志接口和实现
package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Level 定义日志级别
type Level int

const (
	// LevelDebug 调试级别
	LevelDebug Level = iota
	// LevelInfo 信息级别
	LevelInfo
	// LevelWarn 警告级别
	LevelWarn
	// LevelError 错误级别
	LevelError
)

// String 返回日志级别的字符串表示
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger 定义日志接口
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	SetLevel(level Level)
	GetLevel() Level
}

// DefaultLogger 默认日志实现
type DefaultLogger struct {
	level  Level
	logger *log.Logger
}

// NewDefaultLogger 创建新的默认日志器
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		level:  LevelInfo,
		logger: log.New(os.Stdout, "", 0),
	}
}

// SetLevel 设置日志级别
func (l *DefaultLogger) SetLevel(level Level) {
	l.level = level
}

// GetLevel 获取当前日志级别
func (l *DefaultLogger) GetLevel() Level {
	return l.level
}

// log 内部日志方法
func (l *DefaultLogger) log(level Level, msg string, args ...interface{}) {
	if level < l.level {
		return
	}

	// 格式化消息
	formattedMsg := msg
	if len(args) > 0 {
		formattedMsg = fmt.Sprintf(msg, args...)
	}

	// 添加时间戳和日志级别
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	logLine := fmt.Sprintf("%s [%s] %s", timestamp, level.String(), formattedMsg)

	// 输出日志
	l.logger.Println(logLine)
}

// Debug 输出调试日志
func (l *DefaultLogger) Debug(msg string, args ...interface{}) {
	l.log(LevelDebug, msg, args...)
}

// Info 输出信息日志
func (l *DefaultLogger) Info(msg string, args ...interface{}) {
	l.log(LevelInfo, msg, args...)
}

// Warn 输出警告日志
func (l *DefaultLogger) Warn(msg string, args ...interface{}) {
	l.log(LevelWarn, msg, args...)
}

// Error 输出错误日志
func (l *DefaultLogger) Error(msg string, args ...interface{}) {
	l.log(LevelError, msg, args...)
}

// NullLogger 空日志实现，不输出任何日志
type NullLogger struct{}

// Debug 不输出调试日志
func (l *NullLogger) Debug(msg string, args ...interface{}) {}

// Info 不输出信息日志
func (l *NullLogger) Info(msg string, args ...interface{}) {}

// Warn 不输出警告日志
func (l *NullLogger) Warn(msg string, args ...interface{}) {}

// Error 不输出错误日志
func (l *NullLogger) Error(msg string, args ...interface{}) {}

// SetLevel 设置日志级别（空实现）
func (l *NullLogger) SetLevel(level Level) {}

// GetLevel 获取日志级别（空实现）
func (l *NullLogger) GetLevel() Level { return LevelInfo }

// LevelFromString 从字符串解析日志级别
func LevelFromString(levelStr string) Level {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARN", "WARNING":
		return LevelWarn
	case "ERROR":
		return LevelError
	default:
		return LevelInfo
	}
}
