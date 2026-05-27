package jobs

import (
	"context"
	"jobstream/internal/category"
	"jobstream/internal/domain"
	"jobstream/internal/fetcher"
	"log"
	"regexp"
	"sort"
	"strings"
)

// JobService is a 'Service' in DDD.
// It coordinates the work between the domain models, the fetchers, and the database.
type JobService struct {
	repo     domain.JobRepository
	fetchers []fetcher.Fetcher
}

var salaryDigitRegex = regexp.MustCompile(`\d`)

func hasUsableSalary(salary string) bool {
	value := strings.TrimSpace(strings.ToLower(salary))
	if value == "" {
		return false
	}

	// Keep only salaries with at least one digit. This avoids values like
	// "competitive" while preserving ranges and shorthand such as "95k".
	return salaryDigitRegex.MatchString(value)
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
	for _, f := range s.fetchers {
		jobs, err := f.Fetch(ctx)
		if err != nil {
			log.Printf("Error fetching jobs from %s: %v", f.Name(), err)
			continue
		}

		for _, job := range jobs {
			job.Platform = f.Name()
			job.Category = category.Normalize(job.Category, job.Title)

			if !hasUsableSalary(job.Salary) {
				log.Printf(
					"Skipping job without usable salary from %s (id=%s, title=%q)",
					f.Name(),
					job.ID,
					job.Title,
				)
				continue
			}

			if err := s.repo.Save(ctx, &job); err != nil {
				log.Printf(
					"Error saving job from %s (id=%s): %v",
					f.Name(),
					job.ID,
					err,
				)
				continue
			}
		}
	}

	return nil
}

// GetJobs returns all jobs from the database.
func (s *JobService) GetJobs(ctx context.Context, filter domain.JobFilter) ([]domain.Job, int64, error) {
	return s.repo.FindAll(ctx, filter)
}

// GetCategories returns all unique job categories.
func (s *JobService) GetCategories(ctx context.Context) ([]string, error) {
	categories, err := s.repo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	normalizedSet := map[string]struct{}{}
	for _, cat := range categories {
		normalized := category.Normalize(cat, cat)
		if normalized == "" || normalized == "Other" {
			continue
		}
		normalizedSet[normalized] = struct{}{}
	}

	normalizedCategories := make([]string, 0, len(normalizedSet))
	for cat := range normalizedSet {
		normalizedCategories = append(normalizedCategories, cat)
	}
	sort.Strings(normalizedCategories)

	return normalizedCategories, nil
}

// GetPlatforms returns all unique job platforms (sources).
func (s *JobService) GetPlatforms(ctx context.Context) ([]string, error) {
	return s.repo.GetPlatforms(ctx)
}
