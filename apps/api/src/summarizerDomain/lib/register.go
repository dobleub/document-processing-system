package lib

import (
	"context"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/genai"

	sd_interfaces "nx-recipes/dps/lambda/src/summarizerDomain/interfaces"
)

func RegisterSummarizeTool(server *mcp.Server, client *genai.Client) error {
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "summarize_file",
			Description: "Read a local file from the filesystem and generate an AI-powered summary using Google's Gemini model",
		},
		func(_ context.Context, request *mcp.CallToolRequest, input sd_interfaces.SummarizeInput) (result *mcp.CallToolResult, output sd_interfaces.SummarizeOutput, _ error) {
			// Read the file content
			contentBytes, err := os.ReadFile(input.FilePath)
			if err != nil {
				return nil, sd_interfaces.SummarizeOutput{}, fmt.Errorf("failed to read file: %w", err)
			}
			content := string(contentBytes)

			// Generate the summary using Gemini
			summary, err := SummarizeContent(context.Background(), client, content, 0)
			if err != nil {
				return nil, sd_interfaces.SummarizeOutput{}, fmt.Errorf("failed to generate summary: %w", err)
			}

			output = sd_interfaces.SummarizeOutput{
				Summary: summary,
			}
			return &mcp.CallToolResult{}, output, nil
		},
	)
	return nil
}
