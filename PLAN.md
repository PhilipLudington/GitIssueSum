# gitissuesum — GitHub Issue Summarizer CLI (Go)

## Overview
A Go CLI tool that fetches open issues from a GitHub repo and sends them to Claude for an AI-powered summary.

```
gitissuesum owner/repo
```

## Setup
- [x] Initialize Go module and add cobra dependency
- [x] Create `main.go` entry point

## CLI Layer
- [x] `cmd/root.go` — Cobra root command with arg validation, env var reading
- [x] Flags: `--max-issues` (default 200), `--model` (default "claude-sonnet-4-20250514")

## GitHub Client
- [x] `internal/github/types.go` — Issue, Label, User structs
- [x] `internal/github/client.go` — Paginated fetch via REST API, PR filtering, Link header parsing

## Claude Client
- [x] `internal/claude/types.go` — Request/response structs for Messages API
- [x] `internal/claude/client.go` — HTTP call to Anthropic Messages API

## Orchestration
- [x] `internal/summarize/summarize.go` — Wire fetch → build prompt → call Claude
- [x] Prompt: issue metadata (number, title, truncated body, labels, comments, author, date)
- [x] Output: high-level summary, themes/categories, patterns, top 5 issues

## Verification
- [x] Build succeeds: `go build -o gitissuesum .`
- [ ] Manual test against a public repo
- [x] Error cases: missing arg, invalid format, 404, missing API key
