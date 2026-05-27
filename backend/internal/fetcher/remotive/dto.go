package remotive

import (
	"fmt"
	"html"
	"jobstream/internal/category"
	"jobstream/internal/domain"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// APIResponse represents the top-level JSON response from the Remotive API.
type APIResponse struct {
	JobCount int           `json:"job-count"`
	Jobs     []RemotiveJob `json:"jobs"`
}

// RemotiveJob represents a single job from the Remotive API.
// We map only the fields we care about.
type RemotiveJob struct {
	ID                        int    `json:"id"`
	URL                       string `json:"url"`
	Title                     string `json:"title"`
	CompanyName               string `json:"company_name"`
	CandidateRequiredLocation string `json:"candidate_required_location"`
	Category                  string `json:"category"`
	Salary                    string `json:"salary"`
	Description               string `json:"description"`
	PublicationDate           string `json:"publication_date"`
}

var htmlTagRegex = regexp.MustCompile(`<[^>]*>`)

// stripHTML removes all HTML tags and decodes common HTML entities to return clean plain text.
func stripHTML(src string) string {
	plain := htmlTagRegex.ReplaceAllString(src, "")
	return html.UnescapeString(plain)
}

// toDomain converts a Remotive DTO into our core domain Entity.
// This isolates the rest of our application from Remotive's specific JSON structure.
func (r *RemotiveJob) toDomain() domain.Job {
	// The API returns dates like "2026-05-13T06:46:26" (no timezone)
	// We parse it manually instead of relying on the default RFC3339 unmarshaler
	postedAt, err := time.Parse("2006-01-02T15:04:05", r.PublicationDate)
	if err != nil {
		// Fallback to now if parsing fails
		postedAt = time.Now()
	}

	return domain.Job{
		// Generate a deterministic ID based on the platform and source ID
		ID:       fmt.Sprintf("remotive-%d", r.ID),
		SourceID: strconv.Itoa(r.ID),
		// Platform is set by the JobService using the fetcher's Name()
		Title:       r.Title,
		Company:     r.CompanyName,
		Location:    r.CandidateRequiredLocation,
		Category:    category.Normalize(parseCategoryFromURL(r.Category, r.URL), r.Title),
		Description: stripHTML(r.Description),
		URL:         r.URL,
		Salary:      r.Salary,
		PostedAt:    postedAt,
		CreatedAt:   time.Now(),
	}
}

func parseCategoryFromURL(category string, url string) string {
	if category != "" {
		return category
	}
	prefix := "remotive.com/remote-jobs/"
	idx := strings.Index(url, prefix)
	if idx != -1 {
		remaining := url[idx+len(prefix):]
		parts := strings.Split(remaining, "/")
		if len(parts) > 0 && parts[0] != "" {
			formatted := strings.ReplaceAll(parts[0], "-", " ")
			return strings.Title(formatted)
		}
	}
	return ""
}
