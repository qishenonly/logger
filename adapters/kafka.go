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
	logger.RegisterAdapter("kafka", func() logger.LogAdapter {
		return &KafkaAdapter{}
	})
}

// KafkaAdapter 用于将日志输出到Kafka
type KafkaAdapter struct {
	Brokers      []string
	Topic        string
	BatchSize    int
	FlushTimeout time.Duration
	buffer       []logger.LogEntry
	bufferMu     sync.Mutex
	producer     interface{} // 这里用interface{}占位，实际应该是Kafka生产者
}

// Name 返回适配器名称
func (a *KafkaAdapter) Name() string {
	return "kafka"
}

// Init 初始化适配器
func (a *KafkaAdapter) Init(config map[string]interface{}) error {
	// 解析配置参数
	if brokers, ok := config["brokers"].([]interface{}); ok {
		a.Brokers = make([]string, 0, len(brokers))
		for _, broker := range brokers {
			if b, ok := broker.(string); ok {
				a.Brokers = append(a.Brokers, b)
			}
		}
	} else {
		a.Brokers = []string{"localhost:9092"}
	}

	if topic, ok := config["topic"].(string); ok {
		a.Topic = topic
	} else {
		a.Topic = "logs"
	}

	if batchSize, ok := config["batch_size"].(float64); ok {
		a.BatchSize = int(batchSize)
	} else {
		a.BatchSize = 100
	}

	if flushTimeout, ok := config["flush_timeout"].(float64); ok {
		a.FlushTimeout = time.Duration(flushTimeout) * time.Second
	} else {
		a.FlushTimeout = 5 * time.Second
	}

	// 初始化日志缓冲区
	a.buffer = make([]logger.LogEntry, 0, a.BatchSize)

	// 连接到Kafka
	// 实际应该这样:
	// config := sarama.NewConfig()
	// config.Producer.Return.Successes = true
	// producer, err := sarama.NewSyncProducer(a.Brokers, config)
	// if err != nil {
	//     return fmt.Errorf("failed to create kafka producer: %v", err)
	// }
	// a.producer = producer

	// 定期刷新缓冲区
	go a.flushPeriodically()

	return nil
}

// Process 处理日志条目
func (a *KafkaAdapter) Process(ctx context.Context, entry logger.LogEntry) error {
	a.bufferMu.Lock()
	defer a.bufferMu.Unlock()

	// 添加到缓冲区
	a.buffer = append(a.buffer, entry)

	// 如果达到批量大小，刷新缓冲区
	if len(a.buffer) >= a.BatchSize {
		return a.flushBuffer()
	}

	return nil
}

// Flush 刷新缓冲区
func (a *KafkaAdapter) Flush() error {
	a.bufferMu.Lock()
	defer a.bufferMu.Unlock()

	return a.flushBuffer()
}

// flushBuffer 刷新缓冲区（无锁版本，调用前需要获取锁）
func (a *KafkaAdapter) flushBuffer() error {
	if len(a.buffer) == 0 {
		return nil
	}

	// 在实际应用中，这里应该批量发送到Kafka
	// var messages []*sarama.ProducerMessage
	// for _, entry := range a.buffer {
	//     data, _ := json.Marshal(entry)
	//     messages = append(messages, &sarama.ProducerMessage{
	//         Topic: a.Topic,
	//         Value: sarama.ByteEncoder(data),
	//     })
	// }
	//
	// for _, msg := range messages {
	//     _, _, err := a.producer.(*sarama.SyncProducer).SendMessage(msg)
	//     if err != nil {
	//         return err
	//     }
	// }

	// 这里仅作演示，实际打印日志
	for _, entry := range a.buffer {
		data, _ := json.Marshal(entry)
		fmt.Printf("[Kafka Adapter] Would send to topic %s: %s\n", a.Topic, string(data))
	}

	// 清空缓冲区
	a.buffer = a.buffer[:0]

	return nil
}

// flushPeriodically 定期刷新缓冲区
func (a *KafkaAdapter) flushPeriodically() {
	ticker := time.NewTicker(a.FlushTimeout)
	defer ticker.Stop()

	for range ticker.C {
		_ = a.Flush()
	}
}

// Close 关闭适配器
func (a *KafkaAdapter) Close() error {
	// 刷新剩余日志
	err := a.Flush()

	// 关闭Kafka生产者
	// if a.producer != nil {
	//     a.producer.(*sarama.SyncProducer).Close()
	// }

	return err
}
