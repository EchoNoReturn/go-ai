# go-ai

不使用现成框架，手撸 AI 和 AGENT 框架代码。

## 功能特性

- 支持多种 LLM 提供商：OpenAI、Anthropic、DeepSeek
- 会话管理：创建和管理 LLM 对话会话
- 消息类型：System、User、Assistant、Tool 消息
- 请求参数校验：temperature、top_p、frequency_penalty 等
- 配置管理：支持从 JSON 文件加载 API 配置

## 项目结构

```
go-ai/
├── main.go              # 入口文件
├── configs/             # 配置加载
│   ├── json_config.go
│   └── types.go
├── llm/                 # LLM 核心模块
│   ├── base_types.go    # 基础类型定义
│   ├── llm_session.go  # 会话管理
│   ├── openai_style.go # OpenAI 请求
│   └── anthropic_style.go
└── config.json         # 配置文件
```

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/EchoNoReturn/go-ai.git
cd go-ai
```

### 2. 配置 API Key

复制配置示例文件并填入你的 API Key：

```bash
cp config.example.json config.json
```

编辑 `config.json`：

```json
{
  "providers": {
    "deepseek": {
      "name": "deepseek",
      "endpoint": "https://api.deepseek.com",
      "apiKey": "your-api-key"
    }
  }
}
```

### 3. 运行

```bash
go run .
```

## 开发指南

### 代码格式化

```bash
go fmt ./...
go vet ./...
```

### 运行测试

```bash
go test ./...
```

## 目标

- 实现 MCP (Model Context Protocol) 协议
- 支持 Skills 智能体功能
- 跨平台打包支持

## License

MIT
