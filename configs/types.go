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
	Providers     map[string]ProviderConfig `json:"providers"`
	ApiReferences ApiReferences             `json:"apiRefrances"`
}
