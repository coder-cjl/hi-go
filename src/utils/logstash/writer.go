package logstash

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

// Writer 实现了 io.Writer 接口，用于将日志通过 TCP 发送到 Logstash
type Writer struct {
	conn       net.Conn
	host       string
	port       int
	protocol   string
	timeout    time.Duration
	reconnect  bool
	mutex      sync.Mutex
	connected  bool
	buffer     []byte
	bufferSize int
}

// Config Logstash Writer 配置
type Config struct {
	Host       string // Logstash 服务器地址
	Port       int    // Logstash TCP 端口
	Protocol   string // 协议：tcp 或 udp（默认 tcp）
	Timeout    int    // 连接超时（秒，默认 5）
	Reconnect  bool   // 是否自动重连（默认 true）
	BufferSize int    // 缓冲区大小（默认 8192）
}

// NewWriter 创建一个新的 Logstash Writer
func NewWriter(cfg *Config) (*Writer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if cfg.Host == "" {
		return nil, fmt.Errorf("logstash host cannot be empty")
	}

	if cfg.Port == 0 {
		cfg.Port = 5000
	}

	if cfg.Protocol == "" {
		cfg.Protocol = "tcp"
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 5
	}

	if cfg.BufferSize == 0 {
		cfg.BufferSize = 8192
	}

	writer := &Writer{
		host:       cfg.Host,
		port:       cfg.Port,
		protocol:   cfg.Protocol,
		timeout:    time.Duration(cfg.Timeout) * time.Second,
		reconnect:  cfg.Reconnect,
		buffer:     make([]byte, 0, cfg.BufferSize),
		bufferSize: cfg.BufferSize,
	}

	// 初始连接
	if err := writer.connect(); err != nil {
		// 如果启用了重连，初始连接失败不报错
		if !cfg.Reconnect {
			return nil, fmt.Errorf("failed to connect to logstash: %w", err)
		}
		// 记录错误但继续
		fmt.Printf("Warning: initial connection to logstash failed: %v\n", err)
	}

	return writer, nil
}

// connect 连接到 Logstash
func (w *Writer) connect() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// 如果已经连接，先关闭
	if w.conn != nil {
		w.conn.Close()
		w.conn = nil
		w.connected = false
	}

	addr := net.JoinHostPort(w.host, fmt.Sprintf("%d", w.port))
	conn, err := net.DialTimeout(w.protocol, addr, w.timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", addr, err)
	}

	w.conn = conn
	w.connected = true
	return nil
}

// Write 实现 io.Writer 接口
func (w *Writer) Write(p []byte) (n int, err error) {
	// 检查连接状态
	if !w.connected {
		if w.reconnect {
			if err := w.connect(); err != nil {
				// 重连失败，返回错误但不中断
				return len(p), nil
			}
		} else {
			return 0, fmt.Errorf("not connected to logstash")
		}
	}

	// 解析日志为 JSON，确保每行都是有效的 JSON
	var logEntry map[string]interface{}
	if err := json.Unmarshal(p, &logEntry); err != nil {
		// 如果不是 JSON，创建简单的日志条目
		logEntry = map[string]interface{}{
			"message":    string(p),
			"@timestamp": time.Now().Format(time.RFC3339),
		}
	}

	// 重新序列化为 JSON
	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal log entry: %w", err)
	}

	// 添加换行符（Logstash 的 json_lines codec 需要）
	jsonData = append(jsonData, '\n')

	// 发送数据
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.conn == nil {
		if w.reconnect {
			// 尝试重连
			w.mutex.Unlock()
			err := w.connect()
			w.mutex.Lock()
			if err != nil {
				return len(p), nil
			}
		} else {
			return 0, fmt.Errorf("no connection to logstash")
		}
	}

	// 设置写入超时
	if err := w.conn.SetWriteDeadline(time.Now().Add(w.timeout)); err != nil {
		w.connected = false
		return 0, fmt.Errorf("failed to set write deadline: %w", err)
	}

	// 写入数据
	_, err = w.conn.Write(jsonData)
	if err != nil {
		w.connected = false
		w.conn.Close()
		w.conn = nil

		// 如果启用了重连，不报错
		if w.reconnect {
			return len(p), nil
		}
		return 0, fmt.Errorf("failed to write to logstash: %w", err)
	}

	// 返回原始数据长度
	return len(p), nil
}

// Close 关闭连接
func (w *Writer) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.conn != nil {
		err := w.conn.Close()
		w.conn = nil
		w.connected = false
		return err
	}

	return nil
}

// Sync 刷新缓冲区（对于网络连接，这是一个空操作）
func (w *Writer) Sync() error {
	return nil
}

// IsConnected 返回连接状态
func (w *Writer) IsConnected() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.connected
}
