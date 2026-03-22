package server

// Config is a struct containing server configuration values.
type Config struct {
	Title         string `env:"TITLE"          json:"title"`
	Version       string `env:"VERSION"        json:"version"`
	OpenAPIPath   string `env:"OPEN_API_PATH"  json:"openApiPath"`
	DocsPath      string `env:"DOCS_PATH"      json:"docsPath"`
	SchemasPath   string `env:"SCHEMAS_PATH"   json:"schemasPath"`
	DefaultFormat string `env:"DEFAULT_FORMAT" json:"defaultFormat"`
}
