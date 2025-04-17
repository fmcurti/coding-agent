package agent

import (
	"coding-agent/tools"
	"context"
	"fmt"

	"google.golang.org/genai"
)

type Agent struct {
	client         *genai.Client
	getUserMessage func() (string, bool)
	tools          tools.ToolDefinitions
}

func NewAgent(client *genai.Client, getUserMessage func() (string, bool), availableTools tools.ToolDefinitions) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
		tools:          availableTools,
	}
}

func (a *Agent) Run(ctx context.Context) error {
	conversation := []*genai.Content{}

	readUserInput := true
	for {
		if readUserInput {
			fmt.Print("\u001b[94mYou\u001b[0m: ")
			userInput, ok := a.getUserMessage()
			if !ok {
				break
			}
			userMessage := genai.Text(userInput)
			conversation = append(conversation, userMessage...)
		}
		message, err := a.runInference(ctx, conversation)
		if err != nil {
			return err
		}

		conversation = append(conversation, message.Candidates[0].Content)
		toolResults := []genai.FunctionResponse{}
		for _, part := range message.Candidates[0].Content.Parts {
			if part.Text != "" {
				if part.Thought {
					continue
				}
				fmt.Printf("\u001b[92mGemini\u001b[0m: %s\n", part.Text)
			}
			if part.FunctionCall != nil {
				call := part.FunctionCall
				result := a.executeTool(call.ID, call.Name, call.Args)
				toolResults = append(toolResults, result)
			}
		}

		if len(toolResults) == 0 {
			readUserInput = true
			continue
		}

		readUserInput = false
		parts := []*genai.Part{}
		for _, toolResult := range toolResults {
			part := genai.Part{
				FunctionResponse: &toolResult,
			}
			parts = append(parts, &part)
		}
		conversation = append(conversation, &genai.Content{
			Parts: parts,
		})

	}
	return nil
}

func (a *Agent) executeTool(id string, name string, args map[string]any) genai.FunctionResponse {
	var toolDef tools.ToolDefinition
	var found bool
	for _, tool := range a.tools {
		if tool.Name == name {
			toolDef = tool
			found = true
			break
		}
	}
	if !found {
		return genai.FunctionResponse{
			ID:   id,
			Name: name,
			Response: map[string]any{
				"error": "Tool not found",
			},
		}
	}

	fmt.Printf("\u001b[92mtool\u001b[0m: %s(%s)\n", name, args)
	response, err := toolDef.Function(args)
	if err != nil {
		return genai.FunctionResponse{
			ID:   id,
			Name: name,
			Response: map[string]any{
				"error": err.Error(),
			},
		}
	}
	return genai.FunctionResponse{
		ID:   id,
		Name: name,
		Response: map[string]any{
			"output": response,
		},
	}
}

func (a *Agent) runInference(ctx context.Context, conversation []*genai.Content) (*genai.GenerateContentResponse, error) {
	functions := a.tools.ToGeminiFunction()
	tools := make([]*genai.Tool, 1)
	tools[0] = &genai.Tool{
		FunctionDeclarations: functions,
	}
	message, err := a.client.Models.GenerateContent(ctx, "gemini-2.0-flash", conversation, &genai.GenerateContentConfig{
		Tools: tools,
	})
	return message, err
}
