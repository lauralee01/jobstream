package lever

import (
	"context"
	"encoding/json"
	"fmt"
	"jobstream/internal/domain"
	"net/http"
	"time"
)

const platformName = "Lever"

type Client struct {
	httpClient *http.Client
	company    string
	baseURL    string
}

func NewClient(company string) *Client {

	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		company: company,
		baseURL: fmt.Sprintf(
			"https://api.lever.co/v0/postings/%s?mode=json",
			company,
		),
	}
}

func (c *Client) Name() string {
	return "Lever"
}

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
			"unexpected status code %d for company %s",
			resp.StatusCode,
			c.company,
		)
	}

	var leverJobs []LeverJob

	if err := json.NewDecoder(resp.Body).Decode(&leverJobs); err != nil {

		return nil, fmt.Errorf(
			"failed to decode response: %w",
			err,
		)
	}

	jobs := make([]domain.Job, 0, len(leverJobs))

	for _, job := range leverJobs {

		jobs = append(
			jobs,
			job.toDomain(c.company),
		)
	}

	return jobs, nil
}
