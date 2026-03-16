package main

import (
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

	model_list, err := session.ListModels()
	if err != nil {
		panic(err)
	}
	println("Available Models:", model_list)

	session.Model = "deepseek-chat"

	session.AppendSystemMessage("你是一个智能助手", "")
	session.AppendUserMessage("请介绍一下自己", "")

	result, err := session.RunChat()
	if err != nil {
		panic(err)
	}
	println("LLM Response:", result)
}
