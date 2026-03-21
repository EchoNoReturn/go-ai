package configs

type ProviderConfig struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	ApiKey   string `json:"apiKey"`
}

type ApiReference struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

type ApiReferences struct {
	Chat struct {
		OpenAI ApiReference `json:"openai"`
	} `json:"chat"`
}

type RootConfig struct {
	Providers       map[string]ProviderConfig `json:"providers"`
	CurrentProvider string                    `json:"currentProvider"`
	ApiReferences   ApiReferences             `json:"apiReferences"`
}

func (c *RootConfig) GetProvider(name string) (*ProviderConfig, bool) {
	provider, exists := c.Providers[name]
	if !exists {
		return nil, false
	}
	return &provider, true
}

func (c *RootConfig) GetCurrentProvider() (*ProviderConfig, bool) {
	if c == nil {
		return nil, false
	}
	return c.GetProvider(c.CurrentProvider)
}
