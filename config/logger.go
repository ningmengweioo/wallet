package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// 日志级别
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
)

// 日志级别映射到整数，便于比较
var levelMap = map[string]int{
	LevelDebug: 0,
	LevelInfo:  1,
	LevelWarn:  2,
	LevelError: 3,
	LevelFatal: 4,
}

// Logger 日志管理器
type Logger struct {
	level  int
	output io.Writer
	sync.Mutex
	serviceName string
}

var (
	// 全局日志实例
	globalLogger *Logger
	once         sync.Once
)

// InitLogger 初始化日志管理器
func InitLogger() {
	once.Do(func() {
		// 从config包获取配置
		conf := GetConf()
		level := conf.Log.Level
		if level == "" {
			level = LevelInfo // 默认info级别
		}

		// 创建日志目录
		logDir := "log"
		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			os.Mkdir(logDir, 0755)
		}

		// 创建日志文件
		fileName := filepath.Join(logDir, fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02")))
		logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Failed to open log file: %v\n", err)
			// 如果打开文件失败，使用标准输出
			globalLogger = &Logger{
				level:       levelMap[level],
				output:      os.Stdout,
				serviceName: "blog-service",
			}
			return
		}

		// 创建一个io.MultiWriter，同时输出到控制台和文件
		multiWriter := io.MultiWriter(os.Stdout, logFile)

		globalLogger = &Logger{
			level:       levelMap[level],
			output:      multiWriter,
			serviceName: "blog-service",
		}

		// 记录初始化日志
		globalLogger.Info("Logger initialized successfully", "level", level)
	})
}

// GetLogger 获取全局日志实例
func GetLogger() *Logger {
	if globalLogger == nil {
		InitLogger()
	}
	return globalLogger
}

// 设置日志级别
func (l *Logger) SetLevel(level string) {
	l.Lock()
	defer l.Unlock()
	if lvl, ok := levelMap[level]; ok {
		l.level = lvl
	}
}

// Debug 调试级别日志
func (l *Logger) Debug(msg string, fields ...interface{}) {
	if l.level <= levelMap[LevelDebug] {
		writeLog(l, LevelDebug, msg, fields)
	}
}

// Info 信息级别日志
func (l *Logger) Info(msg string, fields ...interface{}) {
	if l.level <= levelMap[LevelInfo] {
		writeLog(l, LevelInfo, msg, fields)
	}
}

// Warn 警告级别日志
func (l *Logger) Warn(msg string, fields ...interface{}) {
	if l.level <= levelMap[LevelWarn] {
		writeLog(l, LevelWarn, msg, fields)
	}
}

// Error 错误级别日志
func (l *Logger) Error(msg string, fields ...interface{}) {
	if l.level <= levelMap[LevelError] {
		writeLog(l, LevelError, msg, fields)
	}
}

// Fatal 致命错误级别日志，记录后程序退出
func (l *Logger) Fatal(msg string, fields ...interface{}) {
	if l.level <= levelMap[LevelFatal] {
		writeLog(l, LevelFatal, msg, fields)
		os.Exit(1)
	}
}

// writeLog 写入日志
func writeLog(l *Logger, level string, msg string, fields []interface{}) {
	l.Lock()
	defer l.Unlock()

	// 获取调用者信息
	_, file, line, ok := runtime.Caller(3)
	caller := "unknown"
	if ok {
		// 获取文件名，不包含路径
		fileName := filepath.Base(file)
		caller = fmt.Sprintf("%s:%d", fileName, line)
	}

	// 格式化时间
	now := time.Now().Format("2006-01-02 15:04:05.000")

	// 构建基本日志信息
	logMsg := map[string]interface{}{
		"time":    now,
		"level":   level,
		"service": l.serviceName,
		"caller":  caller,
		"message": msg,
	}

	// 处理额外字段
	if len(fields) > 0 {
		// 确保字段数量为偶数（key-value对）
		if len(fields)%2 != 0 {
			fields = append(fields, "") // 添加一个空值，避免索引越界
		}

		for i := 0; i < len(fields); i += 2 {
			key, ok := fields[i].(string)
			if !ok {
				key = fmt.Sprintf("field_%d", i/2)
			}
			logMsg[key] = fields[i+1]
		}
	}

	// 转换为JSON
	jsonBytes, err := json.Marshal(logMsg)
	if err != nil {
		// 如果JSON转换失败，使用简单格式
		fmt.Fprintf(l.output, "%s [%s] %s - %s\n", now, level, caller, msg)
		return
	}

	// 写入日志
	fmt.Fprintf(l.output, "%s\n", string(jsonBytes))
}

// Debug 全局调试日志
func Debug(msg string, fields ...interface{}) {
	GetLogger().Debug(msg, fields...)
}

// Info 全局信息日志
func Info(msg string, fields ...interface{}) {
	GetLogger().Info(msg, fields...)
}

// Warn 全局警告日志
func Warn(msg string, fields ...interface{}) {
	GetLogger().Warn(msg, fields...)
}

// Error 全局错误日志
func Error(msg string, fields ...interface{}) {
	GetLogger().Error(msg, fields...)
}

// Fatal 全局致命错误日志
func Fatal(msg string, fields ...interface{}) {
	GetLogger().Fatal(msg, fields...)
}

// WithTrace 记录带有请求追踪信息的日志
func WithTrace(traceID string, level string, msg string, fields ...interface{}) {
	newFields := append([]interface{}{"trace_id", traceID}, fields...)

	logger := GetLogger()
	switch level {
	case LevelDebug:
		logger.Debug(msg, newFields...)
	case LevelInfo:
		logger.Info(msg, newFields...)
	case LevelWarn:
		logger.Warn(msg, newFields...)
	case LevelError:
		logger.Error(msg, newFields...)
	case LevelFatal:
		logger.Fatal(msg, newFields...)
	default:
		logger.Info(msg, newFields...)
	}
}
