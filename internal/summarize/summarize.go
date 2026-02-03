package summarize

import (
	"context"
	"fmt"
	"strings"

	"github.com/mrphil/gitissuesum/internal/claude"
	"github.com/mrphil/gitissuesum/internal/github"
)

const maxBodyChars = 500

func Run(ctx context.Context, owner, repo, apiKey, githubToken, model string, maxIssues int) error {
	fmt.Printf("Fetching issues from %s/%s...\n", owner, repo)

	issues, err := github.FetchIssues(ctx, owner, repo, githubToken, maxIssues)
	if err != nil {
		return fmt.Errorf("failed to fetch issues: %w", err)
	}

	if len(issues) == 0 {
		fmt.Println("No open issues found.")
		return nil
	}

	fmt.Printf("Found %d issues. Sending to Claude for analysis...\n", len(issues))

	prompt := buildPrompt(owner, repo, issues)

	response, err := claude.SendMessage(ctx, apiKey, model, prompt)
	if err != nil {
		return fmt.Errorf("failed to get summary from Claude: %w", err)
	}

	fmt.Println()
	fmt.Println(response)
	return nil
}

func buildPrompt(owner, repo string, issues []github.Issue) string {
	var b strings.Builder

	fmt.Fprintf(&b, "You are analyzing open GitHub issues for the repository %s/%s.\n", owner, repo)
	fmt.Fprintf(&b, "There are %d open issues. Here they are:\n\n", len(issues))

	for _, issue := range issues {
		fmt.Fprintf(&b, "--- Issue #%d ---\n", issue.Number)
		fmt.Fprintf(&b, "Title: %s\n", issue.Title)
		fmt.Fprintf(&b, "Author: %s\n", issue.User.Login)
		fmt.Fprintf(&b, "Created: %s\n", issue.CreatedAt.Format("2006-01-02"))
		fmt.Fprintf(&b, "Comments: %d\n", issue.Comments)

		if len(issue.Labels) > 0 {
			labels := make([]string, len(issue.Labels))
			for i, l := range issue.Labels {
				labels[i] = l.Name
			}
			fmt.Fprintf(&b, "Labels: %s\n", strings.Join(labels, ", "))
		}

		body := truncate(issue.Body, maxBodyChars)
		if body != "" {
			fmt.Fprintf(&b, "Body: %s\n", body)
		}

		b.WriteString("\n")
	}

	b.WriteString(`Please provide:
1. A high-level summary of the open issues (2-3 sentences)
2. Main themes/categories you see, with approximate counts
3. Notable patterns (e.g., recurring problems, areas needing attention)
4. The top 5 most important issues and why they stand out

Be concise and actionable.`)

	return b.String()
}

func truncate(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
