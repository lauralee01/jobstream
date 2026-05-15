package remotive

import (
	"context"
	"encoding/json"
	"fmt"
	"jobstream/internal/domain"
	"net/http"
	"time"
)

const (
	remotiveURL = "https://remotive.com/api/remote-jobs"
	platformName = "Remotive"
)

// Client is a production-grade fetcher for the Remotive API.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Remotive fetcher client.
// We configure a custom HTTP client with explicit timeouts to prevent hanging requests.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // 10 second timeout for the entire request/response cycle
		},
		baseURL: remotiveURL,
	}
}

// Name implements the fetcher.Fetcher interface.
func (c *Client) Name() string {
	return platformName
}

// Fetch implements the fetcher.Fetcher interface.
// It retrieves jobs from Remotive and maps them to our internal domain models.
func (c *Client) Fetch(ctx context.Context) ([]domain.Job, error) {
	// 1. Create a new HTTP request with the provided context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Good practice: set a proper user agent so APIs don't block us as a generic bot
	req.Header.Set("User-Agent", "JobStream-Aggregator/1.0 (https://github.com/example/jobstream)")
	req.Header.Set("Accept", "application/json")

	// 2. Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	// 3. Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 4. Decode the JSON response into our DTO
	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	// 5. Map DTOs to Domain Entities
	jobs := make([]domain.Job, 0, len(apiResp.Jobs))
	for _, remotiveJob := range apiResp.Jobs {
		jobs = append(jobs, remotiveJob.toDomain())
	}

	fmt.Println(jobs)

	return jobs, nil
}
