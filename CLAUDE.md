# GitIssueSum — Claude Code Instructions

## Running Tests

Always use the AirTower wrapper script to run tests:
```bash
./run-tests.sh
```
Do NOT run `go test ./...` directly — use the wrapper script to preserve AirTower integration and result tracking.

## Building

Always use the AirTower wrapper script to build:
```bash
./run-build.sh
```
Do NOT run `go build` directly — use the wrapper script to preserve AirTower integration and result tracking.

## Project Structure

```
go.mod
main.go
cmd/
  root.go                  # Cobra root command
internal/
  github/
    client.go              # GitHub REST API client (pagination)
    types.go               # Issue, Label, User structs
  claude/
    client.go              # Anthropic Messages API client
    types.go               # Request/response structs
  summarize/
    summarize.go           # Orchestration: fetch → build prompt → call Claude
```

## Environment Variables

- `ANTHROPIC_API_KEY` — Required for Claude API access
- `GITHUB_TOKEN` — Optional, recommended for higher GitHub rate limits
