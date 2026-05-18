package http

import (
	"encoding/json"
	"jobstream/internal/domain"
	"jobstream/internal/jobs"
	"log"
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

	// =========================
	// Parse platforms
	// =========================

	platforms := []string{}	

	if query.Get("platforms") != "" {
		platforms = strings.Split(
			query.Get("platforms"),
			",",
		)
	}

	// =========================
	// Parse page
	// =========================

	page := 1

	if p := query.Get("page"); p != "" {
		parsedPage, err := strconv.Atoi(p)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	// =========================
	// Parse limit
	// =========================

	limit := 20

	if l := query.Get("limit"); l != "" {
		parsedLimit, err := strconv.Atoi(l)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// =========================
	// Build filter
	// =========================

	filter := domain.JobFilter{
		Keyword:   query.Get("keyword"),
		Location:  query.Get("location"),
		Category:  query.Get("category"),
		Platforms: platforms,
		Page:      page,
		Limit:     limit,
		SortBy:    "created_at",
		SortOrder: "desc",
	}

	// =========================
	// Fetch jobs
	// =========================

	jobs, total, err := h.service.GetJobs(r.Context(), filter)

	log.Println("Found", total, "jobs")
	log.Println("jobs", jobs)

	if err != nil {
		http.Error(w, "Failed to get jobs", http.StatusInternalServerError)
		return
	}

	// =========================
	// Calculate pagination
	// =========================

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// =========================
	// Response
	// =========================

	w.Header().Set("Content-Type", "application/json")

	response := domain.JobsResponse{
		Metadata: domain.Metadata{
			TotalPages:   totalPages,
			TotalResults: total,
			Page:         page,
			Limit:        limit,
		},
		Data: jobs,
	}

	log.Println("response", response)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
