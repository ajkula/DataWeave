package apigen

import (
	"db_meta/api"
	"db_meta/dbstructs"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

func GenerateOpenAPI(tables []*dbstructs.TableMetadata, config *api.APIConfig) ([]byte, error) {
	openAPI := api.OpenAPI{
		OpenAPI: "3.0.0",
		Info: api.Info{
			Title:       "Generated API",
			Description: "API generated from database schema",
			Version:     "1.0.0",
		},
		Paths: make(map[string]api.PathItem),
		Components: api.Components{
			Schemas: make(map[string]api.Schema),
		},
	}

	for _, table := range tables {
		if config == nil || (*config)[strings.ToLower(table.TableName)] != nil {
			generatePathsForTable(&openAPI, table, config)
			generateSchemaForTable(&openAPI, table)
		}
	}

	return yaml.Marshal(openAPI)
}

func generatePathsForTable(openAPI *api.OpenAPI, table *dbstructs.TableMetadata, config *api.APIConfig) {
	basePath := fmt.Sprintf("/%s", strings.ToLower(table.TableName))
	var tableConfig api.TableConfig
	if config != nil {
		tableConfig = (*config)[strings.ToLower(table.TableName)]
	}

	// List/creation
	if config == nil || (tableConfig != nil && tableConfig[basePath]["GET"].Included) {
		methodConfig := api.MethodConfig{}
		if tableConfig != nil {
			methodConfig = tableConfig[basePath]["GET"]
		}
		getOperation := &api.Operation{
			Summary:     "List " + table.TableName,
			OperationID: generateUniqueOperationID(table.TableName, "list"),
			Parameters:  generateQueryParameters(table, config),
			Responses:   generateStandardResponses(table, true, methodConfig),
		}
		addRequestHeaders(getOperation, methodConfig)

		methodConfig = api.MethodConfig{}
		if tableConfig != nil {
			methodConfig = tableConfig[basePath]["POST"]
		}
		postOperation := &api.Operation{
			Summary:     "Create a new " + table.TableName,
			OperationID: generateUniqueOperationID(table.TableName, "create"),
			RequestBody: &api.RequestBody{
				Required: true,
				Content: map[string]api.MediaType{
					"application/json": {
						Schema: api.Schema{
							Ref: "#/components/schemas/" + table.TableName,
						},
					},
				},
			},
			Responses: generateStandardResponses(table, false, methodConfig),
		}
		addRequestHeaders(postOperation, methodConfig)

		openAPI.Paths[basePath] = api.PathItem{
			Get:  getOperation,
			Post: postOperation,
		}

		// specific elements
		itemPath := fmt.Sprintf("%s/{id}", basePath)
		idParam := api.Parameter{
			Name:        "id",
			In:          "path",
			Required:    true,
			Description: fmt.Sprintf("ID of the %s", table.TableName),
			Schema:      &api.Schema{Type: "string"},
		}

		methodConfig = api.MethodConfig{}
		if tableConfig != nil {
			methodConfig = tableConfig[itemPath]["GET"]
		}
		getSpecificOperation := &api.Operation{
			Summary:     "Get a specific " + table.TableName,
			OperationID: generateUniqueOperationID(table.TableName, "get"),
			Parameters:  []api.Parameter{idParam},
			Responses:   generateStandardResponses(table, false, methodConfig),
		}
		addRequestHeaders(getSpecificOperation, methodConfig)

		methodConfig = api.MethodConfig{}
		if tableConfig != nil {
			methodConfig = tableConfig[itemPath]["PUT"]
		}
		putOperation := &api.Operation{
			Summary:     "Update a " + table.TableName,
			OperationID: generateUniqueOperationID(table.TableName, "update"),
			Parameters:  []api.Parameter{idParam},
			RequestBody: &api.RequestBody{
				Required: true,
				Content: map[string]api.MediaType{
					"application/json": {
						Schema: api.Schema{
							Ref: "#/components/schemas/" + table.TableName,
						},
					},
				},
			},
			Responses: generateStandardResponses(table, false, methodConfig),
		}
		addRequestHeaders(putOperation, methodConfig)

		methodConfig = api.MethodConfig{}
		if tableConfig != nil {
			methodConfig = tableConfig[itemPath]["DELETE"]
		}
		deleteOperation := &api.Operation{
			Summary:     "Delete a " + table.TableName,
			OperationID: generateUniqueOperationID(table.TableName, "delete"),
			Parameters:  []api.Parameter{idParam},
			Responses:   generateStandardResponses(table, false, methodConfig),
		}
		addRequestHeaders(deleteOperation, methodConfig)

		openAPI.Paths[itemPath] = api.PathItem{
			Get:    getSpecificOperation,
			Put:    putOperation,
			Delete: deleteOperation,
		}

		// relations
		for _, relation := range table.Relationships {
			relatedPath := fmt.Sprintf("%s/{id}/%s", basePath, strings.ToLower(relation.RelatedTableName))

			methodConfig = api.MethodConfig{}
			if tableConfig != nil {
				methodConfig = tableConfig[relatedPath]["GET"]
			}
			relatedOperation := &api.Operation{
				Summary:     fmt.Sprintf("List %s for %s", relation.RelatedTableName, table.TableName),
				OperationID: generateUniqueOperationID(table.TableName, "listRelated"+relation.RelatedTableName),
				Parameters:  append([]api.Parameter{idParam}, generateQueryParameters(table, config)...),
				Responses:   generateStandardResponses(table, true, methodConfig),
			}
			addRequestHeaders(relatedOperation, methodConfig)

			openAPI.Paths[relatedPath] = api.PathItem{
				Get: relatedOperation,
			}
		}
	}
}

func addRequestHeaders(operation *api.Operation, methodConfig api.MethodConfig) {
	if operation.Parameters == nil {
		operation.Parameters = []api.Parameter{}
	}
	for header, included := range methodConfig.RequestHeaders {
		if included {
			operation.Parameters = append(operation.Parameters, api.Parameter{
				Name:        header,
				In:          "header",
				Description: fmt.Sprintf("%s header", header),
				Schema:      &api.Schema{Type: "string"},
				Required:    false,
			})
		}
	}
}

func generateUniqueOperationID(tableName, operation string) string {
	return fmt.Sprintf("%s_%s_%s", operation, strings.ToLower(tableName), uuid.New().String())
}

func generateQueryParameters(table *dbstructs.TableMetadata, config *api.APIConfig) []api.Parameter {
	params := []api.Parameter{
		{
			Name:   "page",
			In:     "query",
			Schema: &api.Schema{Type: "integer"},
		},
		{
			Name:   "limit",
			In:     "query",
			Schema: &api.Schema{Type: "integer"},
		},
	}

	tableConfig := (*config)[strings.ToLower(table.TableName)]
	basePath := fmt.Sprintf("/%s", strings.ToLower(table.TableName))

	for _, column := range table.Columns {
		if column.Unique {
			if config == nil || tableConfig[basePath]["GET"].Filters[column.ColumnName] {
				params = append(params, api.Parameter{
					Name:        column.ColumnName,
					In:          "query",
					Description: fmt.Sprintf("Filter by %s", column.ColumnName),
					Schema:      &api.Schema{Type: mapSQLTypeToJSONType(column.DataType)},
				})
			}
		}
	}

	return params
}

func generateStandardResponses(table *dbstructs.TableMetadata, isArray bool, methodConfig api.MethodConfig) map[string]api.Response {
	var successSchema api.Schema
	var successExample *api.Example

	if isArray {
		successSchema = api.Schema{
			Type: "object",
			Properties: map[string]api.Schema{
				"data": {
					Type: "array",
					Items: &api.Schema{
						Ref: "#/components/schemas/" + table.TableName,
					},
				},
				"pagination": {
					Type: "object",
					Properties: map[string]api.Schema{
						"total": {Type: "integer"},
						"pages": {Type: "integer"},
						"page":  {Type: "integer"},
						"limit": {Type: "integer"},
					},
				},
			},
		}
		successExample = &api.Example{
			Value: map[string]interface{}{
				"data": []interface{}{generateExampleForTable(table)},
				"pagination": map[string]int{
					"total": 100, "pages": 10, "page": 1, "limit": 10,
				},
			},
		}
	} else {
		successSchema = api.Schema{Ref: "#/components/schemas/" + table.TableName}
		successExample = &api.Example{Value: generateExampleForTable(table)}
	}

	responses := map[string]api.Response{
		"200": {
			Description: "Successful response",
			Content: map[string]api.MediaType{
				"application/json": {Schema: successSchema, Example: successExample},
			},
			Headers: make(map[string]api.Header),
		},
		"201": {
			Description: "Created successfully",
			Content: map[string]api.MediaType{
				"application/json": {
					Schema:  api.Schema{Ref: "#/components/schemas/" + table.TableName},
					Example: &api.Example{Value: generateExampleForTable(table)},
				},
			},
			Headers: make(map[string]api.Header),
		},
		"400": {
			Description: "Bad Request",
			Content: map[string]api.MediaType{
				"application/json": {
					Schema: api.Schema{
						Type:       "object",
						Properties: map[string]api.Schema{"error": {Type: "string"}},
					},
					Example: &api.Example{Value: map[string]string{"error": "Invalid input"}},
				},
			},
			Headers: make(map[string]api.Header),
		},
		"404": {Description: "Not Found"},
		"500": {Description: "Internal Server Error"},
	}

	if methodConfig.ResponseHeaders != nil {
		for header, included := range methodConfig.ResponseHeaders {
			if included {
				responses["200"].Headers[header] = api.Header{
					Description: fmt.Sprintf("%s header", header),
					Schema:      &api.Schema{Type: "string"},
				}
			}
		}
	}

	return responses
}

func generateSchemaForTable(openAPI *api.OpenAPI, table *dbstructs.TableMetadata) {
	properties := make(map[string]api.Schema)
	required := []string{}

	for _, column := range table.Columns {
		properties[column.ColumnName] = api.Schema{Type: mapSQLTypeToJSONType(column.DataType)}
		if column.NotNull {
			required = append(required, column.ColumnName)
		}
	}

	openAPI.Components.Schemas[table.TableName] = api.Schema{
		Type:       "object",
		Properties: properties,
		Required:   required,
	}
}

func generateExampleForTable(table *dbstructs.TableMetadata) map[string]interface{} {
	example := make(map[string]interface{})
	for _, column := range table.Columns {
		example[column.ColumnName] = generateExampleValue(*column)
	}
	return example
}

func generateExampleValue(column dbstructs.Column) interface{} {
	switch mapSQLTypeToJSONType(column.DataType) {
	case "integer":
		return 42
	case "number":
		return 3.14
	case "boolean":
		return true
	case "string":
		if strings.Contains(strings.ToLower(column.DataType), "date") ||
			strings.Contains(strings.ToLower(column.DataType), "time") {
			return time.Now().Format(time.RFC3339)
		}
		return fmt.Sprintf("Example %s", column.ColumnName)
	default:
		return "Example value"
	}
}

func mapSQLTypeToJSONType(sqlType string) string {
	switch sqlType {
	case "integer", "bigint", "smallint":
		return "integer"
	case "numeric", "decimal", "real", "double precision":
		return "number"
	case "boolean":
		return "boolean"
	case "date", "timestamp":
		return "string"
	default:
		return "string"
	}
}
