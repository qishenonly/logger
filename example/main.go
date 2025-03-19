package main

import (
	"time"

	"github.com/qishenonly/logger"
)

func main() {
	// 使用传统方式初始化日志 - 同时输出到文件和终端
	err := logger.Init(logger.Config{
		Level:      "debug",           // 日志级别：debug, info, warn, error, panic
		Path:       "./logs",          // 日志文件路径
		NodeID:     "node-001",        // 节点ID
		Module:     "api-server",      // 模块名称
		IP:         "192.168.1.10",    // IP地址
		OutputType: logger.OutputBoth, // 同时输出到文件和终端
	})
	if err != nil {
		panic(err)
	}

	// 或者使用函数选项模式初始化日志
	/*
		err := logger.InitWithOptions(
			logger.WithLevel("debug"),
			logger.WithPath("./logs"),
			logger.WithNodeID("node-001"),
			logger.WithModule("api-server"),
			logger.WithIP("192.168.1.10"),
			logger.WithBothOutput(), // 同时输出到文件和终端
			// 或者使用以下选项之一：
			// logger.WithFileOutput(),     // 仅输出到文件
			// logger.WithTerminalOutput(), // 仅输出到终端
		)
		if err != nil {
			panic(err)
		}
	*/

	// 使用全局日志函数
	logger.Info("示例程序启动")
	logger.Infof("当前时间: %s", time.Now().Format("2006-01-02 15:04:05"))

	// 不同级别的日志
	logger.Debug("这是一条调试日志")
	logger.Info("这是一条信息日志")
	logger.Warn("这是一条警告日志")
	logger.Error("这是一条错误日志")
	// logger.Panic("这是一条panic日志") // 这会导致程序崩溃

	// 创建POC模块的日志实例 - 仅输出到文件
	pocLogger, err := logger.New(logger.Config{
		Level:      "info",
		Path:       "./logs", // 共用相同的日志目录
		NodeID:     "node-001",
		Module:     "poc", // 不同的模块名称
		IP:         "192.168.1.10",
		OutputType: logger.OutputFile, // 仅输出到文件
	})
	if err != nil {
		logger.Errorf("创建POC模块日志失败: %v", err)
		return
	}

	// 使用POC模块的日志
	pocLogger.Info("POC模块初始化成功")
	pocLogger.Debugf("这条调试日志不会显示，因为级别设置为info")

	// 使用函数选项模式创建指纹识别模块的日志实例 - 仅输出到终端
	fingerLogger, err := logger.NewWithOptions(
		logger.WithLevel("info"),
		logger.WithPath("./logs"),
		logger.WithNodeID("node-002"),
		logger.WithModule("finger"),
		logger.WithIP("192.168.1.11"),
		logger.WithTerminalOutput(), // 仅输出到终端
	)
	if err != nil {
		logger.Errorf("创建指纹识别模块日志失败: %v", err)
		return
	}

	// 使用指纹识别模块的日志
	fingerLogger.Info("指纹识别模块已启动")

	// 模拟一些错误情况
	err = simulateError()
	if err != nil {
		logger.Errorf("操作失败: %v", err)
	}

	logger.Info("示例程序结束")
}

func simulateError() error {
	// 模拟一个错误
	logger.Debug("正在执行可能失败的操作...")
	// 返回一个模拟错误
	return &customError{message: "模拟的错误情况"}
}

// 自定义错误类型
type customError struct {
	message string
}

func (e *customError) Error() string {
	return e.message
}
