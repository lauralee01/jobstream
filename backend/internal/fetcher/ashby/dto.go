package ashby

import (
	"fmt"
	"jobstream/internal/category"
	"jobstream/internal/domain"
	"time"
)

type APIResponse struct {
	APIVersion string     `json:"apiVersion"`
	Jobs       []AshbyJob `json:"jobs"`
}

type AshbyJob struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Location         string `json:"location"`
	Department       string `json:"department"`
	Team             string `json:"team"`
	IsListed         bool   `json:"isListed"`
	IsRemote         bool   `json:"isRemote"`
	WorkplaceType    string `json:"workplaceType"`
	DescriptionHTML  string `json:"descriptionHtml"`
	DescriptionPlain string `json:"descriptionPlain"`
	PublishedAt      string `json:"publishedAt"`
	EmploymentType   string `json:"employmentType"`
	JobURL           string `json:"jobUrl"`
	ApplyURL         string `json:"applyUrl"`
}

func (j *AshbyJob) toDomain(company string) domain.Job {
	postedAt := time.Now()

	if j.PublishedAt != "" {
		if t, err := time.Parse(time.RFC3339Nano, j.PublishedAt); err == nil {
			postedAt = t
		}
	}

	rawCategory := j.Department
	if rawCategory == "" {
		rawCategory = j.Team
	}

	jobCategory := category.Normalize(rawCategory, j.Title)

	sourceID := j.ID
	if sourceID == "" {
		sourceID = j.JobURL
	}

	return domain.Job{
		ID:          fmt.Sprintf("ashby-%s", sourceID),
		SourceID:    sourceID,
		Title:       j.Title,
		Company:     company,
		Location:    j.Location,
		Category:    jobCategory,
		Description: j.DescriptionPlain,
		URL:         j.JobURL,
		PostedAt:    postedAt,
		CreatedAt:   time.Now(),
	}
}
