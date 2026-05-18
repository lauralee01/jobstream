package fetcher

import (
	"context"
	"jobstream/internal/domain"
	"time"
)

// Fetcher is an interface that represents a job provider (e.g. LinkedIn, Indeed).
// We use an interface here so that the JobService can fetch jobs from MANY sources
type Fetcher interface {
	// Name returns the name of the source (e.g. "LinkedIn")
	Name() string
	// Fetch returns a list of jobs from the source.
	Fetch(ctx context.Context) ([]domain.Job, error)
}

type MockFetcher struct{}

func (m *MockFetcher) Name() string { return "Mock" }

func (m *MockFetcher) Fetch(ctx context.Context) ([]domain.Job, error) {
	return []domain.Job{{ID: "1", SourceID: "1", Platform: "Mock", Title: "Software Engineer Mock", Company: "Mock Company", Location: "Mock Location", Description: "Mock Description", URL: "Mock URL", Salary: "Mock Salary", PostedAt: time.Now(), CreatedAt: time.Now()}}, nil
}
