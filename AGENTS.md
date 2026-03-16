# AGENTS.md - Agent Guidelines for go-ai

## Project Overview

This is a Go project for interacting with LLM providers (OpenAI, Anthropic, DeepSeek). It provides a session-based API for chat completions, model listing, and message handling.

## Build & Development Commands

### Running the Project
```bash
# Build the application
go build -o go-ai .

# Run the application
go run .
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run a single test (use -run flag with test name pattern)
go test -v -run TestFunctionName ./...

# Run tests in a specific package
go test -v ./llm/
go test -v ./configs/
```

### Code Quality
```bash
# Format code (run before committing)
go fmt ./...

# Run linter/vet
go vet ./...

# Run all checks (fmt + vet)
go fmt ./... && go vet ./...
```

## Code Style Guidelines

### Naming Conventions
- **Types**: Use `PascalCase` (e.g., `LLMSession`, `OpenAIChatRequestBody`)
- **Functions/Methods**: Use `PascalCase` (e.g., `NewSession`, `RunChat`)
- **Variables/Constants**: Use `camelCase` (e.g., `apiKey`, `maxTokens`)
- **Packages**: Use short, lowercase names (e.g., `llm`, `configs`)
- **Interfaces**: Use descriptive names ending in `er` where applicable (e.g., `ChatRequest`)

### Import Organization
```go
import (
    "fmt"
    "os"

    "go-ai/configs"
    "go-ai/llm"
)
```
Standard library first, then third-party, then project local packages.

### Code Formatting
- Run `go fmt ./...` before committing
- Use Go's standard formatting (no manual alignment)
- Keep lines under 100 characters when reasonable
- Group related imports together

### Types & Structs
- Use pointer types (`*string`, `*int`) for optional parameters
- Use struct tags for JSON/YAML serialization:
  ```go
  type MessageItem struct {
      Role    string `json:"role" yaml:"role"`
      Content string `json:"content" yaml:"content"`
  }
  ```
- Embed structs for composition (see `SystemMessage`, `UserMessage`)
- Use `uint8` for small enums

### Error Handling
- Use `errors.New()` for simple errors
- Use `fmt.Errorf()` with `%w` for wrapped errors
- Return errors from functions rather than panicking (except in main or init)
- Validate required parameters at the start of functions:
  ```go
  if apiKey == "" {
      return nil, errors.New("apiKey is required")
  }
  ```

### Interfaces & Factories
- Define interfaces for behavior (e.g., `ChatRequest`, `LLMMessage`)
- Use factory functions for construction (e.g., `NewSession`, `CreateMessage`)
- Prefer composition over inheritance

### Constants
- Use `iota` for related constant groups
- Group constants with comments:
  ```go
  const (
      OpenAI LLMType = iota
      Anthropic
  )
  ```

### Method Receivers
- Use short receiver names (`c`, `s`, `m`) consistent with type
- Use pointer receivers (`*LLMSession`) when method modifies state
- Use value receivers for read-only methods

### Comments
- Use Chinese comments (consistent with existing codebase)
- Comment exported functions and types
- Use // ======== Section ======== for major sections

### Git Commits
- When asked to create a git commit, first use git tools to view current changes (git status, git diff)
- Write commit messages in Chinese following mainstream conventions (e.g., `feat:`, `fix:`, `refactor:` prefixes)
- Present the proposed commit message to the user for confirmation before executing the commit
- Never amend commits unless explicitly requested and the commit hasn't been pushed

### Testing
- Create `*_test.go` files in same package
- Use table-driven tests for multiple test cases
- Name test functions `TestFunctionName`

### HTTP Client Patterns
- Create `&http.Client{}` for each request or use shared client
- Always close response bodies with `defer`
- Set appropriate headers (Content-Type, Authorization)
- Handle errors at each step

## Project Structure

```
go-ai/
├── main.go              # Entry point
├── configs/             # Configuration loading
│   ├── json_config.go
│   └── types.go
└── llm/                 # LLM session and API
    ├── base_types.go    # Core types, interfaces
    ├── llm_session.go  # Session implementation
    ├── openai_style.go # OpenAI request types
    └── anthropic_style.go
```

## Configuration

- Configuration is loaded from `config.json`
- Example config in `config.example.json`
- Uses `configs.LoadConfigFromFile()` to load