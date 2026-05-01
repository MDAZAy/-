package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type EmbedClient struct {
	endpoint   string
	httpClient *http.Client
}

type EmbedResponse struct {
	Vector []float64 `json:"vector"`
}

func NewEmbedClient(endpoint string, timeout time.Duration) *EmbedClient {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	return &EmbedClient{
		endpoint:   strings.TrimRight(endpoint, "/"),
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (c *EmbedClient) Embed(ctx context.Context, text string) (*EmbedResponse, error) {
	if c == nil || c.endpoint == "" {
		return nil, fmt.Errorf("embed endpoint is not configured")
	}

	payload, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return nil, fmt.Errorf("marshal embed request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+"/embed", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build embed request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call embed endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("embed endpoint returned status %s", resp.Status)
	}

	var decoded struct {
		Vector []float64 `json:"vector"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("decode embed response: %w", err)
	}

	return &EmbedResponse{Vector: decoded.Vector}, nil
}
