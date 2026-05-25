package greenhouse

import (
	"fmt"
	"jobstream/internal/category"
	"jobstream/internal/domain"
	"log"
	"time"
)

// APIResponse represents Greenhouse API response
type APIResponse struct {
	Jobs []GreenhouseJob `json:"jobs"`
}

// GreenhouseDepartment represents a department within a Greenhouse job
type GreenhouseDepartment struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GreenhouseJob represents a job posting
type GreenhouseJob struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	AbsoluteURL string `json:"absolute_url"`

	Location struct {
		Name string `json:"name"`
	} `json:"location"`

	Departments     []GreenhouseDepartment `json:"departments"`
	Description     string                 `json:"description"`
	UpdatedAt       time.Time              `json:"updated_at"`
	PublicationDate string                 `json:"publication_date"`
}

// toDomain maps Greenhouse job to domain entity
func (j *GreenhouseJob) toDomain(company string) domain.Job {
	postedAt, err := time.Parse("2006-01-02T15:04:05", j.PublicationDate)
	if err != nil {
		postedAt = time.Now()
	}

	log.Println("Found job:", j.Title)

	jobCategory := category.Infer(j.Title)

	if len(j.Departments) > 0 {
		jobCategory = j.Departments[0].Name
	}

	return domain.Job{
		ID:          fmt.Sprintf("greenhouse-%s", fmt.Sprintf("%d", j.ID)),
		SourceID:    fmt.Sprintf("%d", j.ID),
		Title:       j.Title,
		Company:     company,
		Location:    j.Location.Name,
		Category:    jobCategory,
		Description: j.Description,
		URL:         j.AbsoluteURL,
		PostedAt:    postedAt,
		CreatedAt:   time.Now(),
	}
}
