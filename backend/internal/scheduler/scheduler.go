package scheduler

import (
	"jobstream/internal/jobs"
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
// Hint: Use a time.Ticker and a 'for' loop with 'select'.
func (s *Scheduler) Start() {
	// TODO: Implement the ticker logic here.
	// ticker := time.NewTicker(s.interval)
	// go func() {
	//    for range ticker.C {
	//        s.jobService.SyncJobs()
	//    }
	// }()
}
