package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mrphil/gitissuesum/internal/summarize"
	"github.com/spf13/cobra"
)

var (
	maxIssues int
	model     string
)

var rootCmd = &cobra.Command{
	Use:   "gitissuesum owner/repo",
	Short: "Summarize open GitHub issues using Claude",
	Long:  "Fetches open issues from a GitHub repository and generates an AI-powered summary using Claude.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := args[0]
		parts := strings.SplitN(repo, "/", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid repo format %q, expected owner/repo", repo)
		}

		apiKey := os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			return fmt.Errorf("ANTHROPIC_API_KEY environment variable is required")
		}

		githubToken := os.Getenv("GITHUB_TOKEN")

		return summarize.Run(parts[0], parts[1], apiKey, githubToken, model, maxIssues)
	},
}

func init() {
	rootCmd.Flags().IntVar(&maxIssues, "max-issues", 200, "Maximum number of issues to fetch")
	rootCmd.Flags().StringVar(&model, "model", "claude-sonnet-4-20250514", "Claude model to use")
}

func Execute() error {
	return rootCmd.Execute()
}
