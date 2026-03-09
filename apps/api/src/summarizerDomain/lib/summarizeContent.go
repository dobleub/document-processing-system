package lib

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

// Gemini Summarization Logic (shared)
func SummarizeContent(ctx context.Context, client *genai.Client, content string, maxWords int) (string, error) {
	if len(content) > 800_000 { // Gemini has a max input size of around 1 million characters, but we want to leave room for the prompt and response
		content = content[:800_000] + "\n... (content truncated due to size)"
	}

	prompt := fmt.Sprintf(
		"Provide a clear, concise, and accurate summary of the following document. "+
			"Focus on main ideas, key facts, arguments, and conclusions. "+
			"Aim for 200-500 words unless specified otherwise.\n\n%s",
		content,
	)
	if maxWords > 0 {
		prompt += fmt.Sprintf("\nLimit the summary to approximately %d words.", maxWords)
	}

	// List all available models for debugging
	// modelsResp, err := client.Models.List(ctx, nil)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to list models: %w", err)
	// }
	// for _, model := range modelsResp.Items {
	// 	fmt.Printf("Model: %s - %s\n", model.Name, model.Description)
	// }

	resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), nil)
	if err != nil {
		return "", fmt.Errorf("gemini generate failed: %w", err)
	}

	var sb strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		sb.WriteString(part.Text)
	}

	return strings.TrimSpace(sb.String()), nil
}
