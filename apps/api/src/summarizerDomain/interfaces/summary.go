package interfaces

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/genai"
)

// SummarizerHandler is the main handler for summarization-related MCP tools and operations
type SummarizerHandler struct {
	Server *mcp.Server
	Client *genai.Client
}

// MCP Tool: summarize_file (for agents using path)
type SummarizeInput struct {
	FilePath  string `json:"file_path" jsonschema:"Absolute local path to the file"`
	MaxLength int    `json:"max_length,omitempty" jsonschema:"Max summary length in words (optional)"`
}

type SummarizeOutput struct {
	Summary string `json:"summary"`
}
