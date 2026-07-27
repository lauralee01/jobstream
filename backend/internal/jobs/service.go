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
	const (
		saveBatchSize          = 500
		maxConcurrentProviders = 4
	)

	syncStartedAt := time.Now()

	var providerWaitGroup sync.WaitGroup
	var resultMutex sync.Mutex

	providerSemaphore := make(
		chan struct{},
		maxConcurrentProviders,
	)

	syncResult := SyncResult{
		FailedSources: []string{},
		Providers:     []ProviderResult{},
	}

	for _, configuredFetcher := range s.fetchers {
		currentFetcher := configuredFetcher

		providerWaitGroup.Add(1)

		go func() {
			defer providerWaitGroup.Done()

			if err := ctx.Err(); err != nil {
				log.Printf(
					"%s: sync skipped because context ended: %v",
					currentFetcher.Name(),
					err,
				)
				return
			}

			providerSemaphore <- struct{}{}
			defer func() {
				<-providerSemaphore
			}()

			providerSyncResult := ProviderResult{
				Provider: currentFetcher.Name(),
			}

			log.Printf(
				"[job-sync] provider started provider=%s",
				currentFetcher.Name(),
			)

			fetchedJobs, err := currentFetcher.Fetch(ctx)
			if err != nil {
				failureMessage := fmt.Sprintf(
					"%s: fetch failed: %v",
					currentFetcher.Name(),
					err,
				)

				log.Println(failureMessage)

				providerSyncResult.Error = failureMessage

				resultMutex.Lock()
				syncResult.FailedSources = append(
					syncResult.FailedSources,
					failureMessage,
				)
				syncResult.Providers = append(
					syncResult.Providers,
					providerSyncResult,
				)
				resultMutex.Unlock()

				return
			}

			if err := ctx.Err(); err != nil {
				log.Printf(
					"%s: sync cancelled after fetch: %v",
					currentFetcher.Name(),
					err,
				)
				return
			}

			for jobIndex := range fetchedJobs {
				currentJob := &fetchedJobs[jobIndex]

				currentJob.Platform = currentFetcher.Name()
				currentJob.Category = category.Normalize(
					currentJob.Category,
					currentJob.Title,
				)
				currentJob.IsRemote = remote.Detect(*currentJob)
				currentJob.Active = true
				currentJob.LastSeenAt = syncStartedAt

				parsedSalary := salary.Parse(currentJob.Salary)
				currentJob.SalaryMin = parsedSalary.Min
				currentJob.SalaryMax = parsedSalary.Max
			}

			deduplicatedJobs := dedupeJobs(fetchedJobs)
			numberOfSavedJobs := 0

			for batchStartIndex := 0; batchStartIndex < len(deduplicatedJobs); batchStartIndex += saveBatchSize {

				if err := ctx.Err(); err != nil {
					failureMessage := fmt.Sprintf(
						"%s: save cancelled: %v",
						currentFetcher.Name(),
						err,
					)

					log.Println(failureMessage)

					providerSyncResult.Fetched =
						len(deduplicatedJobs)
					providerSyncResult.Saved =
						numberOfSavedJobs
					providerSyncResult.Error =
						failureMessage

					resultMutex.Lock()
					syncResult.FailedSources = append(
						syncResult.FailedSources,
						failureMessage,
					)
					syncResult.Providers = append(
						syncResult.Providers,
						providerSyncResult,
					)
					resultMutex.Unlock()

					return
				}

				batchEndIndex :=
					batchStartIndex + saveBatchSize

				if batchEndIndex > len(deduplicatedJobs) {
					batchEndIndex = len(deduplicatedJobs)
				}

				jobBatch := deduplicatedJobs[batchStartIndex:batchEndIndex]

				if err := s.repo.Save(ctx, jobBatch); err != nil {
					failureMessage := fmt.Sprintf(
						"%s: save failed: %v",
						currentFetcher.Name(),
						err,
					)

					log.Println(failureMessage)

					providerSyncResult.Fetched =
						len(deduplicatedJobs)
					providerSyncResult.Saved =
						numberOfSavedJobs
					providerSyncResult.Error =
						failureMessage

					resultMutex.Lock()
					syncResult.FailedSources = append(
						syncResult.FailedSources,
						failureMessage,
					)
					syncResult.Providers = append(
						syncResult.Providers,
						providerSyncResult,
					)
					resultMutex.Unlock()

					return
				}

				numberOfSavedJobs += len(jobBatch)
			}

			providerSyncResult.Fetched =
				len(deduplicatedJobs)
			providerSyncResult.Saved =
				numberOfSavedJobs

			resultMutex.Lock()
			syncResult.Fetched += len(deduplicatedJobs)
			syncResult.Saved += numberOfSavedJobs
			syncResult.Providers = append(
				syncResult.Providers,
				providerSyncResult,
			)
			resultMutex.Unlock()

			log.Printf(
				"[job-sync] provider completed provider=%s fetched=%d saved=%d",
				currentFetcher.Name(),
				len(deduplicatedJobs),
				numberOfSavedJobs,
			)
		}()
	}

	providerWaitGroup.Wait()

	if syncResult.Saved == 0 &&
		len(syncResult.FailedSources) > 0 {
		return syncResult, fmt.Errorf(
			"all sync attempts failed",
		)
	}

	if err := ctx.Err(); err != nil {
		return syncResult, fmt.Errorf(
			"sync context ended before cleanup: %w",
			err,
		)
	}

	if err := s.repo.MarkStaleInactive(ctx); err != nil {
		log.Printf(
			"failed to mark stale jobs inactive: %v",
			err,
		)
	}

	if err := s.repo.DeleteOldInactive(ctx); err != nil {
		log.Printf(
			"failed to delete old inactive jobs: %v",
			err,
		)
	}

	return syncResult, nil
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
