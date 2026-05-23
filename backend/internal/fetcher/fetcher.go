package fetcher

import (
	"context"
	"jobstream/internal/domain"
)

// Fetcher is an interface that represents a job provider (e.g. LinkedIn, Indeed).
// We use an interface here so that the JobService can fetch jobs from MANY sources
type Fetcher interface {
	// Name returns the name of the source (e.g. "LinkedIn")
	Name() string
	// Fetch returns a list of jobs from the source.
	Fetch(ctx context.Context) ([]domain.Job, error)
}
