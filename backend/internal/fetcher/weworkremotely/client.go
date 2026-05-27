package weworkremotely

import (
	"context"
	"fmt"
	"jobstream/internal/category"
	"jobstream/internal/domain"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	arbeitNowURL = "https://weworkremotely.com/remote-jobs"
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
			Timeout: 10 * time.Second,
		},
		baseURL: arbeitNowURL,
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

	// Headers
	req.Header.Set(
		"User-Agent",
		"JobStream-Aggregator/1.0",
	)

	req.Header.Set(
		"Accept",
		"application/json",
	)

	// Execute request
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf(
			"http request failed: %w",
			err,
		)
	}

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("WeWorkRemotely board not found: %s", c.baseURL)
		return nil, fmt.Errorf("WeWorkRemotely board not found: %s", c.baseURL)
	}

	defer resp.Body.Close()

	// Validate response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"unexpected status code: %d",
			resp.StatusCode,
		)
	}

	// Decode JSON

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var jobs []domain.Job

	// Each job listing is inside <li class="feature">
	doc.Find("li.feature").Each(func(i int, s *goquery.Selection) {

		title := strings.TrimSpace(s.Find("span.title").Text())
		company := strings.TrimSpace(s.Find("span.company").Text())
		location := strings.TrimSpace(s.Find("span.region").Text())
		rawCategory := strings.TrimSpace(s.Find("span.company").Parent().Parent().Find("h2").Text())

		// Extract URL
		href, exists := s.Find("a").Attr("href")
		if !exists {
			return
		}

		url := "https://weworkremotely.com" + href

		// PostedAt is not provided → use now
		postedAt := time.Now()

		job := domain.Job{
			Title:    title,
			Company:  company,
			Location: location,
			URL:      url,
			PostedAt: postedAt,
			Category: category.Normalize(rawCategory, title),
		}

		log.Printf("Job: %v", job)

		jobs = append(jobs, job)
	})

	return jobs, nil
}
