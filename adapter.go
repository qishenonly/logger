package logger

import (
	"context"
	"time"
)

// LogEntry 表示一条完整的日志记录
type LogEntry struct {
	Level      string                 // 日志级别
	Time       time.Time              // 日志时间
	Message    string                 // 日志消息
	Caller     string                 // 调用位置
	NodeID     string                 // 节点ID
	Module     string                 // 模块名称
	IP         string                 // IP地址
	Properties map[string]interface{} // 额外属性
}

// LogAdapter 日志适配器接口，第三方组件可以实现这个接口接收日志
type LogAdapter interface {
	// Name 返回适配器名称
	Name() string

	// Init 初始化适配器
	Init(config map[string]interface{}) error

	// Process 处理一条日志记录
	Process(ctx context.Context, entry LogEntry) error

	// Flush 刷新缓存的日志
	Flush() error

	// Close 关闭适配器
	Close() error
}

// LogAdapterCreator 适配器创建函数类型
type LogAdapterCreator func() LogAdapter

var (
	// adapterRegistry 存储已注册的适配器创建函数
	adapterRegistry = make(map[string]LogAdapterCreator)
)

// RegisterAdapter 注册一个日志适配器
func RegisterAdapter(name string, creator LogAdapterCreator) {
	adapterRegistry[name] = creator
}

// GetAdapter 获取已注册的日志适配器
func GetAdapter(name string) (LogAdapter, bool) {
	creator, exists := adapterRegistry[name]
	if !exists {
		return nil, false
	}
	return creator(), true
}
