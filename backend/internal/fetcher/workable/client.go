package workable

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

const platformName = "Workable"

type Client struct {
	httpClient *http.Client
	company    string
	baseURL    string
}

func NewClient(company string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		company: company,
		baseURL: fmt.Sprintf(
			"https://apply.workable.com/api/v1/widget/accounts/%s",
			company,
		),
	}
}

func (c *Client) Name() string {
	return platformName
}

func (c *Client) Fetch(ctx context.Context) ([]domain.Job, error) {
	var lastErr error

	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
		if err != nil {
			return nil, fmt.Errorf("workable: failed to create request: %w", err)
		}

		req.Header.Set("User-Agent", "JobStream/1.0")
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err

			if uerr, ok := err.(*url.Error); ok {
				if nerr, ok := uerr.Err.(net.Error); ok && (nerr.Timeout() || nerr.Temporary()) {
					time.Sleep(time.Duration(500+attempt*1000) * time.Millisecond)
					continue
				}
			}

			if nerr, ok := err.(net.Error); ok && (nerr.Timeout() || nerr.Temporary()) {
				time.Sleep(time.Duration(500+attempt*1000) * time.Millisecond)
				continue
			}

			return nil, fmt.Errorf("workable: request failed for company %s: %w", c.company, err)
		}

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("workable returned status %d for company %s", resp.StatusCode, c.company)
			resp.Body.Close()
			time.Sleep(time.Duration(500+attempt*1000) * time.Millisecond)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("workable: unexpected status %d for company %s", resp.StatusCode, c.company)
		}

		var apiResp APIResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("workable: failed to decode response for company %s: %w", c.company, err)
		}

		resp.Body.Close()

		jobs := make([]domain.Job, 0, len(apiResp.Jobs))
		for _, job := range apiResp.Jobs {
			jobs = append(jobs, job.toDomain(c.company))
		}

		return jobs, nil
	}

	return nil, fmt.Errorf("workable: request failed after retries for company %s: %w", c.company, lastErr)
}
