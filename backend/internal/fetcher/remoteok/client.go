package remoteok

import (
	"context"
	"encoding/json"
	"fmt"
	"jobstream/internal/domain"
	"net/http"
	"time"
)

const (
	baseURL      = "https://remoteok.com/api"
	platformName = "RemoteOK"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 20 * time.Second},
	}
}

func (c *Client) Name() string { return platformName }

func (c *Client) Fetch(ctx context.Context) ([]domain.Job, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "JobStream-Aggregator/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var apiJobs []RemoteOKJob
	if err := json.NewDecoder(resp.Body).Decode(&apiJobs); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	// First item is usually metadata — skip it
	if len(apiJobs) > 0 && apiJobs[0].ID == "" {
		apiJobs = apiJobs[1:]
	}

	jobs := make([]domain.Job, 0, len(apiJobs))
	for _, j := range apiJobs {
		jobs = append(jobs, j.toDomain())
	}

	return jobs, nil
}