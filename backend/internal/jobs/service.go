package jobs

import (
	"context"
	"jobstream/internal/domain"
	"jobstream/internal/fetcher"
	"log"
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
func (s *JobService) SyncJobs(ctx context.Context) error {
   // 1. Loop through all fetchers
   for _, f := range s.fetchers {
		// 2. Fetch jobs from each source
		jobs, err := f.Fetch()
		if err != nil {
			log.Printf("Error fetching jobs from %s: %v", f.Name(), err)
			continue
		}
		// 3. Save each job to the repo
		for _, job := range jobs {
			job.Platform = f.Name()
            if err := s.repo.Save(ctx, &job); err != nil {
                log.Printf("Error saving job from %s (id=%s): %v", f.Name(), job.ID, err)
                continue
            }
		}
   }
  
   return nil
}
