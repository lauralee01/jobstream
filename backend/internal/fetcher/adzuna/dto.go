package adzuna

import (
	"fmt"
	"jobstream/internal/category"
	"jobstream/internal/domain"
	"strings"
	"time"
)

type APIResponse struct {
	Results []AdzunaJob `json:"results"`
}

type AdzunaJob struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Company        struct {
		DisplayName string `json:"display_name"`
	} `json:"company"`
	Location       struct {
		Area []string `json:"area"`
	} `json:"location"`
	Category       struct {
		Label string `json:"label"`
	} `json:"category"`
	Description    string   `json:"description"`
	URL    string   `json:"redirect_url"`
	Salary      string  `json:"salary"`
	PublicationDate           string `json:"publication_date"`
}

func (aj AdzunaJob) ToJob() domain.Job {
	postedAt, err := time.Parse("2006-01-02T15:04:05", aj.PublicationDate)
	if err != nil {
		// Fallback to now if parsing fails
		postedAt = time.Now()
	}
	return domain.Job{
		ID:             fmt.Sprintf("adzuna-%s", aj.ID),
		SourceID:       aj.ID,
		Title:          aj.Title,
		Company:        aj.Company.DisplayName,
		Location:       strings.Join(aj.Location.Area, ", "),
		Category:       category.Normalize(aj.Category.Label, aj.Title),
		Description:    aj.Description,
		URL:            aj.URL,
		Salary:         aj.Salary,
		PostedAt:    postedAt,
		CreatedAt:   time.Now(),
	}
}
