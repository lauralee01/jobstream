package http

import (
	"fmt"
	"jobstream/internal/jobs"
	"net/http"
)

func NewRouter(jobService *jobs.JobService) *http.ServeMux {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	// Job routes
	jobHandler := NewJobHandler(jobService)
	mux.HandleFunc("/jobs/sync", jobHandler.SyncJobs)

	return mux
}
