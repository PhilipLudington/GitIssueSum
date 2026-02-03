package claude

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	apiURL     = "https://api.anthropic.com/v1/messages"
	httpClient = &http.Client{Timeout: 120 * time.Second}
)

func SendMessage(ctx context.Context, apiKey, model, prompt string) (string, error) {
	reqBody := Request{
		Model:     model,
		MaxTokens: 4096,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Anthropic-Version", "2023-06-01")

	resp, err := doWithRetry(req, body)
	if err != nil {
		return "", fmt.Errorf("Anthropic API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp Response
		msg := fmt.Sprintf("Anthropic API returned status %d", resp.StatusCode)
		if json.NewDecoder(resp.Body).Decode(&errResp) == nil && errResp.Error != nil {
			msg += ": " + errResp.Error.Type
		}
		return "", fmt.Errorf("%s", msg)
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode Anthropic response: %w", err)
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("empty response from Claude")
	}

	var text string
	for _, block := range result.Content {
		if block.Type == "text" {
			text += block.Text
		}
	}
	return text, nil
}

func doWithRetry(req *http.Request, body []byte) (*http.Response, error) {
	backoff := []time.Duration{0, 1 * time.Second, 2 * time.Second}
	var resp *http.Response
	var err error
	for i, wait := range backoff {
		if i > 0 {
			time.Sleep(wait)
			req.Body = io.NopCloser(bytes.NewReader(body))
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
