package http

import (
	"encoding/json"
	"jobstream/internal/jobs"
	"net/http"
)

func NewRouter(jobService *jobs.JobService) http.Handler {
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
	// Make sure to handle preflight OPTIONS requests for CORS
	mux.HandleFunc("/api/v1/jobs/sync", jobHandler.SyncJobs)
	mux.HandleFunc("GET /api/v1/jobs", jobHandler.GetJobs)

	return corsMiddleware(mux)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Nuxt frontend
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
