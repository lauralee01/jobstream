package scheduler

import (
	"context"
	"jobstream/internal/jobs"
	"log"
	"time"
)

// Scheduler represents a background worker that runs periodically.
type Scheduler struct {
	jobService *jobs.JobService
	interval   time.Duration
}

// NewScheduler creates a new scheduler instance.
func NewScheduler(js *jobs.JobService, interval time.Duration) *Scheduler {
	return &Scheduler{
		jobService: js,
		interval:   interval,
	}
}

// runSync performs a single execution of the job synchronization pipeline.
func (s *Scheduler) runSync(ctx context.Context) {
	log.Println("Running scheduled job sync...")

	result, err := s.jobService.SyncJobs(ctx)
	if err != nil {
		log.Printf(
			"Scheduler sync failed. Fetched=%d Saved=%d FailedSources=%v Error=%v",
			result.Fetched,
			result.Saved,
			result.FailedSources,
			err,
		)
		return
	}

	log.Printf(
		"Scheduler sync completed. Fetched=%d Saved=%d FailedSources=%d",
		result.Fetched,
		result.Saved,
		len(result.FailedSources),
	)

	if len(result.FailedSources) > 0 {
		for _, failure := range result.FailedSources {
			log.Printf("Scheduler source failure: %s", failure)
		}
	}
}

// Start starts the background ticker and runs an immediate sync on startup.
func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)

	go func() {
		// Perform an initial sync immediately when scheduler starts
		s.runSync(ctx)

		for {
			select {
			case <-ticker.C:
				s.runSync(ctx)
			case <-ctx.Done():
				log.Println("Scheduler stopped")
				ticker.Stop()
				return
			}
		}
	}()
}
