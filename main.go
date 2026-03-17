package main

import (
	"encoding/json"
	"go-ai/configs"
	"go-ai/llm"
)

func main() {
	const configPath = "config.json"
	config, err := configs.LoadConfigFromFile(configPath)
	if err != nil {
		panic(err)
	}

	// 打印加载的配置
	if config != nil {
		for name, provider := range config.Providers {
			println("Provider Name:", name)
			println("Endpoint:", provider.Endpoint)
			println("API Key:", provider.ApiKey)
			println()
		}
	} else {
		println("No configuration loaded.")
	}

	// 创建一个LLMSession示例
	deepseekConfig, exists := config.Providers["deepseek"]
	if !exists {
		println("Provider 'deepseek' not found in configuration.")
		return
	}
	session, err := llm.NewSession(deepseekConfig.ApiKey, deepseekConfig.Endpoint)
	if err != nil {
		panic(err)
	}

	session.Model = "deepseek-chat"

	session.Tools = []llm.LLMTool{
		{
			Type: "function",
			Function: llm.LLMToolFunction{
				Name:        "get_current_weather",
				Description: "获取指定城市当前的天气信息",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"city": map[string]interface{}{
							"type":        "string",
							"description": "要查询天气的城市名称",
						},
					},
				},
			},
		},
	}

	session.AppendSystemMessage("你是一个智能助手", "")
	session.AppendUserMessage("北京现在是什么天气？", "")

	var num = 0
	resultChan, errorChan := session.RunChatStream()
	for resultChan != nil || errorChan != nil  {
		select {
		case result, ok := <-resultChan:
			if !ok {
				resultChan = nil
				continue
			}
			responseBytes, err := json.Marshal(result)
			if err != nil {
				panic(err)
			}
			println("LLM Response:", num, string(responseBytes))
			num++
		case err, ok := <-errorChan:
			if !ok {
				errorChan = nil
				continue
			}
			panic(err)
		}
	}
}
