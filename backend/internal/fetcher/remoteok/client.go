package remoteok

import (
	"context"
	"encoding/json"
	"fmt"
	"jobstream/internal/domain"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	platformName = "RemoteOK"
	baseURL      = "https://remoteok.com/api"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Name() string { return platformName }

func (c *Client) Fetch(ctx context.Context) ([]domain.Job, error) {
	var lastErr error

	for attempt := 0; attempt < 4; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("User-Agent", "JobStream-Aggregator/1.0 (https://github.com/example/jobstream)")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Referer", "https://remoteok.com/")
		req.Header.Set("Cache-Control", "no-cache")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			// Retry on network errors
			if uerr, ok := err.(*url.Error); ok {
				if nerr, ok := uerr.Err.(net.Error); ok && (nerr.Timeout() || nerr.Temporary()) {
					time.Sleep(time.Duration(400+attempt*600) * time.Millisecond)
					continue
				}
			}
			return nil, fmt.Errorf("RemoteOK request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("status code: %d", resp.StatusCode)
			if resp.StatusCode == 429 || resp.StatusCode >= 500 {
				time.Sleep(time.Duration(700+attempt*800) * time.Millisecond)
				continue
			}
			return nil, fmt.Errorf("RemoteOK returned status %d", resp.StatusCode)
		}

		var apiJobs []RemoteOKJob
		if err := json.NewDecoder(resp.Body).Decode(&apiJobs); err != nil {
			return nil, fmt.Errorf("failed to decode JSON: %w", err)
		}

		// Skip metadata object if present (first item with no id)
		if len(apiJobs) > 0 && apiJobs[0].ID == "" {
			apiJobs = apiJobs[1:]
		}

		jobs := make([]domain.Job, 0, len(apiJobs))
		for _, j := range apiJobs {
			jobs = append(jobs, j.toDomain())
		}

		return jobs, nil
	}

	return nil, fmt.Errorf("RemoteOK fetch failed after retries: %w", lastErr)
}