package category_test

import (
	"jobstream/internal/category"
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name         string
		rawCategory  string
		fallbackText string
		expected     string
	}{
		{
			name:        "Known normalized alias (eng)",
			rawCategory: "eng",
			expected:    "Engineering",
		},
		{
			name:        "Known normalized alias (data)",
			rawCategory: "data & analytics",
			expected:    "Data",
		},
		{
			name:        "Known normalized alias (design)",
			rawCategory: "Product Design",
			expected:    "Design",
		},
		{
			name:         "Fallback inference from title (Senior Go Engineer)",
			rawCategory:  "",
			fallbackText: "Senior Go Engineer",
			expected:     "Engineering",
		},
		{
			name:         "Fallback inference from title (Cybersecurity Analyst)",
			rawCategory:  "General",
			fallbackText: "Cybersecurity Analyst",
			expected:     "Security",
		},
		{
			name:         "Fallback inference from title (HR Specialist)",
			rawCategory:  "",
			fallbackText: "HR Specialist",
			expected:     "People",
		},
		{
			name:         "Unknown category and title fallback to Other",
			rawCategory:  "Miscellaneous",
			fallbackText: "Astronaut Candidate",
			expected:     "Other",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := category.Normalize(tt.rawCategory, tt.fallbackText)
			if got != tt.expected {
				t.Errorf("category.Normalize(%q, %q) = %q, expected %q", tt.rawCategory, tt.fallbackText, got, tt.expected)
			}
		})
	}
}
