package cmd

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/mrphil/gitissuesum/internal/summarize"
	"github.com/spf13/cobra"
)

var validRepoName = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

var (
	maxIssues int
	model     string
)

var rootCmd = &cobra.Command{
	Use:   "gitissuesum <owner/repo or GitHub URL>",
	Short: "Summarize open GitHub issues using Claude",
	Long:  "Fetches open issues from a GitHub repository and generates an AI-powered summary using Claude.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		owner, name, err := parseRepo(args[0])
		if err != nil {
			return err
		}

		apiKey := os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			return fmt.Errorf("ANTHROPIC_API_KEY environment variable is required")
		}

		githubToken := os.Getenv("GITHUB_TOKEN")

		return summarize.Run(cmd.Context(), owner, name, apiKey, githubToken, model, maxIssues)
	},
}

func init() {
	rootCmd.Flags().IntVar(&maxIssues, "max-issues", 200, "Maximum number of issues to fetch")
	rootCmd.Flags().StringVar(&model, "model", "claude-sonnet-4-20250514", "Claude model to use")
}

func parseRepo(arg string) (owner, repo string, err error) {
	if strings.Contains(arg, "://") {
		u, parseErr := url.Parse(arg)
		if parseErr != nil {
			return "", "", fmt.Errorf("invalid URL %q: %w", arg, parseErr)
		}
		arg = strings.TrimPrefix(u.Path, "/")
		arg = strings.TrimSuffix(arg, ".git")
	}

	parts := strings.SplitN(arg, "/", 3)
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid repo format %q, expected owner/repo or a GitHub URL", arg)
	}
	if !validRepoName.MatchString(parts[0]) || !validRepoName.MatchString(parts[1]) {
		return "", "", fmt.Errorf("invalid owner or repo name in %q", arg)
	}
	return parts[0], parts[1], nil
}

func Execute() error {
	return rootCmd.Execute()
}
