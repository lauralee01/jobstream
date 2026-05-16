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
	PostedAt    time.Time `json:"posted_at"`
	CreatedAt   time.Time `json:"created_at"`
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
	Save(ctx context.Context, job *Job) error
	FindAll(ctx context.Context) ([]Job, error)
}
