package greenhouse

import (
	"context"
	"encoding/json"
	"fmt"
	"jobstream/internal/domain"
	"log"
	"net/http"
	"time"
)

const baseURL = "https://boards-api.greenhouse.io/v1/boards"

// Client represents a reusable Greenhouse fetcher
type Client struct {
	httpClient *http.Client
	company    string
	baseURL    string
}

// NewClient creates a new Greenhouse client
func NewClient(company string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		company: company,
		baseURL: fmt.Sprintf(
			"%s/%s/jobs",
			baseURL,
			company,
		),
	}
}

// Name implements Fetcher interface
func (c *Client) Name() string {
	return "Greenhouse"
}

// Fetch retrieves jobs
func (c *Client) Fetch(
	ctx context.Context,
) ([]domain.Job, error) {

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

	req.Header.Set(
		"User-Agent",
		"JobStream/1.0",
	)

	req.Header.Set(
		"Accept",
		"application/json",
	)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf(
			"http request failed: %w",
			err,
		)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"unexpected status code: %d",
			resp.StatusCode,
		)
	}

	var apiResp APIResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf(
			"failed to decode response: %w",
			err,
		)
	}

	jobs := make([]domain.Job, 0, len(apiResp.Jobs))

	for _, job := range apiResp.Jobs {
		jobs = append(
			jobs,
			job.toDomain(c.company),
		)
	}

	log.Println("Fetched %d jobs from %s", len(jobs), jobs)

	return jobs, nil
}
