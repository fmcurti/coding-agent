package main

import (
	"bufio"
	"coding-agent/agent"
	"coding-agent/tools"
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

func main() {
	ctx := context.TODO()
	client, _ := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: os.Getenv("GENAI_API_KEY"),
	})
	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true

	}
	availableTools := tools.ToolDefinitions{}
	availableTools = append(availableTools, tools.ReadFileDefinition)
	availableTools = append(availableTools, tools.ListFilesDefinition)
	availableTools = append(availableTools, tools.EditFileDefinition)
	codingAgent := agent.NewAgent(client, getUserMessage, availableTools)
	err := codingAgent.Run(context.TODO())
	if err != nil {
		fmt.Println("Error:", err)
	}
}
