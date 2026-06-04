package domain

import (
	"context"
	"time"
)

// Job represents a normalized job posting in our system.
// In Domain-Driven Design (DDD), this is an Entity.
type Job struct {
	ID          string    `json:"id"`
	SourceID    string    `json:"source_id"` // ID from the external source (e.g. LinkedIn ID)
	Platform    string    `json:"platform"`  // e.g. "linkedin", "indeed"
	Title       string    `json:"title"`
	Company     string    `json:"company"`
	Location    string    `json:"location"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Salary      string    `json:"salary"`
	SalaryMin   *int64 `json:"salary_min,omitempty"`
	SalaryMax *int64 `json:"salary_max,omitempty"`
	PostedAt    time.Time `json:"posted_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// JobFilter defines search and pagination parameters.
type JobFilter struct {
	Keyword   string
	Location  string
	Category  string
	MinSalary *int
	Platforms []string

	IsRemote *bool

	Page  int
	Limit int

	SortBy    string
	SortOrder string
}

// Source represents an external job provider.
type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"` // "api", "rss", "scraper"
}

// JobRepository is an interface that defines how we interact with the database.
// This allows us to swap PostgreSQL for something else later without changing business logic!
type JobRepository interface {
	Save(ctx context.Context, jobs []Job) error
	FindAll(ctx context.Context, filter JobFilter) ([]Job, int64, error)
	GetCategories(ctx context.Context) ([]string, error)
	GetPlatforms(ctx context.Context) ([]string, error)
}

type JobsResponse struct {
	Metadata Metadata `json:"metadata"`
	Data     []Job    `json:"data"`
}

type Metadata struct {
	TotalPages   int   `json:"total_pages"`
	TotalResults int64 `json:"total_results"`
	Page         int   `json:"page"`
	Limit        int   `json:"limit"`
}
