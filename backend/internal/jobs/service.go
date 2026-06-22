package jobs

import (
	"context"
	"fmt"
	"jobstream/internal/category"
	"jobstream/internal/domain"
	"jobstream/internal/fetcher"
	"jobstream/internal/remote"
	"jobstream/internal/salary"
	"log"
	"sort"
	"sync"
	"time"
)

type JobService struct {
	repo     domain.JobRepository
	fetchers []fetcher.Fetcher
}

type ProviderResult struct {
	Provider string `json:"provider"`
	Fetched  int    `json:"fetched"`
	Saved    int    `json:"saved"`
	Error    string `json:"error,omitempty"`
}

type SyncResult struct {
	Fetched       int              `json:"fetched"`
	Saved         int              `json:"saved"`
	FailedSources []string         `json:"failed_sources"`
	Providers     []ProviderResult `json:"providers"`
}

func NewJobService(repo domain.JobRepository, fetchers []fetcher.Fetcher) *JobService {
	return &JobService{
		repo:     repo,
		fetchers: fetchers,
	}
}

func dedupeJobs(jobs []domain.Job) []domain.Job {
	seen := make(map[string]struct{})
	deduped := make([]domain.Job, 0, len(jobs))

	for _, job := range jobs {
		key := job.SourceID + "|" + job.Platform

		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}
		deduped = append(deduped, job)
	}

	return deduped
}

func (s *JobService) SyncJobs(ctx context.Context) (SyncResult, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	result := SyncResult{
		FailedSources: []string{},
		Providers:     []ProviderResult{},
	}

	for _, f := range s.fetchers {
		fetcher := f

		wg.Add(1)
		go func() {
			defer wg.Done()

			providerResult := ProviderResult{
				Provider: fetcher.Name(),
			}

			jobs, err := fetcher.Fetch(ctx)
			if err != nil {
				msg := fmt.Sprintf("%s: fetch failed: %v", fetcher.Name(), err)
				log.Println(msg)

				providerResult.Error = msg

				mu.Lock()
				result.FailedSources = append(result.FailedSources, msg)
				result.Providers = append(result.Providers, providerResult)
				mu.Unlock()

				return
			}

			for i := range jobs {
				jobs[i].Platform = fetcher.Name()
				jobs[i].Category = category.Normalize(jobs[i].Category, jobs[i].Title)
				jobs[i].IsRemote = remote.Detect(jobs[i])
				jobs[i].Active = true
				jobs[i].LastSeenAt = time.Now()
				parsed := salary.Parse(jobs[i].Salary)
				jobs[i].SalaryMin = parsed.Min
				jobs[i].SalaryMax = parsed.Max
			}

			jobs = dedupeJobs(jobs)

			const batchSize = 500
			savedCount := 0

			for start := 0; start < len(jobs); start += batchSize {
				end := start + batchSize
				if end > len(jobs) {
					end = len(jobs)
				}

				batch := jobs[start:end]

				if err := s.repo.Save(ctx, batch); err != nil {
					msg := fmt.Sprintf("%s: save failed: %v", fetcher.Name(), err)
					log.Println(msg)

					providerResult.Fetched = len(jobs)
					providerResult.Saved = savedCount
					providerResult.Error = msg

					mu.Lock()
					result.FailedSources = append(result.FailedSources, msg)
					result.Providers = append(result.Providers, providerResult)
					mu.Unlock()

					return
				}

				savedCount += len(batch)
			}

			providerResult.Fetched = len(jobs)
			providerResult.Saved = savedCount

			mu.Lock()
			result.Fetched += len(jobs)
			result.Saved += savedCount
			result.Providers = append(result.Providers, providerResult)
			mu.Unlock()
		}()
	}

	wg.Wait()

	if result.Saved == 0 && len(result.FailedSources) > 0 {
		return result, fmt.Errorf("all sync attempts failed")
	}

	if err := s.repo.MarkStaleInactive(ctx); err != nil {
		log.Printf("failed to mark stale jobs inactive: %v", err)
	}

	if err := s.repo.DeleteOldInactive(ctx); err != nil {
		log.Printf("failed to delete old inactive jobs: %v", err)
	}

	return result, nil
}

func (s *JobService) GetJobs(ctx context.Context, filter domain.JobFilter) ([]domain.Job, int64, error) {
	if filter.Category != "" {
		filter.Category = category.Normalize(filter.Category, filter.Category)
	}

	return s.repo.FindAll(ctx, filter)
}

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

func (s *JobService) GetPlatforms(ctx context.Context) ([]string, error) {
	return s.repo.GetPlatforms(ctx)
}
