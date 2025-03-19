package logger

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger 实现Logger接口的zap日志处理器
type ZapLogger struct {
	logger    *zap.Logger
	sugar     *zap.SugaredLogger
	adapters  []LogAdapter
	adapterMu sync.RWMutex
	nodeID    string
	module    string
	ip        string
}

// NewZapLogger 创建一个新的zap日志处理器
func NewZapLogger(logLevel string, logPath string, nodeID string, module string, ip string, outputType OutputType, adapterConfigs []AdapterConfig) (Logger, error) {
	// 解析日志级别
	level := zap.InfoLevel
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "panic":
		level = zap.PanicLevel
	default:
		level = zap.InfoLevel
	}

	// 创建核心编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建多核心日志写入
	cores := []zapcore.Core{}

	// 根据输出类型选择输出目标
	if outputType == OutputTerminal || outputType == OutputBoth {
		// 控制台输出
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			zap.NewAtomicLevelAt(level),
		)
		cores = append(cores, consoleCore)
	}

	// 文件输出（按天）
	if (outputType == OutputFile || outputType == OutputBoth) && logPath != "" {
		// 使用日志旋转器
		rotator, err := NewDailyRotateWriter(logPath)
		if err != nil {
			return nil, fmt.Errorf("create log rotator failed: %v", err)
		}

		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		fileCore := zapcore.NewCore(
			fileEncoder,
			rotator.AsWriteSyncer(),
			zap.NewAtomicLevelAt(level),
		)
		cores = append(cores, fileCore)
	}

	// 如果没有任何有效的输出核心，至少添加一个控制台输出
	if len(cores) == 0 {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			zap.NewAtomicLevelAt(level),
		)
		cores = append(cores, consoleCore)
	}

	// 合并所有核心
	core := zapcore.NewTee(cores...)

	// 添加公共字段
	fields := []zap.Field{}
	if nodeID != "" {
		fields = append(fields, zap.String("nodeId", nodeID))
	}
	if module != "" {
		fields = append(fields, zap.String("module", module))
	}
	if ip != "" {
		fields = append(fields, zap.String("ip", ip))
	}

	// 创建logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.Fields(fields...))

	// 初始化适配器
	adapters := make([]LogAdapter, 0, len(adapterConfigs))
	if adapterConfigs != nil {
		for _, cfg := range adapterConfigs {
			adapter, exists := GetAdapter(cfg.Name)
			if !exists {
				continue
			}

			if err := adapter.Init(cfg.Config); err != nil {
				return nil, fmt.Errorf("init adapter %s failed: %v", cfg.Name, err)
			}

			adapters = append(adapters, adapter)
		}
	}

	return &ZapLogger{
		logger:   logger,
		sugar:    logger.Sugar(),
		adapters: adapters,
		nodeID:   nodeID,
		module:   module,
		ip:       ip,
	}, nil
}

// sendToAdapters 将日志发送到所有适配器
func (l *ZapLogger) sendToAdapters(level string, message string, properties map[string]interface{}) {
	l.adapterMu.RLock()
	defer l.adapterMu.RUnlock()

	if len(l.adapters) == 0 {
		return
	}

	// 创建日志条目
	entry := LogEntry{
		Level:      level,
		Time:       time.Now(),
		Message:    message,
		NodeID:     l.nodeID,
		Module:     l.module,
		IP:         l.ip,
		Properties: properties,
	}

	// 异步发送到适配器
	for _, adapter := range l.adapters {
		go func(a LogAdapter, e LogEntry) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = a.Process(ctx, e)
		}(adapter, entry)
	}
}

// Close 关闭日志记录器及其适配器
func (l *ZapLogger) Close() error {
	l.adapterMu.Lock()
	defer l.adapterMu.Unlock()

	for _, adapter := range l.adapters {
		_ = adapter.Flush()
		_ = adapter.Close()
	}

	l.adapters = nil
	return nil
}

// AddAdapter 添加一个适配器
func (l *ZapLogger) AddAdapter(adapter LogAdapter) {
	l.adapterMu.Lock()
	defer l.adapterMu.Unlock()
	l.adapters = append(l.adapters, adapter)
}

// RemoveAdapter 移除一个适配器
func (l *ZapLogger) RemoveAdapter(name string) {
	l.adapterMu.Lock()
	defer l.adapterMu.Unlock()

	for i, adapter := range l.adapters {
		if adapter.Name() == name {
			l.adapters = append(l.adapters[:i], l.adapters[i+1:]...)
			return
		}
	}
}

// 实现Logger接口方法

func (l *ZapLogger) Panic(args ...any) {
	l.sendToAdapters("panic", fmt.Sprint(args...), nil)
	l.sugar.Panic(args...)
}

func (l *ZapLogger) Panicf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.sendToAdapters("panic", msg, nil)
	l.sugar.Panicf(format, args...)
}

func (l *ZapLogger) Error(args ...any) {
	l.sendToAdapters("error", fmt.Sprint(args...), nil)
	l.sugar.Error(args...)
}

func (l *ZapLogger) Errorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.sendToAdapters("error", msg, nil)
	l.sugar.Errorf(format, args...)
}

func (l *ZapLogger) Warn(args ...any) {
	l.sendToAdapters("warn", fmt.Sprint(args...), nil)
	l.sugar.Warn(args...)
}

func (l *ZapLogger) Warnf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.sendToAdapters("warn", msg, nil)
	l.sugar.Warnf(format, args...)
}

func (l *ZapLogger) Info(args ...any) {
	l.sendToAdapters("info", fmt.Sprint(args...), nil)
	l.sugar.Info(args...)
}

func (l *ZapLogger) Infof(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.sendToAdapters("info", msg, nil)
	l.sugar.Infof(format, args...)
}

func (l *ZapLogger) Debug(args ...any) {
	l.sendToAdapters("debug", fmt.Sprint(args...), nil)
	l.sugar.Debug(args...)
}

func (l *ZapLogger) Debugf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.sendToAdapters("debug", msg, nil)
	l.sugar.Debugf(format, args...)
}
