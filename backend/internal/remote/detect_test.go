package remote_test

import (
	"jobstream/internal/domain"
	"jobstream/internal/remote"
	"testing"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		name     string
		job      domain.Job
		expected bool
	}{
		{
			name: "Remote platform (WeWorkRemotely)",
			job: domain.Job{
				Platform: "WeWorkRemotely",
				Title:    "Software Engineer",
			},
			expected: true,
		},
		{
			name: "Remote platform (Remotive)",
			job: domain.Job{
				Platform: "Remotive",
				Title:    "Backend Developer",
			},
			expected: true,
		},
		{
			name: "Location explicitly Remote",
			job: domain.Job{
				Platform: "Greenhouse",
				Title:    "Product Designer",
				Location: "Remote, US",
			},
			expected: true,
		},
		{
			name: "Strong negative text signal (strictly onsite)",
			job: domain.Job{
				Platform:    "Greenhouse",
				Title:       "Remote Developer",
				Location:    "New York, NY",
				Description: "This role is strictly onsite in our NY office.",
			},
			expected: false,
		},
		{
			name: "Hybrid location signal",
			job: domain.Job{
				Platform: "Lever",
				Title:    "Senior Engineer",
				Location: "Hybrid - Austin, TX",
			},
			expected: false,
		},
		{
			name: "100% Remote in description",
			job: domain.Job{
				Platform:    "Ashby",
				Title:       "Data Scientist",
				Location:    "San Francisco, CA",
				Description: "We are looking for a teammate to work 100% remote anywhere in the US.",
			},
			expected: true,
		},
		{
			name: "Onsite location signal",
			job: domain.Job{
				Platform: "Workable",
				Title:    "Office Manager",
				Location: "Onsite Chicago",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := remote.Detect(tt.job)
			if got != tt.expected {
				t.Errorf("remote.Detect(%+v) = %v, expected %v", tt.job.Title, got, tt.expected)
			}
		})
	}
}
