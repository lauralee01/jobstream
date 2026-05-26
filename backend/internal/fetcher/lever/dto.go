package lever

import (
	"fmt"
	"jobstream/internal/category"
	"jobstream/internal/domain"
	"time"
)

type LeverJob struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	HostedURL string `json:"hostedUrl"`

	Categories struct {
		Team     string `json:"team"`
		Location string `json:"location"`
	} `json:"categories"`

	DescriptionPlain string `json:"descriptionPlain"`

	CreatedAt int64 `json:"createdAt"`
}

func (j *LeverJob) toDomain(company string) domain.Job {

	postedAt := time.UnixMilli(j.CreatedAt)

	jobCategory := j.Categories.Team

	if jobCategory == "" {
		jobCategory = category.Infer(j.Text)
	}

	return domain.Job{
		ID:          fmt.Sprintf("lever-%s", j.ID),
		SourceID:    j.ID,
		Title:       j.Text,
		Company:     company,
		Location:    j.Categories.Location,
		Category:    jobCategory,
		Description: j.DescriptionPlain,
		URL:         j.HostedURL,
		PostedAt:    postedAt,
		CreatedAt:   time.Now(),
	}
}
