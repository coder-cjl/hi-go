package aiservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"hi-go/src/config"
)

// DeepSeekClient DeepSeek客户端
type DeepSeekClient struct {
	config config.AIConfig
	client *http.Client
}

// NewDeepSeekClient 创建DeepSeek客户端
func NewDeepSeekClient(cfg config.AIConfig) *DeepSeekClient {
	return &DeepSeekClient{
		config: cfg,
		client: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

// Chat 对话
func (c *DeepSeekClient) Chat(ctx context.Context, messages []Message, tools []Tool) (*ChatResponse, error) {
	reqBody := map[string]interface{}{
		"model":       c.config.Model,
		"messages":    messages,
		"temperature": c.config.Temperature,
		"max_tokens":  c.config.MaxTokens,
		"stream":      false,
	}

	// 添加工具定义
	if len(tools) > 0 {
		reqBody["tools"] = tools
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("deepseek api error (status %d): %s", resp.StatusCode, string(body))
	}

	var deepseekResp DeepSeekResponse
	if err := json.Unmarshal(body, &deepseekResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w, body: %s", err, string(body))
	}

	if len(deepseekResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from deepseek")
	}

	choice := deepseekResp.Choices[0]
	result := &ChatResponse{
		Content:      choice.Message.Content,
		ToolCalls:    choice.Message.ToolCalls,
		FinishReason: choice.FinishReason,
	}

	return result, nil
}

// ChatStream 流式对话
func (c *DeepSeekClient) ChatStream(ctx context.Context, messages []Message, tools []Tool) (<-chan StreamResponse, error) {
	responseChan := make(chan StreamResponse, 10)

	reqBody := map[string]interface{}{
		"model":       c.config.Model,
		"messages":    messages,
		"temperature": c.config.Temperature,
		"max_tokens":  c.config.MaxTokens,
		"stream":      true, // 启用流式响应
	}

	// 添加工具定义
	if len(tools) > 0 {
		reqBody["tools"] = tools
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		close(responseChan)
		return responseChan, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		close(responseChan)
		return responseChan, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.client.Do(req)
	if err != nil {
		close(responseChan)
		return responseChan, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		close(responseChan)
		return responseChan, fmt.Errorf("deepseek api error (status %d): %s", resp.StatusCode, string(body))
	}

	// 启动goroutine读取SSE流
	go func() {
		defer close(responseChan)
		defer resp.Body.Close()

		reader := &sseReader{reader: resp.Body}
		var accumulatedContent string
		var accumulatedToolCalls []ToolCall

		for {
			select {
			case <-ctx.Done():
				responseChan <- StreamResponse{Error: ctx.Err()}
				return
			default:
				line, err := reader.readLine()
				if err != nil {
					if err != io.EOF {
						responseChan <- StreamResponse{Error: err}
					}
					return
				}

				// SSE格式: data: {...}
				if len(line) < 6 || line[:6] != "data: " {
					continue
				}

				data := line[6:]
				if data == "[DONE]" {
					// 流结束
					return
				}

				var chunk DeepSeekStreamChunk
				if err := json.Unmarshal([]byte(data), &chunk); err != nil {
					continue
				}

				if len(chunk.Choices) == 0 {
					continue
				}

				delta := chunk.Choices[0].Delta

				// 累积内容
				if delta.Content != "" {
					accumulatedContent += delta.Content
				}

				// 累积工具调用
				if len(delta.ToolCalls) > 0 {
					accumulatedToolCalls = append(accumulatedToolCalls, delta.ToolCalls...)
				}

				// 发送流式响应
				streamResp := StreamResponse{
					Content:      delta.Content,
					ToolCalls:    delta.ToolCalls,
					FinishReason: chunk.Choices[0].FinishReason,
				}

				select {
				case responseChan <- streamResp:
				case <-ctx.Done():
					return
				}

				// 如果完成，发送最终累积结果
				if chunk.Choices[0].FinishReason != "" {
					finalResp := StreamResponse{
						Content:      accumulatedContent,
						ToolCalls:    accumulatedToolCalls,
						FinishReason: chunk.Choices[0].FinishReason,
					}
					select {
					case responseChan <- finalResp:
					case <-ctx.Done():
					}
					return
				}
			}
		}
	}()

	return responseChan, nil
}

// sseReader SSE流读取器
type sseReader struct {
	reader io.Reader
	buffer []byte
}

func (r *sseReader) readLine() (string, error) {
	var line []byte
	buf := make([]byte, 1)

	for {
		n, err := r.reader.Read(buf)
		if err != nil {
			if len(line) > 0 {
				return string(line), nil
			}
			return "", err
		}

		if n > 0 {
			if buf[0] == '\n' {
				return string(line), nil
			}
			if buf[0] != '\r' {
				line = append(line, buf[0])
			}
		}
	}
}

// DeepSeekStreamChunk 流式响应块
type DeepSeekStreamChunk struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role      string     `json:"role,omitempty"`
			Content   string     `json:"content,omitempty"`
			ToolCalls []ToolCall `json:"tool_calls,omitempty"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// DeepSeek API 响应结构
type DeepSeekResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role      string     `json:"role"`
			Content   string     `json:"content"`
			ToolCalls []ToolCall `json:"tool_calls,omitempty"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
