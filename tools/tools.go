package tools

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"google.golang.org/genai"
)

func GenerateSchema[T any]() genai.Schema {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T

	schema := reflector.Reflect(v)

	jsonSchema, err := json.Marshal(schema)
	if err != nil {
		panic(err)
	}

	geminiSchema := genai.Schema{}
	err = json.Unmarshal(jsonSchema, &geminiSchema)
	if err != nil {
		panic(err)
	}

	return geminiSchema
}

type ToolDefinition struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	InputSchema genai.Schema `json:"input_schema"`
	Function    func(input map[string]any) (string, error)
}

type ToolDefinitions []ToolDefinition

func (t ToolDefinitions) ToGeminiFunction() []*genai.FunctionDeclaration {
	var functions []*genai.FunctionDeclaration
	for _, tool := range t {
		function := &genai.FunctionDeclaration{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  &tool.InputSchema,
		}
		functions = append(functions, function)
	}
	return functions
}

func TryReadMap(input map[string]any, key string) (string, error) {
	value, found := input[key]
	if !found {
		return "", fmt.Errorf("missing key: %s", key)
	}
	valueStr, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("key %s must be a string", key)
	}
	return valueStr, nil
}
