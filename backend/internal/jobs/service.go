package jobs

import (
	"jobstream/internal/domain"
	"jobstream/internal/fetcher"
)

// JobService is a 'Service' in DDD.
// It coordinates the work between the domain models, the fetchers, and the database.
type JobService struct {
	repo     domain.JobRepository
	fetchers []fetcher.Fetcher
}

// NewJobService creates a new job service.
// This is an example of 'Dependency Injection'. We inject the repository and fetchers!
func NewJobService(repo domain.JobRepository, fetchers []fetcher.Fetcher) *JobService {
	return &JobService{
		repo:     repo,
		fetchers: fetchers,
	}
}

// TODO: Implement a method that runs all fetchers and saves jobs to the database.
// func (s *JobService) SyncJobs() error {
//    // 1. Loop through all fetchers
//    // 2. Fetch jobs from each source
//    // 3. Save each job to the repo
//    return nil
// }
