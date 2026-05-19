package adzuna

import (
	"context"
	"encoding/json"
	"fmt"
	"jobstream/internal/domain"
	"log"
	"net/http"
	"os"
	"time"
)


type Client struct {
	httpClient *http.Client
	baseURL    string
	appID    string
	appKey   string
}

const (
	platformName = "Adzuna"
)


func NewAPIClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL:    "https://api.adzuna.com/v1/api/jobs/us/search/1",
		appID:    os.Getenv("ADZUNA_APP_ID"),
		appKey:   os.Getenv("ADZUNA_APP_KEY"),
	}
}

// Name implements the fetcher.Fetcher interface.
func (c *Client) Name() string {
	return platformName
}

func (c *Client) Fetch(ctx context.Context) ([]domain.Job, error) {
	url := fmt.Sprintf(
		"%s?app_id=%s&app_key=%s&results_per_page=50",
		c.baseURL,
		c.appID,
		c.appKey,
	)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "JobStream-Aggregator/1.0 (https://github.com/example/jobstream)")
	req.Header.Set("Accept", "application/json")
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code received: %d", resp.StatusCode)
	}

	var data APIResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}
	jobs := make([]domain.Job, len(data.Results))
	for i, r := range data.Results {
		jobs[i] = r.ToJob()
	}
	log.Printf("Fetched %d jobs from Adzuna", len(jobs), jobs)
	return jobs, nil
}