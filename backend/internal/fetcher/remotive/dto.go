package remotive

import (
	"jobstream/internal/domain"
	"strconv"
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
	ID             int       `json:"id"`
	URL            string    `json:"url"`
	Title          string    `json:"title"`
	CompanyName    string    `json:"company_name"`
	CandidateRequiredLocation string `json:"candidate_required_location"`
	Salary         string    `json:"salary"`
	Description    string    `json:"description"`
	PublicationDate time.Time `json:"publication_date"`
}

// toDomain converts a Remotive DTO into our core domain Entity.
// This isolates the rest of our application from Remotive's specific JSON structure.
func (r *RemotiveJob) toDomain() domain.Job {
	return domain.Job{
		// Note: We leave ID blank. The DB layer will generate a UUID for us.
		SourceID:    strconv.Itoa(r.ID),
		// Platform is set by the JobService using the fetcher's Name()
		Title:       r.Title,
		Company:     r.CompanyName,
		Location:    r.CandidateRequiredLocation,
		Description: r.Description,
		URL:         r.URL,
		Salary:      r.Salary,
		PostedAt:    r.PublicationDate,
		CreatedAt:   time.Now(),
	}
}
