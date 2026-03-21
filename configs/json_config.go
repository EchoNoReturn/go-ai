package configs

import (
	"encoding/json"
	"fmt"
	"os"
)

// 优化后的配置加载方法，支持 config.json 的 RootConfig 结构体
func LoadConfigFromFile(path string) (*RootConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("配置文件 %s 不存在\n", path)
		return nil, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败： %w", err)
	}
	var config RootConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败： %w", err)
	}
	return &config, nil
}

func UpdateConfig(data *RootConfig, path string) {
	// 判断文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("配置文件 %s 不存在\n", path)
		println("Will create config file: config.json")
		file, err := os.Create(path)
		if err != nil {
			fmt.Printf("创建配置文件失败: %v\n", err)
			return
		}
		defer file.Close()
		return
	}
	dataBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("序列化配置数据失败: %v\n", err)
		return
	}
	err = os.WriteFile(path, dataBytes, 0644)
	if err != nil {
		fmt.Printf("写入配置文件失败: %v\n", err)
		return
	}
}