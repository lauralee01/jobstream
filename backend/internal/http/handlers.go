package http

import (
	"encoding/json"
	"jobstream/internal/domain"
	"jobstream/internal/jobs"
	"math"
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

// GetCategories returns a distinct list of non-empty categories.
func (h *JobHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetCategories(r.Context())
	if err != nil {
		http.Error(w, "Failed to get categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetPlatforms returns a distinct list of non-empty job sources/platforms.
func (h *JobHandler) GetPlatforms(w http.ResponseWriter, r *http.Request) {
	platforms, err := h.service.GetPlatforms(r.Context())
	if err != nil {
		http.Error(w, "Failed to get platforms", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(platforms); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetJobs fetches all jobs and returns them as JSON.
func (h *JobHandler) GetJobs(w http.ResponseWriter, r *http.Request) {
	filter := parseJobFilter(r)

	jobs, total, err := h.service.GetJobs(r.Context(), filter)
	if err != nil {
		http.Error(w, "Failed to get jobs", http.StatusInternalServerError)
		return
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = defaultLimit
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	w.Header().Set("Content-Type", "application/json")

	response := domain.JobsResponse{
		Metadata: domain.Metadata{
			TotalPages:   totalPages,
			TotalResults: total,
			Page:         filter.Page,
			Limit:        limit,
		},
		Data: jobs,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
