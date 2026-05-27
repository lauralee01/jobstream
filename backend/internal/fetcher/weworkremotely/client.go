package weworkremotely

import (
	"context"
	"encoding/xml"
	"fmt"
	"jobstream/internal/domain"
	"net/http"
	"time"
)

const (
	feedURL      = "https://weworkremotely.com/remote-jobs.rss"
	platformName = "WeWorkRemotely"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new WeWorkRemotely fetcher client.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
		baseURL: feedURL,
	}
}

func (c *Client) Name() string {
	return platformName
}

// Fetch jobs from WeWorkRemotely
func (c *Client) Fetch(ctx context.Context) ([]domain.Job, error) {

	// Create request with context
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to create request: %w",
			err,
		)
	}

	req.Header.Set("User-Agent", "JobStream-Aggregator/1.0")
	req.Header.Set("Accept", "application/rss+xml, application/xml;q=0.9, */*;q=0.8")

	// Execute request
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf(
			"http request failed: %w",
			err,
		)
	}

	defer resp.Body.Close()

	// Validate response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"unexpected status code: %d",
			resp.StatusCode,
		)
	}

	var rss RSS
	if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		return nil, fmt.Errorf("failed to decode rss xml: %w", err)
	}

	jobs := make([]domain.Job, 0, len(rss.Channel.Items))
	for _, item := range rss.Channel.Items {
		jobs = append(jobs, item.toDomain())
	}

	return jobs, nil
}
