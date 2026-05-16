package http

import (
	"encoding/json"
	"jobstream/internal/domain"
	"jobstream/internal/jobs"
	"math"
	"net/http"
	"strconv"
	"strings"
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
	query := r.URL.Query()

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10
	}

	var platforms []string
	if p := query.Get("platform"); p != "" {
		platforms = strings.Split(p, ",")
	}

	var isRemote *bool
	if remote := query.Get("is_remote"); remote != "" {
		val := remote == "true"
		isRemote = &val
	}

	sortBy := query.Get("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}

	sortOrder := query.Get("sort_order")
	if sortOrder == "" {
		sortOrder = "desc"
	}

	filter := domain.JobFilter{
		Keyword:   query.Get("keyword"),
		Location:  query.Get("location"),
		Category:  query.Get("category"),
		Platforms: platforms,
		IsRemote:  isRemote,
		Page:      page,
		Limit:     limit,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	jobs, total, err := h.service.GetJobs(r.Context(), filter)
	if err != nil {
		http.Error(w, "Failed to get jobs", http.StatusInternalServerError)
		return
	}

	totalPages := 1
	if limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	response := domain.JobsResponse{
		Metadata: domain.Metadata{
			TotalPages:   totalPages,
			TotalResults: total,
			Page:         page,
			Limit:        limit,
		},
		Data: jobs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}


