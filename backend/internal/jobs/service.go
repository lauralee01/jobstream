package jobs

import (
	"context"
	"jobstream/internal/category"
	"jobstream/internal/domain"
	"jobstream/internal/fetcher"
	"log"
	"sort"
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

func (s *JobService) SyncJobs(ctx context.Context) error {
    var wg sync.WaitGroup
    errCh := make(chan error, len(s.fetchers))

    for _, f := range s.fetchers {
        fetcher := f

        wg.Add(1)
        go func() {
            defer wg.Done()

            jobs, err := fetcher.Fetch(ctx)
            if err != nil {
                log.Printf("Error fetching jobs from %s: %v", fetcher.Name(), err)
                errCh <- err
                return
            }

            // Preprocess jobs: platform, category, salary
            for i := range jobs {
                jobs[i].Platform = fetcher.Name()
                jobs[i].Category = category.Normalize(jobs[i].Category, jobs[i].Title)

                parsed := salary.Parse(jobs[i].Salary)
                jobs[i].SalaryMin = parsed.Min
                jobs[i].SalaryMax = parsed.Max
            }

            // Batch size (tune as needed)
            const batchSize = 500

            for start := 0; start < len(jobs); start += batchSize {
                end := start + batchSize
                if end > len(jobs) {
                    end = len(jobs)
                }

                batch := jobs[start:end]

                if err := s.repo.SaveBatch(ctx, batch); err != nil {
                    log.Printf("Error saving batch from %s: %v", fetcher.Name(), err)
                    errCh <- err
                    return
                }
            }
        }()
    }

    wg.Wait()
    close(errCh)
	
    for err := range errCh {
        if err != nil {
            return err
        }
    }

    return nil
}


// GetJobs returns all jobs from the database.
func (s *JobService) GetJobs(ctx context.Context, filter domain.JobFilter) ([]domain.Job, int64, error) {
	if filter.Category != "" {
		// Normalize incoming category filter so both legacy (e.g. "Accounting")
		// and canonical (e.g. "Finance") values work.
		filter.Category = category.Normalize(filter.Category, filter.Category)
	}
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
