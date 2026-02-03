package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseNextLink_Empty(t *testing.T) {
	if got := parseNextLink(""); got != "" {
		t.Errorf("parseNextLink('') = %q, want ''", got)
	}
}

func TestParseNextLink_SingleNext(t *testing.T) {
	header := `<https://api.github.com/repos/o/r/issues?page=2>; rel="next"`
	want := "https://api.github.com/repos/o/r/issues?page=2"
	if got := parseNextLink(header); got != want {
		t.Errorf("parseNextLink() = %q, want %q", got, want)
	}
}

func TestParseNextLink_MultipleLinks(t *testing.T) {
	header := `<https://api.github.com/repos/o/r/issues?page=2>; rel="next", <https://api.github.com/repos/o/r/issues?page=5>; rel="last"`
	want := "https://api.github.com/repos/o/r/issues?page=2"
	if got := parseNextLink(header); got != want {
		t.Errorf("parseNextLink() = %q, want %q", got, want)
	}
}

func TestParseNextLink_NoNextRel(t *testing.T) {
	header := `<https://api.github.com/repos/o/r/issues?page=1>; rel="prev", <https://api.github.com/repos/o/r/issues?page=5>; rel="last"`
	if got := parseNextLink(header); got != "" {
		t.Errorf("parseNextLink() = %q, want ''", got)
	}
}

func TestFetchIssues_Basic(t *testing.T) {
	issues := []Issue{
		{Number: 1, Title: "bug", User: User{Login: "alice"}},
		{Number: 2, Title: "feature", User: User{Login: "bob"}},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(issues)
	}))
	defer srv.Close()

	old := baseURL
	baseURL = srv.URL
	defer func() { baseURL = old }()

	got, err := FetchIssues("o", "r", "", 100)
	if err != nil {
		t.Fatalf("FetchIssues() error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d issues, want 2", len(got))
	}
	if got[0].Title != "bug" || got[1].Title != "feature" {
		t.Errorf("unexpected issue titles: %v", got)
	}
}

func TestFetchIssues_FiltersPRs(t *testing.T) {
	issues := []Issue{
		{Number: 1, Title: "issue"},
		{Number: 2, Title: "pr", PullRequest: &PullRequest{URL: "https://example.com"}},
		{Number: 3, Title: "another issue"},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(issues)
	}))
	defer srv.Close()

	old := baseURL
	baseURL = srv.URL
	defer func() { baseURL = old }()

	got, err := FetchIssues("o", "r", "", 100)
	if err != nil {
		t.Fatalf("FetchIssues() error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d issues, want 2 (PR should be filtered)", len(got))
	}
	for _, issue := range got {
		if issue.PullRequest != nil {
			t.Errorf("PR not filtered: %v", issue)
		}
	}
}

func TestFetchIssues_Pagination(t *testing.T) {
	page1 := []Issue{{Number: 1, Title: "first"}}
	page2 := []Issue{{Number: 2, Title: "second"}}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") == "2" {
			json.NewEncoder(w).Encode(page2)
			return
		}
		w.Header().Set("Link", fmt.Sprintf(`<%s/repos/o/r/issues?page=2>; rel="next"`, "http://"+r.Host))
		json.NewEncoder(w).Encode(page1)
	}))
	defer srv.Close()

	old := baseURL
	baseURL = srv.URL
	defer func() { baseURL = old }()

	got, err := FetchIssues("o", "r", "", 100)
	if err != nil {
		t.Fatalf("FetchIssues() error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d issues, want 2", len(got))
	}
	if got[0].Title != "first" || got[1].Title != "second" {
		t.Errorf("unexpected order: %v", got)
	}
}

func TestFetchIssues_MaxIssuesLimit(t *testing.T) {
	issues := []Issue{
		{Number: 1, Title: "a"},
		{Number: 2, Title: "b"},
		{Number: 3, Title: "c"},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(issues)
	}))
	defer srv.Close()

	old := baseURL
	baseURL = srv.URL
	defer func() { baseURL = old }()

	got, err := FetchIssues("o", "r", "", 2)
	if err != nil {
		t.Fatalf("FetchIssues() error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d issues, want 2 (maxIssues=2)", len(got))
	}
}

func TestFetchIssues_ErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	old := baseURL
	baseURL = srv.URL
	defer func() { baseURL = old }()

	_, err := FetchIssues("o", "r", "", 100)
	if err == nil {
		t.Fatal("expected error for 404 status")
	}
}

func TestFetchIssues_AuthHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization header = %q, want 'Bearer test-token'", got)
		}
		json.NewEncoder(w).Encode([]Issue{})
	}))
	defer srv.Close()

	old := baseURL
	baseURL = srv.URL
	defer func() { baseURL = old }()

	_, err := FetchIssues("o", "r", "test-token", 100)
	if err != nil {
		t.Fatalf("FetchIssues() error: %v", err)
	}
}
