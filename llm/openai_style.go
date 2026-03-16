package llm

import "errors"

// ======== Request Body 定义 ========

type OpenAIThinking struct {
	Type string `json:"type" yaml:"type"`
}

type OpenAIResponseFormat struct {
	Type string `json:"type" yaml:"type"`
}

type OpenAIChatRequestBody struct {
	Model            string                `json:"model"`    // 必选
	Messages         []MessageItem         `json:"messages"` // 必选
	Thinking         *OpenAIThinking       `json:"thinking,omitempty"`
	FrequencyPenalty *float64              `json:"frequency_penalty,omitempty"`
	MaxTokens        *int                  `json:"max_tokens,omitempty"`
	PresencePenalty  *float64              `json:"presence_penalty,omitempty"`
	ResponseFormat   *OpenAIResponseFormat `json:"response_format,omitempty"`
	Stop             interface{}           `json:"stop,omitempty"`
	Stream           *bool                 `json:"stream,omitempty"`
	StreamOptions    interface{}           `json:"stream_options,omitempty"`
	Temperature      *float64              `json:"temperature,omitempty"`
	TopP             *float64              `json:"top_p,omitempty"`
	Tools            []LLMTool             `json:"tools,omitempty"`
	ToolChoice       interface{}           `json:"tool_choice,omitempty"`
	Logprobs         *bool                 `json:"logprobs,omitempty"`
	TopLogprobs      *int                  `json:"top_logprobs,omitempty"`
}

func (r *OpenAIChatRequestBody) Validate() (bool, error) {
	if r.Model == "" {
		return false, errors.New("model is required")
	}
	if len(r.Messages) == 0 {
		return false, errors.New("at least one message is required")
	}
	// 校验 frequency_penalty [-2,2]
	if r.FrequencyPenalty != nil {
		if *r.FrequencyPenalty < -2 || *r.FrequencyPenalty > 2 {
			return false, errors.New("frequency_penalty must be between -2 and 2")
		}
	}
	// 校验 presence_penalty [-2,2]
	if r.PresencePenalty != nil {
		if *r.PresencePenalty < -2 || *r.PresencePenalty > 2 {
			return false, errors.New("presence_penalty must be between -2 and 2")
		}
	}
	// 校验 temperature [0,2]
	if r.Temperature != nil {
		if *r.Temperature < 0 || *r.Temperature > 2 {
			return false, errors.New("temperature must be between 0 and 2")
		}
	}
	// 校验 top_p [0,1]
	if r.TopP != nil {
		if *r.TopP < 0 || *r.TopP > 1 {
			return false, errors.New("top_p must be between 0 and 1")
		}
	}
	// 校验 top_logprobs [0,20]
	if r.TopLogprobs != nil {
		if *r.TopLogprobs < 0 || *r.TopLogprobs > 20 {
			return false, errors.New("top_logprobs must be between 0 and 20")
		}
	}

	if r.StreamOptions != nil {
		// 判断是否开启了 stream，如果没有开启，则 stream_options 无效
		if r.Stream == nil || !*r.Stream {
			return false, errors.New("stream_options is invalid when stream is not enabled")
		}
	}
	return true, nil
}
