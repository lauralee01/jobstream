package remoteok

import (
	"fmt"
	"jobstream/internal/category"
	"jobstream/internal/domain"
	"time"
)

type RemoteOKJob struct {
	ID          string   `json:"id"`
	Title       string   `json:"position"`
	Company     string   `json:"company"`
	Location    string   `json:"location"`
	URL         string   `json:"url"`
	Description string   `json:"description"` 
	Date        string   `json:"date"`       
	Tags        []string `json:"tags"`
}

func (j *RemoteOKJob) toDomain() domain.Job {
	postedAt := time.Now()
	if t, err := time.Parse(time.RFC3339, j.Date); err == nil {
		postedAt = t
	}

	cat := ""
	if len(j.Tags) > 0 {
		cat = j.Tags[0]
	}
	jobCategory := category.Normalize(cat, j.Title)

	return domain.Job{
		ID:          fmt.Sprintf("remoteok-%s", j.ID),
		SourceID:    j.ID,
		Title:       j.Title,
		Company:     j.Company,
		Location:    j.Location, 
		Category:    jobCategory,
		Description: j.Description,
		URL:         j.URL,
		PostedAt:    postedAt,
		CreatedAt:   time.Now(),
	}
}