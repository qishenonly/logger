package logger

type Logger interface {
	Panic(args ...any)
	Panicf(fmt string, args ...any)
	Error(args ...any)
	Errorf(fmt string, args ...any)
	Warn(args ...any)
	Warnf(fmt string, args ...any)
	Info(args ...any)
	Infof(fmt string, args ...any)
	Debug(args ...any)
	Debugf(fmt string, args ...any)

	// Close 关闭日志记录器
	Close() error

	// AddAdapter 添加一个适配器
	AddAdapter(adapter LogAdapter)

	// RemoveAdapter 移除一个适配器
	RemoveAdapter(name string)
}
