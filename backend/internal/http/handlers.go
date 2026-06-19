package http

import (
	"context"
	"encoding/json"
	"jobstream/internal/cache"
	"jobstream/internal/domain"
	"jobstream/internal/jobs"
	"log"
	"math"
	"net/http"
	"time"
)

type JobHandler struct {
	service       *jobs.JobService
	metadataCache *cache.MetadataCache
}

func NewJobHandler(service *jobs.JobService) *JobHandler {
	return &JobHandler{
		service:       service,
		metadataCache: cache.NewMetadataCache(10 * time.Minute), //Cache categories/platforms for 10 min
	}
}

func (h *JobHandler) SyncJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	result, err := h.service.SyncJobs(r.Context())
	if err != nil {
		log.Printf("SyncJobs failed: %v", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Failed to sync jobs",
			"result":  result,
		})
		return
	}

	h.metadataCache.Invalidate()

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Jobs synced successfully",
		"result":  result,
	})
}

// GetCategories returns a distinct list of non-empty categories.
func (h *JobHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.metadataCache.GetCategories(r.Context(), func(ctx context.Context) ([]string, error) {
		return h.service.GetCategories(ctx)
	})
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
	platforms, err := h.metadataCache.GetPlatforms(r.Context(), func(ctx context.Context) ([]string, error) {
		return h.service.GetPlatforms(ctx)
	})

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
	// Add 15-second timeout to prevent slow queries from blocking indefinitely
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	filter := parseJobFilter(r)

	jobs, total, err := h.service.GetJobs(r.Context(), filter)
	if err != nil {
		// Check if error was due to timeout
		if ctx.Err() == context.DeadlineExceeded {
			http.Error(w, "Request took too long to process. Try refining your filters.", http.StatusGatewayTimeout)
			return
		}
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

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Jobs fetched successfully",
		"result":  response,
	})

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
