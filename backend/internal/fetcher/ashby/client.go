package ashby

import (
	"context"
	"encoding/json"
	"fmt"
	"jobstream/internal/domain"
	"net/http"
	"time"
)

const platformName = "Ashby"

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
			"https://api.ashbyhq.com/posting-api/job-board/%s?includeCompensation=true",
			company,
		),
	}
}

func (c *Client) Name() string {
	return platformName
}

func (c *Client) Fetch(ctx context.Context) ([]domain.Job, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ashby: failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "JobStream/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ashby: request failed for company %s: %w", c.company, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ashby: unexpected status %d for company %s", resp.StatusCode, c.company)
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("ashby: failed to decode response: %w", err)
	}

	jobs := make([]domain.Job, 0, len(apiResp.Jobs))
	for _, job := range apiResp.Jobs {
		if !job.IsListed {
			continue
		}
		jobs = append(jobs, job.toDomain(c.company))
	}

	return jobs, nil
}
