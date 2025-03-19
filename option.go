package logger

// Option 定义日志配置选项
type Option func(*Config)

// WithLevel 设置日志级别
func WithLevel(level string) Option {
	return func(c *Config) {
		c.Level = level
	}
}

// WithPath 设置日志文件路径
func WithPath(path string) Option {
	return func(c *Config) {
		c.Path = path
	}
}

// WithNodeID 设置节点ID
func WithNodeID(nodeID string) Option {
	return func(c *Config) {
		c.NodeID = nodeID
	}
}

// WithModule 设置模块名称
func WithModule(module string) Option {
	return func(c *Config) {
		c.Module = module
	}
}

// WithIP 设置IP地址
func WithIP(ip string) Option {
	return func(c *Config) {
		c.IP = ip
	}
}

// WithOutputType 设置输出类型
func WithOutputType(outputType OutputType) Option {
	return func(c *Config) {
		c.OutputType = outputType
	}
}

// WithFileOutput 设置输出仅到文件
func WithFileOutput() Option {
	return WithOutputType(OutputFile)
}

// WithTerminalOutput 设置输出仅到终端
func WithTerminalOutput() Option {
	return WithOutputType(OutputTerminal)
}

// WithBothOutput 设置同时输出到文件和终端
func WithBothOutput() Option {
	return WithOutputType(OutputBoth)
}

// WithAdapter 添加一个日志适配器
func WithAdapter(name string, config map[string]interface{}) Option {
	return func(c *Config) {
		c.Adapters = append(c.Adapters, AdapterConfig{
			Name:   name,
			Config: config,
		})
	}
}

// WithElasticsearchAdapter 添加Elasticsearch适配器
func WithElasticsearchAdapter(config map[string]interface{}) Option {
	return WithAdapter("elasticsearch", config)
}

// WithKafkaAdapter 添加Kafka适配器
func WithKafkaAdapter(config map[string]interface{}) Option {
	return WithAdapter("kafka", config)
}

// WithPrometheusAdapter 添加Prometheus适配器
func WithPrometheusAdapter(config map[string]interface{}) Option {
	return WithAdapter("prometheus", config)
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		Level:      "info",
		Path:       "",
		NodeID:     "",
		Module:     "default",
		IP:         "",
		OutputType: GetDefaultOutputType(),
		Adapters:   []AdapterConfig{},
	}
}

// NewConfig 通过选项创建新的配置
func NewConfig(opts ...Option) Config {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}
