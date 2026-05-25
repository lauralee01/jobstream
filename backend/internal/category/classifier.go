package category

import "strings"

// Infer attempts to determine a normalized category
// from a job title.
func Infer(title string) string {

	title = strings.ToLower(title)

	for category, keywords := range categoryKeywords {

		for _, keyword := range keywords {

			if strings.Contains(title, keyword) {
				return category
			}
		}
	}

	return "Other"
}
