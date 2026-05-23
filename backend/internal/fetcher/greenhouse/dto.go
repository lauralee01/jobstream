package greenhouse

import (
	"fmt"
	"jobstream/internal/domain"
	"time"
)

// APIResponse represents Greenhouse API response
type APIResponse struct {
	Jobs []GreenhouseJob `json:"jobs"`
}

// GreenhouseJob represents a job posting
type GreenhouseJob struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	AbsoluteURL string `json:"absolute_url"`

	Location struct {
		Name string `json:"name"`
	} `json:"location"`

	Description     string    `json:"description"`
	UpdatedAt       time.Time `json:"updated_at"`
	PublicationDate string    `json:"publication_date"`
}

// toDomain maps Greenhouse job to domain entity
func (j *GreenhouseJob) toDomain(company string) domain.Job {
	postedAt, err := time.Parse("2006-01-02T15:04:05", j.PublicationDate)
	if err != nil {
		postedAt = time.Now()
	}

	category := ""

	return domain.Job{
		ID:          fmt.Sprintf("greenhouse-%s", fmt.Sprintf("%d", j.ID)),
		SourceID:    fmt.Sprintf("%d", j.ID),
		Title:       j.Title,
		Company:     company,
		Location:    j.Location.Name,
		Category:    category,
		Description: j.Description,
		URL:         j.AbsoluteURL,
		PostedAt:    postedAt,
		CreatedAt:   time.Now(),
	}
}
