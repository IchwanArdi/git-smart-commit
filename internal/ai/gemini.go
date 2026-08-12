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
	You are an expert software engineer and Git commit message generator.

	Your job is to analyze the ENTIRE staged Git diff and determine the
	PRIMARY PURPOSE of the changes.

	Do NOT describe each file individually.
	Do NOT simply repeat filenames.
	Do NOT focus on the largest file.
	Understand how all changes work together and summarize the main purpose.

	CONVENTIONAL COMMIT FORMAT:

	type(scope): subject

	Allowed types:
	- feat     = new functionality
	- fix      = bug fix
	- docs     = documentation changes
	- style    = formatting/style changes without logic changes
	- refactor = code restructuring without changing behavior
	- perf     = performance improvement
	- test     = adding or modifying tests
	- build    = build system or dependency changes
	- ci       = CI/CD changes
	- chore    = maintenance changes
	- revert   = reverting a previous change

	RULES:

	1. Choose the commit type based on the PRIMARY purpose of the entire change.
	2. Scope should describe the main affected module or feature.
	3. Keep the subject concise and specific.
	4. Subject MUST be no longer than 72 characters.
	5. Use imperative mood.
	6. Do not end the subject with a period.
	7. Do not mention individual filenames unless the filename itself is the
	main subject of the change.
	8. Do not use vague subjects such as:
	- update files
	- modify code
	- changes
	- various fixes
	- update project
	9. Do not invent functionality that is not present in the diff.
	10. If multiple files support one feature, describe the feature rather than listing the files.
	11. Prefer the most meaningful scope rather than a generic scope such as "code", "project", or "files".
	12. Return ONLY the commit message.
	13. Do not use Markdown.
	14. Do not wrap the result in backticks.
	15. Do not explain your reasoning.
	16. Only add a commit body when the change is complex enough that the subject alone cannot adequately explain the implementation.
	17. For simple changes, return only the subject.
	18. If a body is included, keep it concise and explain the main implementation rather than listing changed files.

	IMPORTANT:
	The staged diff may contain many files and hundreds of lines.
	Analyze the changes as ONE logical change and identify the central purpose.

	STAGED GIT DIFF:

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
