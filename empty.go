package logger

// emptyLogger 是一个空的日志实现，不做任何操作
// 用于在极端情况下避免空指针异常
type emptyLogger struct{}

func (l *emptyLogger) Panic(args ...any) {}

func (l *emptyLogger) Panicf(format string, args ...any) {}

func (l *emptyLogger) Error(args ...any) {}

func (l *emptyLogger) Errorf(format string, args ...any) {}

func (l *emptyLogger) Warn(args ...any) {}

func (l *emptyLogger) Warnf(format string, args ...any) {}

func (l *emptyLogger) Info(args ...any) {}

func (l *emptyLogger) Infof(format string, args ...any) {}

func (l *emptyLogger) Debug(args ...any) {}

func (l *emptyLogger) Debugf(format string, args ...any) {}

func (l *emptyLogger) Close() error { return nil }

func (l *emptyLogger) AddAdapter(adapter LogAdapter) {}

func (l *emptyLogger) RemoveAdapter(name string) {}
