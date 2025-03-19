package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap/zapcore"
)

// DailyRotateWriter 按天旋转的日志写入器，并按月归档
type DailyRotateWriter struct {
	logPath     string
	file        *os.File
	currentDate string
	mutex       sync.Mutex
}

// NewDailyRotateWriter 创建一个按天旋转、按月归档的日志写入器
func NewDailyRotateWriter(logPath string) (*DailyRotateWriter, error) {
	writer := &DailyRotateWriter{
		logPath: logPath,
	}

	if err := writer.rotateFile(); err != nil {
		return nil, err
	}

	return writer, nil
}

// Write 实现io.Writer接口
func (w *DailyRotateWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	today := time.Now().Format("2006-01-02")
	if today != w.currentDate {
		if err := w.rotateFile(); err != nil {
			return 0, err
		}
	}

	return w.file.Write(p)
}

// Sync 实现zapcore.WriteSyncer接口
func (w *DailyRotateWriter) Sync() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}

// Close 关闭文件
func (w *DailyRotateWriter) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.file != nil {
		err := w.file.Close()
		w.file = nil
		return err
	}
	return nil
}

// rotateFile 旋转日志文件
func (w *DailyRotateWriter) rotateFile() error {
	// 关闭旧文件
	if w.file != nil {
		err := w.file.Close()
		if err != nil {
			return err
		}
		w.file = nil
	}

	// 更新当前日期
	now := time.Now()
	w.currentDate = now.Format("2006-01-02")

	// 按月归档的目录结构: logs/2006-01/02.log
	monthDir := now.Format("2006-01")
	dayFile := now.Format("01-02.log")

	// 月度目录路径
	monthDirPath := filepath.Join(w.logPath, monthDir)

	// 确保月度目录存在
	if err := os.MkdirAll(monthDirPath, 0755); err != nil {
		return fmt.Errorf("create month directory failed: %v", err)
	}

	// 创建新的日志文件
	logFileName := filepath.Join(monthDirPath, dayFile)
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open log file failed: %v", err)
	}

	w.file = file
	return nil
}

// AsWriteSyncer 将DailyRotateWriter转换为zapcore.WriteSyncer
func (w *DailyRotateWriter) AsWriteSyncer() zapcore.WriteSyncer {
	return zapcore.AddSync(w)
}
