package config

type HTTPAction struct {
	ValidStatus []uint            `yaml:"valid-status" validate:"required,min=1,max=20,dive,max=1000"`
	Target      string            `yaml:"target" validate:"required,max=255,min=1"`
	Method      string            `yaml:"method" validate:"required,oneof=GET POST PUT DELETE HEAD"`
	Port        uint              `yaml:"port" validate:"required,max=65535,min=1"`
	Redirect    bool              `yaml:"redirect"`
	Timeout     string            `yaml:"timeout" validate:"required"`
	Body        string            `yaml:"body"`
	Query       map[string]string `json:"query,omitempty"`
	Headers     map[string]string `yaml:"headers,omitempty" validate:"max=5"`
	Protocol    string            `yaml:"protocol" validate:"oneof=http https"`
	Path        string            `yaml:"path,omitempty"`
	Key         string            `yaml:"key,omitempty"`
	Cert        string            `yaml:"cert,omitempty"`
	Cacert      string            `yaml:"cacert,omitempty"`
	ServerName  string            `yaml:"server-name"`
	Insecure    bool              `yaml:"insecure"`
	Extractors  HTTPExtractors    `yaml:"extractors"`
}

type HTTPExtractors struct {
	Headers  map[string]string `yaml:"headers"`
	BodyJSON map[string][]any  `yaml:"body-json"`
	Body     string            `yaml:"body"`
}

type Action struct {
	Name         string     `yaml:"name"`
	Description  string     `yaml:"description"`
	HTTP         HTTPAction `yaml:"http"`
	Log          LogAction  `yaml:"log"`
	Transformers []string   `yaml:"transformers"`
}

type LogAction struct {
	Message string `yaml:"message"`
}

type Test struct {
	Variables    map[string]any         `yaml:"variables"`
	Actions      []Action               `yaml:"actions" validate:"required,min=1"`
	Transformers map[string]Transformer `yaml:"transformers"`
}

type Transformer struct {
	Exoscale TransformerExoscaleConfig `yaml:"exoscale"`
}

type TransformerExoscaleConfig struct {
	APIKey    string `yaml:"api-key"`
	APISecret string `yaml:"api-secret"`
}
