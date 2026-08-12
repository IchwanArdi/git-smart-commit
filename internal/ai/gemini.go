package ai

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

const modelName = "gemini-3.5-flash"

func GenerateCommitMessage(diff string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")

	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %w", err)
	}

	prompt := fmt.Sprintf(`
You are an expert Git commit message generator.

Analyze the following staged Git diff.

Your task is to determine the PRIMARY purpose of the changes.

Generate a Conventional Commit message.

Rules:
- Use this format: type(scope): subject
- Allowed types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert
- Scope is optional.
- Keep the subject concise.
- Do not invent changes.
- Base your answer ONLY on the provided diff.
- Return ONLY the commit message.
- Do not use Markdown.
- Do not explain your reasoning.

Git diff:

%s
`, diff)

	response, err := client.Models.GenerateContent(
		ctx,
		modelName,
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("Gemini API request failed: %w", err)
	}

	if response == nil || response.Text() == "" {
		return "", fmt.Errorf("Gemini returned an empty response")
	}

	return response.Text(), nil
}
