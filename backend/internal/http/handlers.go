package http

import (
	"encoding/json"
	"jobstream/internal/jobs"
	"net/http"
)

type JobHandler struct {
	service *jobs.JobService
}

func NewJobHandler(service *jobs.JobService) *JobHandler {
	return &JobHandler{
		service: service,
	}
}

func (h *JobHandler) SyncJobs(w http.ResponseWriter, r *http.Request) {
	err := h.service.SyncJobs(r.Context())
	if err != nil {
		http.Error(w, "Failed to sync jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Jobs synced successfully",
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetJobs fetches all jobs and returns them as JSON.
func (h *JobHandler) GetJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.service.GetJobs(r.Context())
	if err != nil {
		http.Error(w, "Failed to get jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(jobs); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

