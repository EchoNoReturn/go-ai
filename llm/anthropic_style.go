package llm

import "errors"

type AnthropicChatRequest struct {
	Model         string       `json:"model"`    // 必选
	Messages      []LLMMessage `json:"messages"` // 必选
	System        interface{}  `json:"system,omitempty"`
	MaxTokens     *int         `json:"max_tokens,omitempty"`
	Temperature   *float64     `json:"temperature,omitempty"`
	TopP          *float64     `json:"top_p,omitempty"`
	StopSequences []string     `json:"stop_sequences,omitempty"`
	Stream        *bool        `json:"stream,omitempty"`
	Thinking      interface{}  `json:"thinking,omitempty"`
	ToolChoice    interface{}  `json:"tool_choice,omitempty"`
	Tools         []LLMTool    `json:"tools,omitempty"`
}

func (r *AnthropicChatRequest) Validate() (bool, error) {
	if r.Model == "" {
		return false, errors.New("model is required")
	}
	if len(r.Messages) == 0 {
		return false, errors.New("at least one message is required")
	}
	if r.MaxTokens != nil {
		if *r.MaxTokens < 0 || *r.MaxTokens > 131072 {
			return false, errors.New("max_tokens must be between 0 and 131072")
		}
	}
	if r.Temperature != nil {
		if *r.Temperature < 0 || *r.Temperature > 1.5 {
			return false, errors.New("temperature must be between 0 and 1.5")
		}
	}
	if r.TopP != nil {
		if *r.TopP < 0.01 || *r.TopP > 1.0 {
			return false, errors.New("top_p must be between 0.01 and 1.0")
		}
	}
	return true, nil
}

// ======== AnthropicResponse 定义 ========

// Anthropic风格响应结构体
type AnthropicChatResponse struct {
	ID         string              `json:"id" yaml:"id"`                           // 响应唯一标识
	Type       string              `json:"type" yaml:"type"`                       // 类型，如"completion"
	Model      string              `json:"model" yaml:"model"`                     // 模型名称
	Content    string              `json:"content" yaml:"content"`                 // 生成内容
	StopReason string              `json:"stop_reason" yaml:"stop_reason"`         // 结束原因
	Usage      *AnthropicChatUsage `json:"usage,omitempty" yaml:"usage,omitempty"` // token统计
}

type AnthropicChatUsage struct {
	InputTokens  int `json:"input_tokens" yaml:"input_tokens"`   // 输入tokens
	OutputTokens int `json:"output_tokens" yaml:"output_tokens"` // 输出tokens
}

// 流式响应结构体
type AnthropicChatStreamResponse struct {
	ID         string              `json:"id" yaml:"id"`
	Type       string              `json:"type" yaml:"type"`
	Model      string              `json:"model" yaml:"model"`
	Content    string              `json:"content" yaml:"content"` // 增量内容
	StopReason string              `json:"stop_reason" yaml:"stop_reason"`
	Usage      *AnthropicChatUsage `json:"usage,omitempty" yaml:"usage,omitempty"`
}
