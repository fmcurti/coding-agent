# My Gemini Agent

This project is strongly based on the tutorial from [How to Build an Agent](https://ampcode.com/how-to-build-an-agent). However, this implementation uses the Gemini model instead of Anthropic's.

## Overview

This project implements a coding agent powered by the Gemini language model. The agent can perform the following actions:

*   Read files
*   List files in a directory
*   Edit files

The agent takes user input from the console, sends it to Gemini, and then executes tool calls based on Gemini's response. The results of the tool calls are then sent back to Gemini to refine its response. This process continues until the agent has completed the user's request.

## Core Components

*   `main.go`: The entry point of the application. It initializes the Gemini client, sets up the available tools, and starts the agent's main loop.
*   `agent/agent.go`: This file contains the `Agent` struct and its methods. The `Agent` struct holds the Gemini client, a function to get user input, and the available tools. The `Run` method is the main loop of the agent.
*   `tools/tools.go`: This file defines the `ToolDefinition` struct and related helper functions. It includes functions to generate JSON schemas for tools and to convert the tool definitions into a format that Gemini can understand.
*   `tools/editFile.go`: This file defines the `edit_file` tool, which allows the agent to edit files. It takes the file path, the old string to replace, and the new string as input.
*   `tools/listFiles.go`: This file defines the `list_files` tool, which allows the agent to list files in a directory. It takes an optional path as input; if no path is provided, it lists the files in the current directory.
*   `tools/readFile.go`: This file defines the `read_file` tool, which allows the agent to read the contents of a file. It takes the file path as input.

## Getting Started

1.  Set the `GENAI_API_KEY` environment variable to your Gemini API key.
2.  Run `go run main.go` to start the agent.

## Usage

Type your instructions into the console, and the agent will attempt to follow them, using the available tools as needed.
