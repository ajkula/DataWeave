package api

import "db_meta/dbstructs"

// EndpointDefinition represents the structure for defining an API endpoint.
type EndpointDefinition struct {
	Path            string
	GET             *MethodDefinition // Pointers allow for a nil if the method is not defined
	POST            *MethodDefinition
	PUT             *MethodDefinition
	PATCH           *MethodDefinition
	DELETE          *MethodDefinition
	Parameters      []ParameterDefinition
	Responses       map[int]ResponseDefinition // A map with the status code as the key
	Tags            []string
	Security        []SecurityRequirement
	Description     string
	OperationId     string
	RequestBody     *dbstructs.TableMetadata
	QueryParameters []dbstructs.Column
}

// MethodDefinition defines the details of a HTTP method for an endpoint.
type MethodDefinition struct {
	Enabled         bool
	Description     string
	OperationId     string
	RequestBody     *RequestBodyDefinition
	QueryParameters []ParameterDefinition
	Responses       map[int]ResponseDefinition
}

// RequestBodyDefinition represents the expected structure of a request body.
type RequestBodyDefinition struct {
	Description string
	Required    bool
	Content     map[string]MediaTypeDefinition // e.g., "application/json", "application/xml"
}

// MediaTypeDefinition specifies the media type and example for a request or response body.
type MediaTypeDefinition struct {
	Schema  interface{}
	Example interface{}
}

// ParameterDefinition defines a single operation parameter.
type ParameterDefinition struct {
	Name        string
	In          string // "query", "header", "path", "cookie"
	Description string
	Required    bool
	Schema      interface{}
	Example     interface{}
}

// ResponseDefinition represents the structure of a response from an API endpoint.
type ResponseDefinition struct {
	Description string
	Headers     map[string]HeaderDefinition
	Content     map[string]MediaTypeDefinition // "application/json", "application/whatever"
}

// HeaderDefinition defines the expected structure of a header in a request or response.
type HeaderDefinition struct {
	Description string
	Required    bool
	Schema      interface{}
}

// SecurityRequirement specifies the security scheme and its requirements.
type SecurityRequirement struct {
	Type         string // "apiKey", "http", "oauth2", "openIdConnect"
	Scheme       string
	BearerFormat string
	In           string // "query", "header"
	Name         string // Name of the header or query parameter
}

// PaginationInfo is used to store pagination information in API responses.
type PaginationInfo struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalCount int `json:"totalCount"`
	TotalPages int `json:"totalPages"`
}

// PagedResponse is a wrapper around table data with pagination information.
type PagedResponse struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}
