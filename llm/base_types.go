package llm

import (
	"encoding/json"
)

// ======== ChatRequest 接口 ========
type ChatRequest interface {
	Validate() (bool, error)
}

// ======== LLMType 定义 ========
type LLMType uint8

const (
	OpenAI LLMType = iota
	Anthropic
)

func (t LLMType) String() string {
	switch t {
	case OpenAI:
		return "openai"
	case Anthropic:
		return "anthropic"
	default:
		return "unknown"
	}
}

func StrToLLMType(s string) LLMType {
	switch s {
	case "openai":
		return OpenAI
	case "anthropic":
		return Anthropic
	default:
		return 0
	}
}

// ======== 消息定义 ========
type LLMMessageType uint8

const (
	MSystem LLMMessageType = iota
	MUser
	MAssistant
	MTool
)

func (t LLMMessageType) String() string {
	switch t {
	case MSystem:
		return "system"
	case MUser:
		return "user"
	case MAssistant:
		return "assistant"
	case MTool:
		return "tool"
	default:
		return ""
	}
}

type LLMMessage interface {
	ToJsonString() string
}

type MessageItem struct {
	Role    string `json:"role" yaml:"role"`
	Content string `json:"content" yaml:"content"`
}

func toJsonString(m any) string {
	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(data)
}

// ======== SystemMessage ========

type SystemMessage struct {
	MessageItem
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

func (m *SystemMessage) ToJsonString() string {
	return toJsonString(m)
}

// ======== UserMessage ========

type UserMessage struct {
	MessageItem
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

func (m *UserMessage) ToJsonString() string {
	return toJsonString(m)
}

// ======== AssistantMessage ========

type AssistantMessage struct {
	MessageItem
	Name             string `json:"name,omitempty" yaml:"name,omitempty"`
	Prefix           *bool  `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	ReasoningContent string `json:"reasoning_content,omitempty" yaml:"reasoning_content,omitempty"`
}

func (m *AssistantMessage) ToJsonString() string {
	return toJsonString(m)
}

// ======== ToolMessage ========

type ToolMessage struct {
	MessageItem
	ToolCallId string `json:"tool_call_id,omitempty" yaml:"tool_call_id,omitempty"`
}

func (m *ToolMessage) ToJsonString() string {
	return toJsonString(m)
}

// ======= 消息工厂函数 ========

func CreateMessage(role LLMMessageType, content string) LLMMessage {
	switch role {
	case MSystem:
		return &SystemMessage{
			MessageItem: MessageItem{
				Role:    role.String(),
				Content: content,
			},
		}
	case MUser:
		return &UserMessage{
			MessageItem: MessageItem{
				Role:    role.String(),
				Content: content,
			},
		}
	case MAssistant:
		return &AssistantMessage{
			MessageItem: MessageItem{
				Role:    role.String(),
				Content: content,
			},
		}
	case MTool:
		return &ToolMessage{
			MessageItem: MessageItem{
				Role:    role.String(),
				Content: content,
			},
		}
	default:
		return nil
	}
}

// ======= LLM工具定义 ========

type LLMToolFunction struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
}

type LLMTool struct {
	Type     string          `json:"type" yaml:"type"`
	Function LLMToolFunction `json:"function" yaml:"function"`
}
