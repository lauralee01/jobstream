package linkedin

import (
	"jobstream/internal/domain"
)

type Fetcher struct {}

func NewFetcher() *Fetcher {
	return &Fetcher{}
}

func (f *Fetcher) Name() string {
	return "LinkedIn"
}

func (f *Fetcher) Fetch() ([]domain.Job, error) {
	return []domain.Job{}, nil
}