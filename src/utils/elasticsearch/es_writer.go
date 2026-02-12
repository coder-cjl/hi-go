package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// ESWriter 实现了 io.Writer 接口，用于将日志写入 Elasticsearch
type ESWriter struct {
	client    *elasticsearch.Client
	index     string
	buffer    []LogEntry
	mutex     sync.Mutex
	batchSize int
	ticker    *time.Ticker
	ctx       context.Context
	cancel    context.CancelFunc
}

// LogEntry 日志条目结构
type LogEntry struct {
	Timestamp string                 `json:"@timestamp"`
	Level     string                 `json:"level"`
	Logger    string                 `json:"logger"`
	Caller    string                 `json:"caller"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// Config Elasticsearch Writer 配置
type Config struct {
	Addrs     []string // ES 集群地址
	Username  string   // 用户名（可选）
	Password  string   // 密码（可选）
	Index     string   // 索引名称
	MaxRetry  int      // 最大重试次数
	BatchSize int      // 批量写入大小
	FlushTime int      // 刷新间隔（秒）
}

// NewESWriter 创建一个新的 Elasticsearch Writer
func NewESWriter(cfg *Config) (*ESWriter, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if len(cfg.Addrs) == 0 {
		return nil, fmt.Errorf("elasticsearch addresses cannot be empty")
	}

	if cfg.Index == "" {
		cfg.Index = "logs"
	}

	if cfg.MaxRetry == 0 {
		cfg.MaxRetry = 3
	}

	if cfg.BatchSize == 0 {
		cfg.BatchSize = 100
	}

	if cfg.FlushTime == 0 {
		cfg.FlushTime = 5
	}

	// 创建 ES 客户端
	esCfg := elasticsearch.Config{
		Addresses:  cfg.Addrs,
		Username:   cfg.Username,
		Password:   cfg.Password,
		MaxRetries: cfg.MaxRetry,
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	// 测试连接
	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch returned error: %s", res.String())
	}

	ctx, cancel := context.WithCancel(context.Background())

	esWriter := &ESWriter{
		client:    client,
		index:     cfg.Index,
		buffer:    make([]LogEntry, 0, cfg.BatchSize),
		batchSize: cfg.BatchSize,
		ticker:    time.NewTicker(time.Duration(cfg.FlushTime) * time.Second),
		ctx:       ctx,
		cancel:    cancel,
	}

	// 启动定时刷新
	go esWriter.autoFlush()

	return esWriter, nil
}

// Write 实现 io.Writer 接口
func (w *ESWriter) Write(p []byte) (n int, err error) {
	// 解析日志条目
	var entry map[string]interface{}
	if err := json.Unmarshal(p, &entry); err != nil {
		// 如果解析失败，创建简单的日志条目
		entry = map[string]interface{}{
			"message":    string(p),
			"@timestamp": time.Now().Format(time.RFC3339),
		}
	}

	// 确保有时间戳
	if _, ok := entry["@timestamp"]; !ok {
		if ts, ok := entry["ts"]; ok {
			entry["@timestamp"] = ts
		} else if ts, ok := entry["time"]; ok {
			entry["@timestamp"] = ts
		} else {
			entry["@timestamp"] = time.Now().Format(time.RFC3339)
		}
	}

	// 转换为 LogEntry
	logEntry := LogEntry{
		Timestamp: fmt.Sprintf("%v", entry["@timestamp"]),
		Level:     fmt.Sprintf("%v", entry["level"]),
		Logger:    fmt.Sprintf("%v", entry["logger"]),
		Caller:    fmt.Sprintf("%v", entry["caller"]),
		Message:   fmt.Sprintf("%v", entry["msg"]),
		Fields:    make(map[string]interface{}),
	}

	// 提取其他字段
	for k, v := range entry {
		if k != "@timestamp" && k != "level" && k != "logger" && k != "caller" && k != "msg" {
			logEntry.Fields[k] = v
		}
	}

	w.mutex.Lock()
	w.buffer = append(w.buffer, logEntry)
	shouldFlush := len(w.buffer) >= w.batchSize
	w.mutex.Unlock()

	if shouldFlush {
		if err := w.flush(); err != nil {
			// 记录错误但不中断写入
			fmt.Printf("Failed to flush logs to elasticsearch: %v\n", err)
		}
	}

	return len(p), nil
}

// flush 刷新缓冲区到 Elasticsearch
func (w *ESWriter) flush() error {
	w.mutex.Lock()
	if len(w.buffer) == 0 {
		w.mutex.Unlock()
		return nil
	}

	// 复制缓冲区并清空
	entries := make([]LogEntry, len(w.buffer))
	copy(entries, w.buffer)
	w.buffer = w.buffer[:0]
	w.mutex.Unlock()

	// 批量索引
	var buf bytes.Buffer
	for _, entry := range entries {
		// 写入元数据
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": w.index,
			},
		}
		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return fmt.Errorf("failed to encode metadata: %w", err)
		}

		// 写入文档
		if err := json.NewEncoder(&buf).Encode(entry); err != nil {
			return fmt.Errorf("failed to encode log entry: %w", err)
		}
	}

	// 执行批量索引
	req := esapi.BulkRequest{
		Body: bytes.NewReader(buf.Bytes()),
	}

	res, err := req.Do(context.Background(), w.client)
	if err != nil {
		return fmt.Errorf("bulk request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk request returned error: %s", res.String())
	}

	return nil
}

// autoFlush 自动定时刷新
func (w *ESWriter) autoFlush() {
	for {
		select {
		case <-w.ticker.C:
			if err := w.flush(); err != nil {
				fmt.Printf("Auto flush failed: %v\n", err)
			}
		case <-w.ctx.Done():
			return
		}
	}
}

// Close 关闭 writer 并刷新剩余日志
func (w *ESWriter) Close() error {
	w.cancel()
	w.ticker.Stop()
	return w.flush()
}

// Sync 同步刷新缓冲区
func (w *ESWriter) Sync() error {
	return w.flush()
}
