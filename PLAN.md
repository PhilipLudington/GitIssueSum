# gitissuesum — GitHub Issue Summarizer CLI (Go)

## Overview
A Go CLI tool that fetches open issues from a GitHub repo and sends them to Claude for an AI-powered summary.

```
gitissuesum owner/repo
```

## Setup
- [ ] Initialize Go module and add cobra dependency
- [ ] Create `main.go` entry point

## CLI Layer
- [ ] `cmd/root.go` — Cobra root command with arg validation, env var reading
- [ ] Flags: `--max-issues` (default 200), `--model` (default "claude-sonnet-4-20250514")

## GitHub Client
- [ ] `internal/github/types.go` — Issue, Label, User structs
- [ ] `internal/github/client.go` — Paginated fetch via REST API, PR filtering, Link header parsing

## Claude Client
- [ ] `internal/claude/types.go` — Request/response structs for Messages API
- [ ] `internal/claude/client.go` — HTTP call to Anthropic Messages API

## Orchestration
- [ ] `internal/summarize/summarize.go` — Wire fetch → build prompt → call Claude
- [ ] Prompt: issue metadata (number, title, truncated body, labels, comments, author, date)
- [ ] Output: high-level summary, themes/categories, patterns, top 5 issues

## Verification
- [ ] Build succeeds: `go build -o gitissuesum .`
- [ ] Manual test against a public repo
- [ ] Error cases: missing arg, invalid format, 404, missing API key
