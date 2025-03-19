# 日志工具包

本日志工具包基于 [go.uber.org/zap](https://github.com/uber-go/zap) 实现，提供了灵活、可扩展的日志记录系统，支持多种输出适配器。

## 主要特点

1. 支持多级别日志：debug、info、warn、error、panic、fatal
2. 支持多种输出适配器：控制台、文件、Kafka 和 Elasticsearch 等
3. 支持结构化日志输出，便于与各种日志分析工具集成
4. 文件日志按日自动切换，按月归档
5. 分布式系统支持（节点 ID、模块名称和 IP 字段）
6. 支持错误详情记录，包括堆栈跟踪
7. 灵活的配置方式（结构体配置和函数选项模式）

## 日志存储结构

日志文件会按照以下结构自动组织：

```
logs/
  ├── 2023-01/
  │     ├── 01-01.log
  │     ├── 01-02.log
  │     └── ...
  ├── 2023-02/
  │     ├── 02-01.log
  │     ├── 02-02.log
  │     └── ...
  └── ...
```

- 顶级目录是配置的日志路径（例如 `./logs`）
- 第二级是按月归档的目录，格式为 `YYYY-MM`
- 日志文件按天生成，格式为 `MM-DD.log`

所有模块和节点的日志都会统一存放在这个目录结构中，通过日志内容中的`nodeid`和`module`字段区分。

## 快速开始

### 1. 初始化日志系统

#### 使用结构体配置

```go
import "github.com/qishenonly/logger"

func main() {
    // 初始化日志，设置日志级别和输出路径
    err := logger.Init(logger.Config{
        Level:      "info",              // 日志级别: debug, info, warn, error, panic
        Path:       "./logs",            // 日志文件路径，为空则只输出到控制台
        NodeID:     "node-001",          // 节点ID，用于分布式系统标识
        Module:     "api-server",        // 模块名称，如poc、finger等
        IP:         "192.168.1.10",      // IP地址
        OutputType: logger.OutputBoth,  // 输出类型: file(仅文件)、terminal(仅终端)、both(两者)
    })
    if err != nil {
        panic(err)
    }
    
    // 使用全局日志函数
    logger.Info("应用启动成功")
    logger.Infof("服务监听端口: %d", 8080)
}
```

#### 使用函数选项模式

```go
import "github.com/qishenonly/logger"

func main() {
    // 使用函数选项模式初始化日志
    err := logger.InitWithOptions(
        logger.WithLevel("info"),
        logger.WithPath("./logs"),
        logger.WithNodeID("node-001"),
        logger.WithModule("api-server"),
        logger.WithIP("192.168.1.10"),
        logger.WithBothOutput(),  // 同时输出到文件和终端
        // 或者使用以下选项之一:
        // logger.WithFileOutput(),     // 仅输出到文件
        // logger.WithTerminalOutput(), // 仅输出到终端
    )
    if err != nil {
        panic(err)
    }
    
    // 使用全局日志函数
    logger.Info("应用启动成功")
    logger.Infof("服务监听端口: %d", 8080)
}
```

### 2. 记录日志

#### 基本用法

```go
ctx := context.Background()

// 不同级别的日志
logger.Debug(ctx, "这是一条调试日志")
logger.Info(ctx, "这是一条信息日志")
logger.Warn(ctx, "这是一条警告日志")
logger.Error(ctx, "这是一条错误日志", logger.ErrorInfo{Err: err})
logger.Fatal(ctx, "这是一条致命错误日志", logger.ErrorInfo{Err: err}) // 会导致程序退出

// 格式化日志函数
logger.Debugf("调试信息: %v", data)
logger.Infof("用户 %s 登录成功，IP: %s", username, ip)
logger.Warnf("资源使用率高: %d%%", usage)
logger.Errorf("操作失败: %v", err)
logger.Panicf("严重错误: %v", err)
```

### 3. 创建独立的日志实例

#### 使用结构体配置

```go
// 创建独立的日志实例（使用相同的日志目录）
myLogger, err := logger.New(logger.Config{
    Level:      "debug",
    Path:       "./logs",         // 与全局相同的日志目录
    NodeID:     "node-002",
    Module:     "finger",
    IP:         "192.168.1.11",
    OutputType: logger.OutputFile, // 仅输出到文件
})
if err != nil {
    panic(err)
}

// 使用独立的日志实例
myLogger.Debug("这是调试信息")
myLogger.Infof("用户 %s 登录成功", username)
```

#### 使用函数选项模式

```go
// 创建独立的日志实例（使用函数选项模式）
myLogger, err := logger.NewWithOptions(
    logger.WithLevel("debug"),
    logger.WithPath("./logs"),
    logger.WithNodeID("node-002"),
    logger.WithModule("finger"),
    logger.WithIP("192.168.1.11"),
    logger.WithTerminalOutput(), // 仅输出到终端
)
if err != nil {
    panic(err)
}

// 使用独立的日志实例
myLogger.Debug("这是调试信息")
myLogger.Infof("用户 %s 登录成功", username)
```

## 日志级别

支持以下日志级别（按严重程度递增排序）:

- `debug`: 调试信息
- `info`: 一般信息
- `warn`: 警告信息
- `error`: 错误信息
- `panic`: 导致程序崩溃的严重错误
- `fatal`: 致命错误，记录后会导致程序退出

只有大于或等于配置级别的日志才会被输出。例如，如果配置级别为 `info`，则 `debug` 级别的日志不会输出。

## 日志输出类型

支持以下输出类型：

- `OutputFile`: 仅输出到文件
- `OutputTerminal`: 仅输出到终端
- `OutputBoth`: 同时输出到文件和终端（默认）

选择合适的输出类型可以满足不同场景的需求：
- 开发环境可能希望只输出到终端方便调试
- 生产环境可能需要同时输出到文件和终端
- 某些模块可能只需记录到文件而不需在终端显示

## 日志输出示例

### 控制台输出

控制台输出采用更加友好的格式:

```
2023-03-15T10:24:15.123+0800    INFO    应用启动成功    {"nodeid": "node-001", "module": "api-server", "ip": "192.168.1.10"}
2023-03-15T10:24:15.124+0800    INFO    服务监听端口: 8080    {"nodeid": "node-001", "module": "api-server", "ip": "192.168.1.10"}
2023-03-15T10:24:20.456+0800    ERROR   数据库连接失败    {"error": "connection refused", "nodeid": "node-001", "module": "api-server", "ip": "192.168.1.10"}
```

### 文件输出

文件输出采用JSON格式，便于解析和分析:

```json
{"level":"INFO","time":"2023-03-15T10:24:15.123+0800","caller":"app/main.go:42","msg":"应用启动成功","nodeid":"node-001","module":"api-server","ip":"192.168.1.10"}
{"level":"INFO","time":"2023-03-15T10:24:15.124+0800","caller":"app/main.go:43","msg":"服务监听端口: 8080","nodeid":"node-001","module":"api-server","ip":"192.168.1.10"}
{"level":"ERROR","time":"2023-03-15T10:24:20.456+0800","caller":"app/db.go:28","msg":"数据库连接失败","nodeid":"node-001","module":"api-server","ip":"192.168.1.10"}
```

## 函数选项模式

日志工具包支持函数选项模式，提供以下选项函数：

- `WithLevel(level string)`: 设置日志级别
- `WithPath(path string)`: 设置日志文件路径
- `WithNodeID(nodeID string)`: 设置节点ID
- `WithModule(module string)`: 设置模块名称
- `WithIP(ip string)`: 设置IP地址
- `WithOutputType(outputType OutputType)`: 设置输出类型
- `WithFileOutput()`: 设置仅输出到文件
- `WithTerminalOutput()`: 设置仅输出到终端
- `WithBothOutput()`: 设置同时输出到文件和终端

使用函数选项模式可以更灵活地配置日志，不需要每次都创建完整的Config结构体。

## 自定义适配器

您可以通过实现`LogAdapter`接口来创建自定义适配器：

```go
type LogAdapter interface {
    // Name 返回适配器名称
    Name() string
    
    // Init 初始化适配器
    Init(config map[string]interface{}) error
    
    // Process 处理单条日志
    Process(ctx context.Context, entry LogEntry) error
    
    // Flush 刷新缓冲区
    Flush() error
    
    // Close 关闭适配器
    Close() error
}
```

然后通过`RegisterAdapter`函数注册适配器：

```go
func init() {
    logger.RegisterAdapter("custom", func() logger.LogAdapter {
        return &CustomAdapter{}
    })
}
```

## 最佳实践

1. **合理设置日志级别**：生产环境通常使用info或warn级别，开发环境可使用debug级别。

2. **使用结构化字段**：尽量使用结构化字段而不是将所有信息塞入消息字符串，这样有利于后期分析。

3. **请求上下文传递**：确保在整个请求链路中传递上下文，以便跟踪完整的请求流程。

4. **记录关键操作**：记录系统的关键操作和状态变化，但避免记录敏感信息如密码、令牌等。

5. **错误日志详情**：记录错误时提供足够的上下文信息，包括操作名称、请求数据等。

6. **性能考虑**：对于高吞吐量服务，考虑使用批处理适配器并适当调整批处理参数。

7. **分布式环境**：在分布式环境中，请确保每个节点设置正确的NodeID、Module和IP值。

8. **日志目录共享**：不同模块可以创建不同的日志实例，但它们共用同一个日志目录，通过日志内容中的字段区分。

9. **输出类型选择**：可以为不同模块设置不同的输出类型，如某些模块仅记录到文件，而其他模块同时输出到终端和文件。

## 常见问题

### Q: 如何在单元测试中使用日志？

A: 可以使用`logger.InitForTest()`创建一个专用于测试的日志实例，它默认仅输出到控制台且不会影响测试结果。

### Q: 如何控制日志输出格式？

A: 通过适配器的`format`配置控制，目前支持`text`和`json`两种格式。

### Q: 日志系统是否支持异步写入？

A: 是的，Kafka和Elasticsearch适配器支持批量异步处理，文件适配器则是直接写入，但底层使用了缓冲。

### Q: 日志文件按天自动切换，如何配置？

A: 这是默认行为，日志文件会按照`YYYY-MM/MM-DD.log`的格式自动组织，无需额外配置。

### Q: 如何同时使用多种适配器？

A: 在配置中的`adapters`数组中添加多个适配器配置即可，每个适配器可以有不同的级别和配置。