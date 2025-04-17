package tools

import "os"

var ReadFileDefinition = ToolDefinition{
	Name:        "read_file",
	Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
	InputSchema: ReadFileInputSchema,
	Function:    ReadFile,
}

type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory."`
}

var ReadFileInputSchema = GenerateSchema[ReadFileInput]()

func ReadFile(input map[string]any) (string, error) {
	pathStr, err := TryReadMap(input, "path")
	if err != nil {
		return "", err
	}
	content, err := os.ReadFile(pathStr)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
