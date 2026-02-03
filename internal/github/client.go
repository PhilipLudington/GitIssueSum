package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	linkNextRe = regexp.MustCompile(`<([^>]+)>;\s*rel="next"`)
	baseURL    = "https://api.github.com"
	httpClient = &http.Client{Timeout: 30 * time.Second}
)

func FetchIssues(ctx context.Context, owner, repo, token string, maxIssues int) ([]Issue, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues?state=open&per_page=100", baseURL, owner, repo)

	var allIssues []Issue
	for url != "" && len(allIssues) < maxIssues {
		issues, nextURL, err := fetchPage(ctx, url, token)
		if err != nil {
			return nil, err
		}
		for _, issue := range issues {
			if issue.PullRequest != nil {
				continue
			}
			allIssues = append(allIssues, issue)
			if len(allIssues) >= maxIssues {
				break
			}
		}
		url = nextURL
	}

	return allIssues, nil
}

func fetchPage(ctx context.Context, url, token string) ([]Issue, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "gitissuesum")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := doWithRetry(req)
	if err != nil {
		return nil, "", fmt.Errorf("GitHub API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("GitHub API returned status %d for %s", resp.StatusCode, url)
	}

	var issues []Issue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, "", fmt.Errorf("failed to decode GitHub response: %w", err)
	}

	nextURL := parseNextLink(resp.Header.Get("Link"))
	return issues, nextURL, nil
}

func doWithRetry(req *http.Request) (*http.Response, error) {
	backoff := []time.Duration{0, 1 * time.Second, 2 * time.Second}
	var resp *http.Response
	var err error
	for i, wait := range backoff {
		if i > 0 {
			time.Sleep(wait)
		}
		resp, err = httpClient.Do(req)
		if err != nil {
			continue
		}
		switch resp.StatusCode {
		case http.StatusTooManyRequests, http.StatusBadGateway,
			http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			resp.Body.Close()
			continue
		}
		return resp, nil
	}
	return resp, err
}

func parseNextLink(header string) string {
	if header == "" {
		return ""
	}
	for _, part := range strings.Split(header, ",") {
		if matches := linkNextRe.FindStringSubmatch(part); len(matches) == 2 {
			return matches[1]
		}
	}
	return ""
}
