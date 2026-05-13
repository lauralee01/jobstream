package http

import (
	"encoding/json"
	"jobstream/internal/jobs"
	"net/http"
)

func NewRouter(jobService *jobs.JobService) *http.ServeMux {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	// Job routes
	jobHandler := NewJobHandler(jobService)
	mux.HandleFunc("/api/v1/jobs/sync", jobHandler.SyncJobs)
	mux.HandleFunc("GET /api/v1/jobs", jobHandler.GetJobs)
	return mux
}
