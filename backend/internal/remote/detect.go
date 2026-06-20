package remote

import (
	"jobstream/internal/domain"
	"strings"
)

var remoteOnlyPlatforms = map[string]bool{
	"WeWorkRemotely": true,
	"Remotive":       true,
}

func Detect(job domain.Job) bool {
	if remoteOnlyPlatforms[job.Platform] {
		return true
	}

	text := strings.ToLower(
		job.Location + " " +
			job.Title + " " +
			job.Description,
	)

	remoteKeywords := []string{
		"remote",
		"remote-first",
		"fully remote",
		"work from home",
		"wfh",
		"anywhere",
		"distributed",
	}

	for _, keyword := range remoteKeywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}

	return false
}
