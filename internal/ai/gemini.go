package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"google.golang.org/genai"
)

const modelName = "gemini-3.5-flash"

// CommitSuggestion adalah hasil analisis AI terhadap staged Git diff.
type CommitSuggestion struct {
	Type    string `json:"type"`
	Scope   string `json:"scope"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// GenerateCommitSuggestion menganalisis staged Git diff menggunakan Gemini
// dan mengembalikan commit suggestion dalam bentuk data terstruktur.
func GenerateCommitSuggestion(diff string) (*CommitSuggestion, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")

	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	prompt := fmt.Sprintf(`
You are an expert software engineer and Git commit message generator.

Your job is to analyze the ENTIRE staged Git diff and determine the
PRIMARY PURPOSE of the changes.

Do NOT describe each file individually.
Do NOT simply repeat filenames.
Do NOT focus on the largest file.
Understand how all changes work together and identify the central purpose.

CONVENTIONAL COMMIT TYPES:

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
3. Scope may be empty if there is no meaningful scope.
4. Keep the subject concise and specific.
5. Subject MUST be no longer than 72 characters.
6. Use imperative mood.
7. Do not end the subject with a period.
8. Do not mention individual filenames unless the filename itself is the
   main subject of the change.
9. Do not use vague subjects such as:
   - update files
   - modify code
   - changes
   - various fixes
   - update project
10. Do not invent functionality that is not present in the diff.
11. If multiple files support one feature, describe the feature rather
    than listing the files.
12. Prefer a meaningful scope over generic scopes such as:
    "code", "project", or "files".
13. Only include a body when the change is complex enough that the
    subject alone cannot adequately explain the implementation.
14. Keep the body concise.
15. Do not explain your reasoning.

OUTPUT FORMAT:

Return ONLY valid JSON.

The JSON MUST contain exactly these fields:

{
  "type": "feat",
  "scope": "auth",
  "subject": "add passwordless authentication",
  "body": ""
}

IMPORTANT:

- type MUST be one of:
  feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert
- scope must be a short meaningful identifier or an empty string.
- subject must be no longer than 72 characters.
- body must be an empty string for simple changes.
- Do NOT use Markdown.
- Do NOT wrap the JSON in backticks.
- Do NOT add any text before or after the JSON.

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
		return nil, fmt.Errorf("Gemini API request failed: %w", err)
	}

	if response == nil {
		return nil, fmt.Errorf("Gemini returned an empty response")
	}

	text := strings.TrimSpace(response.Text())

	if text == "" {
		return nil, fmt.Errorf("Gemini returned an empty response")
	}

	var suggestion CommitSuggestion

	if err := json.Unmarshal([]byte(text), &suggestion); err != nil {
		return nil, fmt.Errorf("Gemini returned invalid JSON: %w", err)
	}

	if err := validateCommitSuggestion(&suggestion); err != nil {
		return nil, err
	}

	return &suggestion, nil
}

// validateCommitSuggestion memastikan hasil AI mengikuti aturan
// Conventional Commits sebelum digunakan untuk membuat commit.
func validateCommitSuggestion(suggestion *CommitSuggestion) error {
	allowedTypes := map[string]bool{
		"feat":     true,
		"fix":      true,
		"docs":     true,
		"style":    true,
		"refactor": true,
		"perf":     true,
		"test":     true,
		"build":    true,
		"ci":       true,
		"chore":    true,
		"revert":   true,
	}

	suggestion.Type = strings.TrimSpace(suggestion.Type)
	suggestion.Scope = strings.TrimSpace(suggestion.Scope)
	suggestion.Subject = strings.TrimSpace(suggestion.Subject)
	suggestion.Body = strings.TrimSpace(suggestion.Body)

	if !allowedTypes[suggestion.Type] {
		return fmt.Errorf(
			"Gemini returned invalid commit type: %q",
			suggestion.Type,
		)
	}

	if suggestion.Subject == "" {
		return fmt.Errorf("Gemini returned an empty commit subject")
	}

	if len([]rune(suggestion.Subject)) > 72 {
		return fmt.Errorf(
			"Gemini subject exceeds 72 characters (%d characters)",
			len([]rune(suggestion.Subject)),
		)
	}

	// Subject tidak boleh diakhiri titik.
	suggestion.Subject = strings.TrimSuffix(
		suggestion.Subject,
		".",
	)

	return nil
}

// FormatCommitMessage mengubah CommitSuggestion menjadi
// Conventional Commit message yang siap diberikan kepada git commit.
func FormatCommitMessage(suggestion *CommitSuggestion) string {
	var msg strings.Builder

	msg.WriteString(suggestion.Type)

	if suggestion.Scope != "" {
		msg.WriteString("(")
		msg.WriteString(suggestion.Scope)
		msg.WriteString(")")
	}

	msg.WriteString(": ")
	msg.WriteString(suggestion.Subject)

	if suggestion.Body != "" {
		msg.WriteString("\n\n")
		msg.WriteString(suggestion.Body)
	}

	return msg.String()
}
