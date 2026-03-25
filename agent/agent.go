package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-ai/configs"
	"go-ai/llm"
)

type RunTaskOptions struct {
	EnvConfig configs.RootConfig
	// StreamCallback 是一个可选的回调函数，用于把流式输出结果额外提供给外部使用
	StreamCallback func(result string)
}

func initMemory(session *llm.LLMSession) {
	// 这里可以根据需要添加系统提示词
	session.AppendSystemMessage("你是一个智能助手，会始终使用中文回答问题", "")
}

func RunTask(task string, options RunTaskOptions) error {
	// 加载配置
	APPLICATION_CONFIG := options.EnvConfig

	// 创建一个LLMSession示例
	providerConfig, exists := APPLICATION_CONFIG.GetCurrentProvider()
	if !exists {
		err := errors.New("Current provider not found in configuration.")
		return err
	}
	// 初始化 LLM 会话
	session, err := llm.NewSession(providerConfig.ApiKey, providerConfig.Endpoint)
	if err != nil {
		// 处理错误
		return err
	}
	modelName, err := APPLICATION_CONFIG.GetCurrentModel()
	if err != nil {
		// 处理错误
		return err
	}
	session.Model = modelName

	// 初始化记忆系统
	initMemory(session)

	session.AppendUserMessage(task, "")

	// 运行会话并获取结果
	resultChan, errorChan := session.RunChatStream(true)
	handleByType := (func(result any) (string, error) {
		var content string
		if session.Type == llm.OpenAI {
			// 处理 OpenAI 风格的结果
			println("origin: ", toJsonString(result))
			result, ok := result.(*llm.OpenAIChatResponseBody)
			if !ok {
				return "", fmt.Errorf("Unexpected result type")
			}
			if len(result.Choices) > 0 {
				content = result.Choices[0].Message.Content
			} else {
				content = ""
			}
			
		} else {
			// 处理Anthropic风格的结果
			result, ok := result.(*llm.AnthropicChatResponse)
			if !ok {
				return "", fmt.Errorf("Unexpected result type")
			}
			content = result.Content.Text
		}
		return content, nil
	})

	for {
		select {
		case result, ok := <-resultChan:
			if !ok {
				resultChan = nil
			} else {
				// 根据 resultType 区分风格处理处理结果
				println("-->", toJsonString(result))
				content, err := handleByType(result)
				if err != nil {
					println("Error:", err.Error())
					continue
				}
				println("Streamed Result:", content)
			}
		case err, ok := <-errorChan:
			if !ok {
				errorChan = nil
			} else {
				// 处理错误
				println("Error:", err.Error())
				return err
			}
		}

		if resultChan == nil && errorChan == nil {
			break
		}
	}

	return nil
}

func toJsonString(result interface{}) string {
	bytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Sprintf("Error converting to JSON: %v", err)
	}
	return string(bytes)
}
