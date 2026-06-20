package http

import (
	"jobstream/internal/domain"
	"net/http"
	"strconv"
	"strings"
)

const (
	defaultPage  = 1
	defaultLimit = 20
	maxLimit     = 100
)

var allowedSortColumns = map[string]bool{
	"posted_at":  true,
	"created_at": true,
	"title":      true,
	"company":    true,
}

// parseJobFilter maps HTTP query parameters to a domain JobFilter.
func parseJobFilter(r *http.Request) domain.JobFilter {
	query := r.URL.Query()

	platforms := []string{}
	if raw := query.Get("platforms"); raw != "" {
		platforms = strings.Split(raw, ",")
	}

	page := defaultPage
	if p := query.Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	limit := defaultLimit
	if l := query.Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
			if limit > maxLimit {
				limit = maxLimit
			}
		}
	}

	var isRemote *bool
	if rVal := query.Get("remote"); rVal != "" {
		if parsed, err := strconv.ParseBool(rVal); err == nil {
			isRemote = &parsed
		}
	}

	var minSalary *int
	if sVal := query.Get("min_salary"); sVal != "" {
		if parsed, err := strconv.Atoi(sVal); err == nil {
			minSalary = &parsed
		}
	}

	sortBy := "posted_at"
	sortOrder := "desc"
	if s := query.Get("sort_by"); allowedSortColumns[s] {
		sortBy = s
	}
	if o := strings.ToLower(query.Get("sort_order")); o == "asc" || o == "desc" {
		sortOrder = o
	}

	return domain.JobFilter{
		Keyword:   query.Get("keyword"),
		Location:  query.Get("location"),
		Category:  query.Get("category"),
		MinSalary: minSalary,
		Platforms: platforms,
		IsRemote:  isRemote,
		Page:      page,
		Limit:     limit,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
}
