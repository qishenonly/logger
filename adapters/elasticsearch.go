package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/qishenonly/logger"
)

func init() {
	// 注册适配器
	logger.RegisterAdapter("elasticsearch", func() logger.LogAdapter {
		return &ElasticsearchAdapter{}
	})
}

// ElasticsearchAdapter 用于将日志输出到Elasticsearch
type ElasticsearchAdapter struct {
	Hosts         []string
	Index         string
	Username      string
	Password      string
	BulkSize      int
	FlushInterval time.Duration
	buffer        []logger.LogEntry
	bufferMu      sync.Mutex
	client        interface{} // 这里用interface{}占位，实际应该是ES客户端
}

// Name 返回适配器名称
func (a *ElasticsearchAdapter) Name() string {
	return "elasticsearch"
}

// Init 初始化适配器
func (a *ElasticsearchAdapter) Init(config map[string]interface{}) error {
	// 解析配置参数
	if hosts, ok := config["hosts"].([]interface{}); ok {
		a.Hosts = make([]string, 0, len(hosts))
		for _, host := range hosts {
			if h, ok := host.(string); ok {
				a.Hosts = append(a.Hosts, h)
			}
		}
	} else {
		a.Hosts = []string{"http://localhost:9200"}
	}

	if index, ok := config["index"].(string); ok {
		a.Index = index
	} else {
		a.Index = "logs-" + time.Now().Format("2006.01.02")
	}

	if username, ok := config["username"].(string); ok {
		a.Username = username
	}

	if password, ok := config["password"].(string); ok {
		a.Password = password
	}

	if bulkSize, ok := config["bulk_size"].(float64); ok {
		a.BulkSize = int(bulkSize)
	} else {
		a.BulkSize = 200
	}

	if flushInterval, ok := config["flush_interval"].(float64); ok {
		a.FlushInterval = time.Duration(flushInterval) * time.Second
	} else {
		a.FlushInterval = 10 * time.Second
	}

	// 初始化日志缓冲区
	a.buffer = make([]logger.LogEntry, 0, a.BulkSize)

	// 连接到Elasticsearch
	// 实际应该这样:
	// cfg := elasticsearch.Config{
	//     Addresses: a.Hosts,
	//     Username:  a.Username,
	//     Password:  a.Password,
	// }
	// client, err := elasticsearch.NewClient(cfg)
	// if err != nil {
	//     return fmt.Errorf("failed to create elasticsearch client: %v", err)
	// }
	// a.client = client

	// 定期刷新缓冲区
	go a.flushPeriodically()

	return nil
}

// Process 处理日志条目
func (a *ElasticsearchAdapter) Process(ctx context.Context, entry logger.LogEntry) error {
	a.bufferMu.Lock()
	defer a.bufferMu.Unlock()

	// 添加到缓冲区
	a.buffer = append(a.buffer, entry)

	// 如果达到批量大小，刷新缓冲区
	if len(a.buffer) >= a.BulkSize {
		return a.flushBuffer()
	}

	return nil
}

// Flush 刷新缓冲区
func (a *ElasticsearchAdapter) Flush() error {
	a.bufferMu.Lock()
	defer a.bufferMu.Unlock()

	return a.flushBuffer()
}

// flushBuffer 刷新缓冲区（无锁版本，调用前需要获取锁）
func (a *ElasticsearchAdapter) flushBuffer() error {
	if len(a.buffer) == 0 {
		return nil
	}

	// 在实际应用中，这里应该批量发送到Elasticsearch
	// var buf bytes.Buffer
	// for _, entry := range a.buffer {
	//     // 为每个文档添加索引元数据
	//     meta := []byte(fmt.Sprintf(`{ "index" : { "_index" : "%s" } }%s`, a.Index, "\n"))
	//     data, err := json.Marshal(entry)
	//     if err != nil {
	//         continue
	//     }
	//     data = append(data, '\n')
	//
	//     buf.Write(meta)
	//     buf.Write(data)
	// }
	//
	// res, err := a.client.(*elasticsearch.Client).Bulk(bytes.NewReader(buf.Bytes()))
	// if err != nil {
	//     return err
	// }
	// defer res.Body.Close()
	//
	// if res.IsError() {
	//     return fmt.Errorf("elasticsearch bulk request failed: %s", res.String())
	// }

	// 这里仅作演示，实际打印日志
	for _, entry := range a.buffer {
		data, _ := json.Marshal(entry)
		fmt.Printf("[Elasticsearch Adapter] Would index to %s: %s\n", a.Index, string(data))
	}

	// 清空缓冲区
	a.buffer = a.buffer[:0]

	return nil
}

// flushPeriodically 定期刷新缓冲区
func (a *ElasticsearchAdapter) flushPeriodically() {
	ticker := time.NewTicker(a.FlushInterval)
	defer ticker.Stop()

	for range ticker.C {
		_ = a.Flush()
	}
}

// Close 关闭适配器
func (a *ElasticsearchAdapter) Close() error {
	// 刷新剩余日志
	return a.Flush()
}
