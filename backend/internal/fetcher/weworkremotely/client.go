package weworkremotely

import (
	"context"
	"encoding/xml"
	"fmt"
	"jobstream/internal/domain"
	"log"
	"net/http"
	"time"
)

const (
	weworkremotelyURL = "https://weworkremotely.com/remote-jobs.rss"
	platformName      = "WeWorkRemotely"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new WeWorkRemotely fetcher client.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 15 * time.Second, // WWR RSS feed can be large, use a slightly longer timeout
		},
		baseURL: weworkremotelyURL,
	}
}

func (c *Client) Name() string {
	return platformName
}

// Fetch retrieves jobs from WeWorkRemotely's public RSS feed and normalizes them.
func (c *Client) Fetch(ctx context.Context) ([]domain.Job, error) {
	// Create request with context
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Headers to represent a healthy user-agent
	req.Header.Set("User-Agent", "JobStream-Aggregator/1.0 (https://github.com/example/jobstream)")
	req.Header.Set("Accept", "application/xml, text/xml, */*")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	// Validate response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode XML RSS feed
	var feed RSS
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("failed to decode XML RSS response: %w", err)
	}

	// Map RSS items to domain job entities
	items := feed.Channel.Items
	jobs := make([]domain.Job, 0, len(items))

	for _, item := range items {
		job := item.toDomain()
		jobs = append(jobs, job)
	}

	log.Printf("Successfully fetched and normalized %d jobs from %s", len(jobs), c.Name())
	return jobs, nil
}
