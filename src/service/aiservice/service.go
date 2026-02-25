package aiservice

import (
	"context"
	"encoding/json"
	"fmt"

	"hi-go/src/config"
	"hi-go/src/utils/logger"

	"go.uber.org/zap"
)

// Service AI服务
type Service struct {
	config   config.AIConfig
	registry *SkillRegistry
	client   AIClient
}

// AIClient AI客户端接口
type AIClient interface {
	Chat(ctx context.Context, messages []Message, tools []Tool) (*ChatResponse, error)
	ChatStream(ctx context.Context, messages []Message, tools []Tool) (<-chan StreamResponse, error)
}

// StreamResponse 流式响应
type StreamResponse struct {
	Content      string     `json:"content"`
	ToolCalls    []ToolCall `json:"tool_calls,omitempty"`
	FinishReason string     `json:"finish_reason"`
	Error        error      `json:"-"`
}

// Message 消息结构
type Message struct {
	Role       string     `json:"role"`                   // system, user, assistant, tool
	Content    string     `json:"content"`                // 消息内容
	ToolCallID string     `json:"tool_call_id,omitempty"` // 工具调用ID（仅在 role=tool 时使用）
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`   // 工具调用列表（仅在 role=assistant 时使用）
}

// Tool 工具定义
type Tool struct {
	Type     string   `json:"type"`     // 固定为 "function"
	Function Function `json:"function"` // 函数定义
}

// Function 函数定义
type Function struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ToolCall 工具调用
type ToolCall struct {
	ID       string       `json:"id"`       // 工具调用ID
	Type     string       `json:"type"`     // 固定为 "function"
	Function FunctionCall `json:"function"` // 函数调用详情
}

// FunctionCall 函数调用
type FunctionCall struct {
	Name      string `json:"name"`      // 函数名
	Arguments string `json:"arguments"` // JSON格式的参数
}

// ChatResponse 对话响应
type ChatResponse struct {
	Content      string     `json:"content"`
	ToolCalls    []ToolCall `json:"tool_calls,omitempty"`
	FinishReason string     `json:"finish_reason"`
}

// NewService 创建AI服务
func NewService(cfg config.AIConfig, registry *SkillRegistry, client AIClient) *Service {
	return &Service{
		config:   cfg,
		registry: registry,
		client:   client,
	}
}

// Chat 对话
func (s *Service) Chat(ctx context.Context, userMessage string) (string, error) {
	messages := []Message{
		{
			Role:    "system",
			Content: s.config.SystemPrompt,
		},
		{
			Role:    "user",
			Content: userMessage,
		},
	}

	// 获取所有可用的技能作为工具
	tools := s.getTools()

	// 最多执行5轮对话（防止无限循环）
	maxRounds := 5
	for i := 0; i < maxRounds; i++ {
		// 调用AI
		resp, err := s.client.Chat(ctx, messages, tools)
		if err != nil {
			return "", err
		}

		// 如果没有工具调用，直接返回结果
		if len(resp.ToolCalls) == 0 {
			return resp.Content, nil
		}

		// 执行工具调用
		toolResults, err := s.executeToolCalls(ctx, resp.ToolCalls)
		if err != nil {
			logger.Error("执行工具调用失败", zap.Error(err))
			return "", fmt.Errorf("执行工具调用失败: %w", err)
		}

		// 将助手的响应（包含工具调用）添加到消息历史
		messages = append(messages, Message{
			Role:      "assistant",
			Content:   resp.Content,
			ToolCalls: resp.ToolCalls,
		})

		// 将工具调用结果添加到消息历史
		messages = append(messages, toolResults...)
	}

	return "", fmt.Errorf("超过最大对话轮数")
}

// ChatStream 流式对话（支持工具调用的多轮对话）
func (s *Service) ChatStream(ctx context.Context, userMessage string) (<-chan StreamResponse, error) {
	outputChan := make(chan StreamResponse, 10)

	go func() {
		defer close(outputChan)

		messages := []Message{
			{
				Role:    "system",
				Content: s.config.SystemPrompt,
			},
			{
				Role:    "user",
				Content: userMessage,
			},
		}

		// 获取所有可用的技能作为工具
		tools := s.getTools()

		// 最多执行5轮对话（防止无限循环）
		maxRounds := 5
		for round := 0; round < maxRounds; round++ {
			// 调用AI流式接口
			streamChan, err := s.client.ChatStream(ctx, messages, tools)
			if err != nil {
				outputChan <- StreamResponse{Error: err}
				return
			}

			var accumulatedContent string
			var accumulatedToolCalls []ToolCall
			var lastFinishReason string

			// 读取流式响应
			for resp := range streamChan {
				if resp.Error != nil {
					outputChan <- resp
					return
				}

				// 累积内容
				if resp.Content != "" {
					accumulatedContent += resp.Content
					// 转发给客户端
					outputChan <- StreamResponse{
						Content: resp.Content,
					}
				}

				// 累积工具调用
				if len(resp.ToolCalls) > 0 {
					accumulatedToolCalls = append(accumulatedToolCalls, resp.ToolCalls...)
				}

				// 记录完成原因
				if resp.FinishReason != "" {
					lastFinishReason = resp.FinishReason
				}
			}

			// 如果没有工具调用，说明对话完成
			if len(accumulatedToolCalls) == 0 {
				// 发送完成信号
				outputChan <- StreamResponse{
					FinishReason: lastFinishReason,
				}
				return
			}

			// 有工具调用，执行工具
			logger.Info("检测到工具调用，开始执行", zap.Int("count", len(accumulatedToolCalls)))

			// 将助手的响应（包含工具调用）添加到消息历史
			messages = append(messages, Message{
				Role:      "assistant",
				Content:   accumulatedContent,
				ToolCalls: accumulatedToolCalls,
			})

			// 执行工具调用
			toolResults, err := s.executeToolCalls(ctx, accumulatedToolCalls)
			if err != nil {
				logger.Error("执行工具调用失败", zap.Error(err))
				outputChan <- StreamResponse{
					Error: fmt.Errorf("执行工具调用失败: %w", err),
				}
				return
			}

			// 将工具调用结果添加到消息历史
			messages = append(messages, toolResults...)

			// 继续下一轮对话，AI会基于工具结果生成最终回答
		}

		// 超过最大轮数
		outputChan <- StreamResponse{
			Error: fmt.Errorf("超过最大对话轮数"),
		}
	}()

	return outputChan, nil
}

// executeToolCalls 执行工具调用
func (s *Service) executeToolCalls(ctx context.Context, toolCalls []ToolCall) ([]Message, error) {
	results := make([]Message, 0, len(toolCalls))

	for _, tc := range toolCalls {
		skill, ok := s.registry.Get(tc.Function.Name)
		if !ok {
			return nil, fmt.Errorf("skill not found: %s", tc.Function.Name)
		}

		// 解析参数
		var params map[string]interface{}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &params); err != nil {
			return nil, fmt.Errorf("invalid function arguments: %w", err)
		}

		// 执行技能
		result, err := skill.Execute(ctx, params)
		if err != nil {
			logger.Error("技能执行失败", zap.String("skill", tc.Function.Name), zap.Error(err))
			return nil, fmt.Errorf("skill execution failed: %w", err)
		}

		// 将结果转为JSON
		resultJSON, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}

		// 添加工具调用结果
		results = append(results, Message{
			Role:       "tool",
			Content:    string(resultJSON),
			ToolCallID: tc.ID,
		})
	}

	return results, nil
}

// getTools 获取所有技能的工具定义
func (s *Service) getTools() []Tool {
	skills := s.registry.GetAll()
	tools := make([]Tool, 0, len(skills))

	for _, skill := range skills {
		tools = append(tools, Tool{
			Type: "function",
			Function: Function{
				Name:        skill.Name(),
				Description: skill.Description(),
				Parameters:  skill.Parameters(),
			},
		})
	}

	return tools
}
