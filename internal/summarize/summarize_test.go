package summarize

import (
	"strings"
	"testing"
	"time"

	"github.com/mrphil/gitissuesum/internal/github"
)

func TestTruncate_Short(t *testing.T) {
	if got := truncate("hello", 10); got != "hello" {
		t.Errorf("truncate('hello', 10) = %q, want 'hello'", got)
	}
}

func TestTruncate_ExactLength(t *testing.T) {
	if got := truncate("hello", 5); got != "hello" {
		t.Errorf("truncate('hello', 5) = %q, want 'hello'", got)
	}
}

func TestTruncate_OverLength(t *testing.T) {
	got := truncate("hello world", 5)
	if got != "hello..." {
		t.Errorf("truncate('hello world', 5) = %q, want 'hello...'", got)
	}
}

func TestTruncate_WhitespaceTrimming(t *testing.T) {
	if got := truncate("  hi  ", 10); got != "hi" {
		t.Errorf("truncate('  hi  ', 10) = %q, want 'hi'", got)
	}
}

func TestBuildPrompt_SingleIssue(t *testing.T) {
	issues := []github.Issue{
		{
			Number:    42,
			Title:     "Test issue",
			Body:      "Some body text",
			User:      github.User{Login: "alice"},
			CreatedAt: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			Comments:  3,
		},
	}

	prompt := buildPrompt("owner", "repo", issues)

	checks := []string{
		"owner/repo",
		"1 open issues",
		"Issue #42",
		"Title: Test issue",
		"Author: alice",
		"Created: 2025-01-15",
		"Comments: 3",
		"Body: Some body text",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing %q", want)
		}
	}
}

func TestBuildPrompt_MultipleIssues(t *testing.T) {
	issues := []github.Issue{
		{Number: 1, Title: "First", User: github.User{Login: "a"}, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Number: 2, Title: "Second", User: github.User{Login: "b"}, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)},
	}

	prompt := buildPrompt("o", "r", issues)

	if !strings.Contains(prompt, "2 open issues") {
		t.Error("prompt should mention 2 open issues")
	}
	if !strings.Contains(prompt, "Issue #1") || !strings.Contains(prompt, "Issue #2") {
		t.Error("prompt should contain both issues")
	}
}

func TestBuildPrompt_Labels(t *testing.T) {
	issues := []github.Issue{
		{
			Number:    1,
			Title:     "Labeled",
			User:      github.User{Login: "a"},
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Labels:    []github.Label{{Name: "bug"}, {Name: "urgent"}},
		},
	}

	prompt := buildPrompt("o", "r", issues)

	if !strings.Contains(prompt, "Labels: bug, urgent") {
		t.Errorf("prompt missing labels, got:\n%s", prompt)
	}
}

func TestBuildPrompt_EmptyBody(t *testing.T) {
	issues := []github.Issue{
		{
			Number:    1,
			Title:     "No body",
			User:      github.User{Login: "a"},
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Body:      "",
		},
	}

	prompt := buildPrompt("o", "r", issues)

	if strings.Contains(prompt, "Body:") {
		t.Error("prompt should not contain Body: line for empty body")
	}
}

func TestBuildPrompt_BodyTruncation(t *testing.T) {
	longBody := strings.Repeat("x", 600)
	issues := []github.Issue{
		{
			Number:    1,
			Title:     "Long body",
			User:      github.User{Login: "a"},
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Body:      longBody,
		},
	}

	prompt := buildPrompt("o", "r", issues)

	if !strings.Contains(prompt, "...") {
		t.Error("long body should be truncated with ...")
	}
	// Body in prompt should be maxBodyChars (500) + "..."
	if strings.Contains(prompt, longBody) {
		t.Error("full 600-char body should not appear in prompt")
	}
}
