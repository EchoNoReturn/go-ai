package main

import (
	// "encoding/json"
	"go-ai/agent"
	"go-ai/configs"
	// "go-ai/llm"
)

const _config_path = "config.json"

var APPLICATION_CONFIG *configs.RootConfig

func main() {
	APPLICATION_CONFIG, err := configs.LoadConfigFromFile(_config_path)
	if err != nil {
		panic(err)
	}

	// 打印加载的配置
	if APPLICATION_CONFIG == nil {
		panic("No configuration loaded.")
	}

	agent.RunTask("你好，你是谁啊？", agent.RunTaskOptions{
		EnvConfig: *APPLICATION_CONFIG,
		StreamCallback: func(result string) {
			println("回答:", result)
		},
	})
	// 创建一个LLMSession示例
	// providerConfig, exists := APPLICATION_CONFIG.GetCurrentProvider()
	// if !exists {
	// 	panic("Current provider not found in configuration.")
	// }
	// session, err := llm.NewSession(providerConfig.ApiKey, providerConfig.Endpoint)
	// if err != nil {
	// 	panic(err)
	// }

	// session.Model = APPLICATION_CONFIG.CurrentModel

	// session.Tools = []llm.LLMTool{
	// 	{
	// 		Type: "function",
	// 		Function: llm.LLMToolFunction{
	// 			Name:        "get_current_weather",
	// 			Description: "获取指定城市当前的天气信息",
	// 			Parameters: map[string]interface{}{
	// 				"type": "object",
	// 				"properties": map[string]interface{}{
	// 					"city": map[string]interface{}{
	// 						"type":        "string",
	// 						"description": "要查询天气的城市名称",
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// session.AppendSystemMessage("你是一个智能助手", "")
	// session.AppendUserMessage("北京现在是什么天气？", "")

	// var num = 0
	// resultChan, errorChan := session.RunChatStream(true)
	// for resultChan != nil || errorChan != nil {
	// 	select {
	// 	case result, ok := <-resultChan:
	// 		if !ok {
	// 			resultChan = nil
	// 			continue
	// 		}
	// 		responseBytes, err := json.Marshal(result)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		println("LLM Response:", num, string(responseBytes))
	// 		num++
	// 	case err, ok := <-errorChan:
	// 		if !ok {
	// 			errorChan = nil
	// 			continue
	// 		}
	// 		panic(err)
	// 	}
	// }
}
