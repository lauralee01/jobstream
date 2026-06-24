package workable

import (
	"fmt"
	"jobstream/internal/category"
	"jobstream/internal/domain"
	"strings"
	"time"
)

type APIResponse struct {
	Name        string        `json:"name"`
	Subdomain   string        `json:"subdomain"`
	Description string        `json:"description"`
	Website     string        `json:"website"`
	Jobs        []WorkableJob `json:"jobs"`
}

type WorkableJob struct {
	ID              string `json:"id"`
	Shortcode       string `json:"shortcode"`
	Title           string `json:"title"`
	FullTitle       string `json:"full_title"`
	Department      string `json:"department"`
	URL             string `json:"url"`
	ApplicationURL  string `json:"application_url"`
	Shortlink       string `json:"shortlink"`
	EmploymentType  string `json:"employment_type"`
	Description     string `json:"description"`
	DescriptionHTML string `json:"description_html"`
	PublishedOn     string `json:"published_on"`
	CreatedAt       string `json:"created_at"`

	Location WorkableLocation `json:"location"`
}

type WorkableLocation struct {
	City        string `json:"city"`
	Region      string `json:"region"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Telecommute bool   `json:"telecommute"`
}

func (j *WorkableJob) toDomain(company string) domain.Job {
	title := strings.TrimSpace(j.Title)
	if title == "" {
		title = strings.TrimSpace(j.FullTitle)
	}

	sourceID := j.Shortcode
	if sourceID == "" {
		sourceID = j.ID
	}
	if sourceID == "" {
		sourceID = j.URL
	}

	jobURL := strings.TrimSpace(j.URL)
	if jobURL == "" {
		jobURL = strings.TrimSpace(j.Shortlink)
	}
	if jobURL == "" {
		jobURL = strings.TrimSpace(j.ApplicationURL)
	}

	description := strings.TrimSpace(j.Description)
	if description == "" {
		description = strings.TrimSpace(j.DescriptionHTML)
	}

	postedAt := time.Now()
	if j.PublishedOn != "" {
		if t, err := time.Parse("2006-01-02", j.PublishedOn); err == nil {
			postedAt = t
		}
	} else if j.CreatedAt != "" {
		if t, err := time.Parse(time.RFC3339, j.CreatedAt); err == nil {
			postedAt = t
		}
	}

	location := formatLocation(j.Location)

	rawCategory := j.Department
	jobCategory := category.Normalize(rawCategory, title)

	return domain.Job{
		ID:          fmt.Sprintf("workable-%s", sourceID),
		SourceID:    sourceID,
		Title:       title,
		Company:     company,
		Location:    location,
		Category:    jobCategory,
		Description: description,
		URL:         jobURL,
		PostedAt:    postedAt,
		CreatedAt:   time.Now(),
	}
}

func formatLocation(location WorkableLocation) string {
	if location.Telecommute {
		return "Remote"
	}

	parts := []string{}

	if strings.TrimSpace(location.City) != "" {
		parts = append(parts, strings.TrimSpace(location.City))
	}

	if strings.TrimSpace(location.Region) != "" {
		parts = append(parts, strings.TrimSpace(location.Region))
	}

	if strings.TrimSpace(location.Country) != "" {
		parts = append(parts, strings.TrimSpace(location.Country))
	}

	if len(parts) == 0 {
		return "Not specified"
	}

	return strings.Join(parts, ", ")
}
