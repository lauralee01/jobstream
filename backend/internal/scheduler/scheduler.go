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

// Start starts the background ticker.
func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)

	go func() {
		for {
			select {
			case <-ticker.C:
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
					continue
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

			case <-ctx.Done():
				log.Println("Scheduler stopped")
				ticker.Stop()
				return
			}
		}
	}()
}
