package logger

import (
	"fmt"
	"sync"
)

var (
	defaultLogger Logger
	loggerMu      sync.Mutex
)

// AdapterConfig 定义适配器配置
type AdapterConfig struct {
	Name   string                 // 适配器名称
	Config map[string]interface{} // 适配器配置
}

// Config 定义日志配置
type Config struct {
	Level      string          // 日志级别: debug, info, warn, error, panic
	Path       string          // 日志文件路径，为空则只输出到控制台
	NodeID     string          // 节点ID，用于分布式系统标识当前节点
	Module     string          // 模块名称，如poc、finger等
	IP         string          // IP地址
	OutputType OutputType      // 输出类型：file、terminal、both
	Adapters   []AdapterConfig // 日志适配器配置
}

// Init 初始化默认日志
func Init(config Config) error {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	// 如果没有指定输出类型，设置默认值
	if !IsValidOutputType(config.OutputType) {
		config.OutputType = GetDefaultOutputType()
	}

	logger, err := NewZapLogger(config.Level, config.Path, config.NodeID, config.Module, config.IP, config.OutputType, config.Adapters)
	if err != nil {
		return fmt.Errorf("init logger failed: %v", err)
	}

	defaultLogger = logger
	return nil
}

// InitWithOptions 使用选项模式初始化默认日志
func InitWithOptions(opts ...Option) error {
	config := NewConfig(opts...)
	return Init(config)
}

// Default 获取默认日志实例
func Default() Logger {
	if defaultLogger == nil {
		// 如果默认日志未初始化，创建一个只输出到控制台的默认日志
		loggerMu.Lock()
		if defaultLogger == nil {
			logger, err := NewZapLogger("info", "", "", "default", "", OutputTerminal, nil) // 默认日志只输出到终端
			if err != nil {
				// 在极端情况下，如果创建日志失败，使用一个空实现避免空指针
				defaultLogger = &emptyLogger{}
			} else {
				defaultLogger = logger
			}
		}
		loggerMu.Unlock()
	}
	return defaultLogger
}

// New 创建一个新的日志实例
func New(config Config) (Logger, error) {
	// 如果没有指定输出类型，设置默认值
	if !IsValidOutputType(config.OutputType) {
		config.OutputType = GetDefaultOutputType()
	}

	return NewZapLogger(config.Level, config.Path, config.NodeID, config.Module, config.IP, config.OutputType, config.Adapters)
}

// NewWithOptions 使用选项模式创建一个新的日志实例
func NewWithOptions(opts ...Option) (Logger, error) {
	config := NewConfig(opts...)
	return New(config)
}

// 以下是全局日志函数，使用默认日志实例
func Panic(args ...any) {
	Default().Panic(args...)
}

func Panicf(format string, args ...any) {
	Default().Panicf(format, args...)
}

func Error(args ...any) {
	Default().Error(args...)
}

func Errorf(format string, args ...any) {
	Default().Errorf(format, args...)
}

func Warn(args ...any) {
	Default().Warn(args...)
}

func Warnf(format string, args ...any) {
	Default().Warnf(format, args...)
}

func Info(args ...any) {
	Default().Info(args...)
}

func Infof(format string, args ...any) {
	Default().Infof(format, args...)
}

func Debug(args ...any) {
	Default().Debug(args...)
}

func Debugf(format string, args ...any) {
	Default().Debugf(format, args...)
}
