package jobs

import (
	"context"
	"errors"
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
	var hasErrors bool

	for _, f := range s.fetchers {
		jobs, err := f.Fetch(ctx)
		if err != nil {
			log.Printf("Error fetching jobs from %s: %v", f.Name(), err)
			hasErrors = true
			continue
		}

		for _, job := range jobs {
			job.Platform = f.Name()

			if err := s.repo.Save(ctx, &job); err != nil {
				log.Printf(
					"Error saving job from %s (id=%s): %v",
					f.Name(),
					job.ID,
					err,
				)

				hasErrors = true
				continue
			}
		}
	}

	if hasErrors {
		return errors.New("some operations failed during sync")
	}

	return nil
}

// GetJobs returns all jobs from the database.
func (s *JobService) GetJobs(ctx context.Context, filter domain.JobFilter) ([]domain.Job, int64, error) {
	return s.repo.FindAll(ctx, filter)
}

// GetCategories returns all unique job categories.
func (s *JobService) GetCategories(ctx context.Context) ([]string, error) {
	return s.repo.GetCategories(ctx)
}

// GetPlatforms returns all unique job platforms (sources).
func (s *JobService) GetPlatforms(ctx context.Context) ([]string, error) {
	return s.repo.GetPlatforms(ctx)
}
