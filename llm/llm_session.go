package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
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
	StreamOptions    interface{}
	Temperature      *float64
	TopP             *float64
	ToolChoice       interface{}
	Logprobs         *bool
	TopLogprobs      *int
	Header           map[string]string // 新增：自定义请求头
}

// ======== 策略拆分辅助函数 ========
// 构建请求体
func (c *LLMSession) buildRequestBody(stream bool, requestMethodInfo RequestMethodInfo) ([]byte, string, error) {
	var reqBody []byte
	var err error
	var url string
	switch c.Type {
	case OpenAI:
		openaiReq := OpenAIChatRequestBody{
			Model:            c.Model,
			Messages:         toMessageItemSlice(c.MessageList),
			Thinking:         c.Thinking,
			FrequencyPenalty: c.FrequencyPenalty,
			MaxTokens:        c.MaxTokens,
			PresencePenalty:  c.PresencePenalty,
			ResponseFormat:   c.ResponseFormat,
			Stop:             c.Stop,
			Stream:           &stream,
			StreamOptions:    c.StreamOptions,
			Temperature:      c.Temperature,
			TopP:             c.TopP,
			Tools:            c.Tools,
			ToolChoice:       c.ToolChoice,
			Logprobs:         c.Logprobs,
			TopLogprobs:      c.TopLogprobs,
		}
		reqBody, err = json.Marshal(openaiReq)
		url = c.Endpoint + requestMethodInfo.Path
	case Anthropic:
		anthropicReq := AnthropicChatRequest{
			Model:         c.Model,
			Messages:      c.MessageList,
			System:        c.Thinking,
			MaxTokens:     c.MaxTokens,
			Temperature:   c.Temperature,
			TopP:          c.TopP,
			StopSequences: nil,
			Stream:        &stream,
			Thinking:      c.Thinking,
			ToolChoice:    c.ToolChoice,
			Tools:         c.Tools,
		}
		reqBody, err = json.Marshal(anthropicReq)
		url = c.Endpoint + requestMethodInfo.Path
	default:
		return nil, "", errors.New("unsupported LLM type")
	}
	return reqBody, url, err
}

// 发送HTTP请求
func (c *LLMSession) sendRequest(url string, reqBody []byte, reqMethodInfo RequestMethodInfo) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(reqMethodInfo.Method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)
	for k, v := range c.Header {
		req.Header.Set(k, v)
	}
	return client.Do(req)
}

// 解析响应体
func (c *LLMSession) parseResponse(respBody []byte) (any, error) {
	switch c.Type {
	case OpenAI:
		var result OpenAIChatResponseBody
		if err := json.Unmarshal(respBody, &result); err != nil {
			return nil, err
		}
		return &result, nil
	case Anthropic:
		var result AnthropicChatResponse
		if err := json.Unmarshal(respBody, &result); err != nil {
			return nil, err
		}
		return &result, nil
	default:
		return nil, errors.New("unsupported LLM type")
	}
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

func checkResp(resp *http.Response) error {
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var bodyJson map[string]interface{}
		json.Unmarshal(body, &bodyJson)
		reason := bodyJson["error"].(map[string]interface{})["message"].(string)
		return errors.New("unexpected status code: " + resp.Status + ", reason: " + reason)
	}
	return nil
}

func (c *LLMSession) RunChat() (interface{}, error) {
	// 检查 Model 是否设置
	if c.Model == "" {
		return nil, errors.New("model is required, but not set")
	}
	requestMethodInfo := GetRequestInfo(c.Type, "chat")
	// 构建请求体
	reqBody, url, err := c.buildRequestBody(false, requestMethodInfo)
	if err != nil {
		return nil, err
	}
	// 发送请求
	resp, err := c.sendRequest(url, reqBody, requestMethodInfo)
	if err != nil {
		return nil, err
	}
	// 非流式响应直接解析并返回
	defer resp.Body.Close()

	if err := checkResp(resp); err != nil {
		return nil, err
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// 解析响应，返回结构体对象
	respObj, err := c.parseResponse(respBytes)
	if err != nil {
		return nil, err
	}
	return respObj, nil
}

func (c *LLMSession) RunChatStream(noHandle bool) (<-chan any, <-chan error) {
	resultChan := make(chan any)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		// 检查 Model 是否设置
		if c.Model == "" {
			errorChan <- errors.New("model is required, but not set")
			return
		}
		// 构建请求体
		requestMethodInfo := GetRequestInfo(c.Type, "chat")
		reqBody, url, err := c.buildRequestBody(true, requestMethodInfo)
		if err != nil {
			errorChan <- err
			return
		}

		// 发送请求
		resp, err := c.sendRequest(url, reqBody, requestMethodInfo)
		if err != nil {
			errorChan <- err
			return
		}
		defer resp.Body.Close()

		// 检查响应状态码
		if err := checkResp(resp); err != nil {
			errorChan <- err
			return
		}

		// 流式响应解析
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					return
				}
				errorChan <- err
				return
			}

			line = strings.TrimSpace(line)

			if noHandle {
				resultChan <- line
				continue
			}

			if !strings.HasPrefix(line, "data:") {
				continue
			}
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if data == "" || data == "[DONE]" {
				continue
			}

			// 解析 SSE 数据段
			respObj, err := c.parseResponse([]byte(data))
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- respObj
		}
	}()

	return resultChan, errorChan
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
