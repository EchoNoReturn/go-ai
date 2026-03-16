package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type LLMSession struct {
	Endpoint         string
	ApiKey           string
	Type             LLMType
	MessageList      []LLMMessage
	Tools            []LLMTool
	Model            string
	Thinking         *OpenAIThinking
	FrequencyPenalty *float64
	MaxTokens        *int
	PresencePenalty  *float64
	ResponseFormat   *OpenAIResponseFormat
	Stop             interface{}
	Stream           *bool
	StreamOptions    interface{}
	Temperature      *float64
	TopP             *float64
	ToolChoice       interface{}
	Logprobs         *bool
	TopLogprobs      *int
	Header           map[string]string // 新增：自定义请求头
}

// ListModels 查询支持的模型列表（OpenAI/DeepSeek），Anthropic不支持
func (c *LLMSession) ListModels() (string, error) {
	var url string
	switch c.Type {
	case OpenAI:
		url = c.Endpoint + "/models"
	case Anthropic:
		return "", errors.New("Anthropic API does not support listing models")
	default:
		return "", errors.New("unsupported LLM type")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)
	for k, v := range c.Header {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (s *LLMSession) SetHeader(key, value string) *LLMSession {
	if s.Header == nil {
		s.Header = make(map[string]string)
	}
	s.Header[key] = value
	return s
}

// 支持 MessageList 复制、覆盖、引用替换
func (c *LLMSession) SetMessageListCopy(list []LLMMessage) *LLMSession {
	c.MessageList = append([]LLMMessage(nil), list...)
	return c
}

func (c *LLMSession) AppendSystemMessage(content, name string) *LLMSession {
	systemMessage := CreateMessage(MSystem, content)
	if msg, ok := systemMessage.(*SystemMessage); ok {
		msg.Name = name
	}
	c.MessageList = append(c.MessageList, systemMessage)
	return c
}

func (c *LLMSession) AppendUserMessage(content, name string) *LLMSession {
	userMessage := CreateMessage(MUser, content)
	if msg, ok := userMessage.(*UserMessage); ok {
		msg.Name = name
	}
	c.MessageList = append(c.MessageList, userMessage)
	return c
}

func (c *LLMSession) AppendAssistantMessage(content, name, reasoningContent string, prefix bool) *LLMSession {
	assistantMessage := CreateMessage(MAssistant, content)
	if msg, ok := assistantMessage.(*AssistantMessage); ok {
		msg.Name = name
		msg.Prefix = &prefix
		msg.ReasoningContent = reasoningContent
	}
	c.MessageList = append(c.MessageList, assistantMessage)
	return c
}

func (c *LLMSession) AppendToolMessage(content, toolCallId string) *LLMSession {
	toolMessage := CreateMessage(MTool, content)
	if msg, ok := toolMessage.(*ToolMessage); ok {
		msg.ToolCallId = toolCallId
	}
	c.MessageList = append(c.MessageList, toolMessage)
	return c
}

func (c *LLMSession) RunChat() (string, error) {
	// 检查 Model 是否设置
	if c.Model == "" {
		return "", errors.New("model is required, but not set")
	}
	// 根据类型构建请求体
	var reqBody []byte
	var err error
	var url string
	var respContent string
	switch c.Type {
	case OpenAI:
		// OpenAI/DeepSeek风格
		openaiReq := OpenAIChatRequestBody{
			Model:            c.Model,
			Messages:         toMessageItemSlice(c.MessageList),
			Thinking:         c.Thinking,
			FrequencyPenalty: c.FrequencyPenalty,
			MaxTokens:        c.MaxTokens,
			PresencePenalty:  c.PresencePenalty,
			ResponseFormat:   c.ResponseFormat,
			Stop:             c.Stop,
			Stream:           c.Stream,
			StreamOptions:    c.StreamOptions,
			Temperature:      c.Temperature,
			TopP:             c.TopP,
			Tools:            c.Tools,
			ToolChoice:       c.ToolChoice,
			Logprobs:         c.Logprobs,
			TopLogprobs:      c.TopLogprobs,
		}
		reqBody, err = json.Marshal(openaiReq)
		url = c.Endpoint + "/chat/completions"
	case Anthropic:
		// Anthropic风格
		anthropicReq := AnthropicChatRequest{
			Model:         c.Model,
			Messages:      c.MessageList,
			System:        c.Thinking, // 可根据实际需求调整
			MaxTokens:     c.MaxTokens,
			Temperature:   c.Temperature,
			TopP:          c.TopP,
			StopSequences: nil,
			Stream:        c.Stream,
			Thinking:      c.Thinking,
			ToolChoice:    c.ToolChoice,
			Tools:         c.Tools,
		}
		reqBody, err = json.Marshal(anthropicReq)
		url = c.Endpoint + "/anthropic/v1/messages"
	default:
		return "", errors.New("unsupported LLM type")
	}
	if err != nil {
		return "", err
	}
	// 发起 HTTP 请求
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)
	// 合并自定义 Header
	for k, v := range c.Header {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	respContent = string(body)
	return respContent, nil
}

// 辅助函数：将 []LLMMessage 转换为 []MessageItem
func toMessageItemSlice(msgs []LLMMessage) []MessageItem {
	items := make([]MessageItem, 0, len(msgs))
	for _, m := range msgs {
		if item, ok := m.(*SystemMessage); ok {
			items = append(items, item.MessageItem)
		} else if item, ok := m.(*UserMessage); ok {
			items = append(items, item.MessageItem)
		} else if item, ok := m.(*AssistantMessage); ok {
			items = append(items, item.MessageItem)
		} else if item, ok := m.(*ToolMessage); ok {
			items = append(items, item.MessageItem)
		}
	}
	return items
}

func NewSession(apiKey, endpoint string) (*LLMSession, error) {
	// 必填参数校验
	if apiKey == "" {
		return nil, errors.New("apiKey is required")
	}
	if endpoint == "" {
		return nil, errors.New("endpoint is required")
	}

	return &LLMSession{
		Endpoint:    endpoint,
		ApiKey:      apiKey,
		Type:        OpenAI,
		MessageList: make([]LLMMessage, 0),
		Header:      make(map[string]string),
	}, nil
}
