package claude

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMessage_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Response{
			Content: []ContentBlock{{Type: "text", Text: "hello world"}},
		})
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	got, err := SendMessage(context.Background(), "key", "model", "prompt")
	if err != nil {
		t.Fatalf("SendMessage() error: %v", err)
	}
	if got != "hello world" {
		t.Errorf("got %q, want 'hello world'", got)
	}
}

func TestSendMessage_MultiBlock(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Response{
			Content: []ContentBlock{
				{Type: "text", Text: "part1"},
				{Type: "text", Text: "part2"},
			},
		})
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	got, err := SendMessage(context.Background(), "key", "model", "prompt")
	if err != nil {
		t.Fatalf("SendMessage() error: %v", err)
	}
	if got != "part1part2" {
		t.Errorf("got %q, want 'part1part2'", got)
	}
}

func TestSendMessage_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Error: &APIError{Type: "invalid_request", Message: "bad prompt"},
		})
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	_, err := SendMessage(context.Background(), "key", "model", "prompt")
	if err == nil {
		t.Fatal("expected error for 400 status")
	}
	if got := err.Error(); got != "Anthropic API returned status 400: invalid_request" {
		t.Errorf("error = %q, want 'Anthropic API returned status 400: invalid_request'", got)
	}
}

func TestSendMessage_Non200NoError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{})
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	_, err := SendMessage(context.Background(), "key", "model", "prompt")
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}

func TestSendMessage_EmptyContent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Response{Content: []ContentBlock{}})
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	_, err := SendMessage(context.Background(), "key", "model", "prompt")
	if err == nil {
		t.Fatal("expected error for empty content")
	}
}

func TestSendMessage_RequestValidation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if got := r.Header.Get("X-API-Key"); got != "test-key" {
			t.Errorf("X-API-Key = %q, want 'test-key'", got)
		}
		if got := r.Header.Get("Anthropic-Version"); got != "2023-06-01" {
			t.Errorf("Anthropic-Version = %q, want '2023-06-01'", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want 'application/json'", got)
		}

		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if req.Model != "test-model" {
			t.Errorf("model = %q, want 'test-model'", req.Model)
		}
		if len(req.Messages) != 1 || req.Messages[0].Content != "test-prompt" {
			t.Errorf("unexpected messages: %v", req.Messages)
		}

		json.NewEncoder(w).Encode(Response{
			Content: []ContentBlock{{Type: "text", Text: "ok"}},
		})
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	_, err := SendMessage(context.Background(), "test-key", "test-model", "test-prompt")
	if err != nil {
		t.Fatalf("SendMessage() error: %v", err)
	}
}
