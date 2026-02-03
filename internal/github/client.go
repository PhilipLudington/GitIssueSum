package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var linkNextRe = regexp.MustCompile(`<([^>]+)>;\s*rel="next"`)

func FetchIssues(owner, repo, token string, maxIssues int) ([]Issue, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?state=open&per_page=100", owner, repo)

	var allIssues []Issue
	for url != "" && len(allIssues) < maxIssues {
		issues, nextURL, err := fetchPage(url, token)
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

func fetchPage(url, token string) ([]Issue, string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
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
