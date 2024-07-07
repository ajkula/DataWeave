package api

type OpenAPI struct {
	OpenAPI    string                `yaml:"openapi"`
	Info       Info                  `yaml:"info"`
	Paths      map[string]PathItem   `yaml:"paths"`
	Components Components            `yaml:"components"`
	Security   []SecurityRequirement `yaml:"security,omitempty"`
}

type Info struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description,omitempty"`
	Version     string `yaml:"version"`
}

type PathItem struct {
	Get        *Operation  `yaml:"get,omitempty"`
	Post       *Operation  `yaml:"post,omitempty"`
	Put        *Operation  `yaml:"put,omitempty"`
	Patch      *Operation  `yaml:"patch,omitempty"`
	Delete     *Operation  `yaml:"delete,omitempty"`
	Parameters []Parameter `yaml:"parameters,omitempty"`
}

type Operation struct {
	Summary     string                `yaml:"summary"`
	Description string                `yaml:"description,omitempty"`
	OperationID string                `yaml:"operationId"`
	Tags        []string              `yaml:"tags,omitempty"`
	Parameters  []Parameter           `yaml:"parameters,omitempty"`
	RequestBody *RequestBody          `yaml:"requestBody,omitempty"`
	Responses   map[string]Response   `yaml:"responses"`
	Security    []SecurityRequirement `yaml:"security,omitempty"`
}

type Parameter struct {
	Name        string  `yaml:"name"`
	In          string  `yaml:"in"`
	Description string  `yaml:"description,omitempty"`
	Required    bool    `yaml:"required,omitempty"`
	Schema      *Schema `yaml:"schema"`
}

type RequestBody struct {
	Description string               `yaml:"description,omitempty"`
	Required    bool                 `yaml:"required,omitempty"`
	Content     map[string]MediaType `yaml:"content"`
}

type Response struct {
	Description string               `yaml:"description"`
	Headers     map[string]Header    `yaml:"headers,omitempty"`
	Content     map[string]MediaType `yaml:"content,omitempty"`
}

type MediaType struct {
	Schema  Schema   `yaml:"schema"`
	Example *Example `yaml:"example,omitempty"`
}

type Schema struct {
	Type       string            `yaml:"type,omitempty"`
	Properties map[string]Schema `yaml:"properties,omitempty"`
	Items      *Schema           `yaml:"items,omitempty"`
	Ref        string            `yaml:"$ref,omitempty"`
	Required   []string          `yaml:"required,omitempty"`
}

type Header struct {
	Description string  `yaml:"description,omitempty"`
	Schema      *Schema `yaml:"schema"`
	Required    bool    `yaml:"required,omitempty"`
}

type Components struct {
	Schemas         map[string]Schema         `yaml:"schemas,omitempty"`
	SecuritySchemes map[string]SecurityScheme `yaml:"securitySchemes,omitempty"`
}

type SecurityScheme struct {
	Type         string `yaml:"type"`
	Description  string `yaml:"description,omitempty"`
	Name         string `yaml:"name,omitempty"`
	In           string `yaml:"in,omitempty"`
	Scheme       string `yaml:"scheme,omitempty"`
	BearerFormat string `yaml:"bearerFormat,omitempty"`
}

type SecurityRequirement map[string][]string

type Example struct {
	Value interface{} `yaml:"value"`
}

// params structs
type APIConfig map[string]TableConfig

type TableConfig map[string]EndpointConfig

type EndpointConfig map[string]MethodConfig

type MethodConfig struct {
	Included        bool            `json:"included"`
	Security        string          `json:"security"`
	Filters         map[string]bool `json:"filters,omitempty"`
	RequestHeaders  map[string]bool `json:"requestHeaders,omitempty"`
	ResponseHeaders map[string]bool `json:"responseHeaders,omitempty"`
}
