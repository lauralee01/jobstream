package http

import (
	"fmt"
	"jobstream/internal/jobs"
	"net/http"
)

func NewRouter(jobService *jobs.JobService) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	return mux
}
