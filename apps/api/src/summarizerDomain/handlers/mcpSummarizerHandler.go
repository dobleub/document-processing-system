package handlers

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
	"google.golang.org/genai"

	"nx-recipes/dps/lambda/config"
	sd_interfaces "nx-recipes/dps/lambda/src/summarizerDomain/interfaces"
	summarizerDomainLib "nx-recipes/dps/lambda/src/summarizerDomain/lib"
)

// @BasePath /summarizer

// MCPHandler godoc
// @Summary MCP Handler for Summarizer
// @Schemes
// @Description Handle MCP requests for summarization
// @Tags Summarizer
// @Accept json
// @Produce json
// @Param data body sd_interfaces.SummarizeInput true "Summarize Input"
// @Success 200 {string} string "Summarize Output"
// @Failure 400 {string} string "Bad Request"
// @Router /summarizer/mcp [post]
func SetUpMCPHandler(ctx context.Context, env *config.Config, logger *zap.Logger) *sd_interfaces.SummarizerHandler {
	// add mcp client to context
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: env.GeminiAPIKey,
	})
	if err != nil {
		logger.Fatal("error creating MCP client", zap.Error(err))
	}

	// implement mcp
	impl := &mcp.Implementation{
		Name:    "gemini-file-summarizer",
		Title:   "Gemini-powered file summarization (local paths + direct upload API)",
		Version: "v1.0.0",
	}
	server := mcp.NewServer(impl, nil)

	// register mcp tools
	summarizerDomainLib.RegisterSummarizeTool(server, client)

	return &sd_interfaces.SummarizerHandler{
		Server: server,
		Client: client,
	}
}
