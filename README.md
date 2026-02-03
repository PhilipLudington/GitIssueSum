# GitIssueSum

A CLI tool that summarizes open GitHub issues using Claude. Built as a live demo of a Claude Code workflow.

- [Slide deck](https://example.com/deck)
- [Demo video](https://example.com/video)

## What This Is

This project was created during a live demonstration of how I use Claude Code for everyday development. The entire workflow — writing code, triaging bugs, adding features, running tests, and committing — was done through Claude Code.

## Usage

```bash
export ANTHROPIC_API_KEY="your-key"
export GITHUB_TOKEN="your-token"   # optional, recommended for rate limits

# Using owner/repo format
./gitissuesum anthropics/claude-code

# Using a GitHub URL
./gitissuesum https://github.com/anthropics/claude-code
```

### Options

```
--max-issues int   Maximum number of issues to fetch (default 200)
--model string     Claude model to use (default "claude-sonnet-4-20250514")
```

## Building

```bash
go build -o gitissuesum .
```

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `ANTHROPIC_API_KEY` | Yes | Your Anthropic API key |
| `GITHUB_TOKEN` | No | GitHub personal access token for higher rate limits |
