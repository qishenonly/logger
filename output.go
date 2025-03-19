package logger

// OutputType 定义日志输出类型
type OutputType string

const (
	// OutputFile 仅输出到文件
	OutputFile OutputType = "file"
	// OutputTerminal 仅输出到终端
	OutputTerminal OutputType = "terminal"
	// OutputBoth 同时输出到文件和终端
	OutputBoth OutputType = "both"
)

// IsValidOutputType 检查输出类型是否有效
func IsValidOutputType(output OutputType) bool {
	return output == OutputFile || output == OutputTerminal || output == OutputBoth
}

// GetDefaultOutputType 获取默认输出类型
func GetDefaultOutputType() OutputType {
	return OutputTerminal // 默认输出到终端
}
