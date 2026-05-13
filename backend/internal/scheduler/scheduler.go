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
                if err := s.jobService.SyncJobs(ctx); err != nil {
                    log.Printf("Scheduler error: %v\n", err)
                }
            case <-ctx.Done():
                log.Println("Scheduler stopped")
                ticker.Stop()
                return
            }
        }
    }()
}
