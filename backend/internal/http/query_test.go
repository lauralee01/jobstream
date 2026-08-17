package http

import (
	"net/http/httptest"
	"testing"
)

func TestParseJobFilter(t *testing.T) {
	t.Run("Default values when query is empty", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/jobs", nil)
		filter := parseJobFilter(req)

		if filter.Page != defaultPage {
			t.Errorf("Page = %d, expected %d", filter.Page, defaultPage)
		}
		if filter.Limit != defaultLimit {
			t.Errorf("Limit = %d, expected %d", filter.Limit, defaultLimit)
		}
		if filter.SortBy != "posted_at" {
			t.Errorf("SortBy = %s, expected 'posted_at'", filter.SortBy)
		}
		if filter.SortOrder != "desc" {
			t.Errorf("SortOrder = %s, expected 'desc'", filter.SortOrder)
		}
		if filter.IsRemote != nil {
			t.Errorf("IsRemote = %v, expected nil", filter.IsRemote)
		}
	})

	t.Run("Parse filters with custom query parameters", func(t *testing.T) {
		req := httptest.NewRequest(
			"GET",
			"/api/v1/jobs?keyword=golang&location=remote&min_salary=120000&remote=true&platforms=remotive,adzuna&page=2&limit=50&sort_by=company&sort_order=asc",
			nil,
		)
		filter := parseJobFilter(req)

		if filter.Keyword != "golang" {
			t.Errorf("Keyword = %s, expected 'golang'", filter.Keyword)
		}
		if filter.Location != "remote" {
			t.Errorf("Location = %s, expected 'remote'", filter.Location)
		}
		if filter.MinSalary == nil || *filter.MinSalary != 120000 {
			t.Errorf("MinSalary = %v, expected 120000", filter.MinSalary)
		}
		if filter.IsRemote == nil || *filter.IsRemote != true {
			t.Errorf("IsRemote = %v, expected true", filter.IsRemote)
		}
		if len(filter.Platforms) != 2 || filter.Platforms[0] != "remotive" || filter.Platforms[1] != "adzuna" {
			t.Errorf("Platforms = %v, expected ['remotive', 'adzuna']", filter.Platforms)
		}
		if filter.Page != 2 {
			t.Errorf("Page = %d, expected 2", filter.Page)
		}
		if filter.Limit != 50 {
			t.Errorf("Limit = %d, expected 50", filter.Limit)
		}
		if filter.SortBy != "company" {
			t.Errorf("SortBy = %s, expected 'company'", filter.SortBy)
		}
		if filter.SortOrder != "asc" {
			t.Errorf("SortOrder = %s, expected 'asc'", filter.SortOrder)
		}
	})

	t.Run("Enforce maximum limit capping", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/jobs?limit=500", nil)
		filter := parseJobFilter(req)

		if filter.Limit != maxLimit {
			t.Errorf("Limit = %d, expected capped %d", filter.Limit, maxLimit)
		}
	})
}
