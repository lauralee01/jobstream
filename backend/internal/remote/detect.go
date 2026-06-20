package remote

import (
	"jobstream/internal/domain"
	"regexp"
	"strings"
)

var remoteOnlyPlatforms = map[string]bool{
	"WeWorkRemotely": true,
	"Remotive":       true,
}

var (
	remoteRegex = regexp.MustCompile(`(?i)\b(remote|remotely|remote-first|fully remote|work from home|wfh|anywhere|worldwide|distributed)\b`)

	notRemoteRegex = regexp.MustCompile(`(?i)\b(not\s+remote|no\s+remote|non-remote|not\s+a\s+remote|remote:\s*no|remote:\s*false|temporary\s+remote|remote\s+not\s+available)\b`)
)

func Detect(job domain.Job) bool {
	if remoteOnlyPlatforms[job.Platform] {
		return true
	}

	location := strings.ToLower(strings.TrimSpace(job.Location))

	if location != "" {
		if notRemoteRegex.MatchString(location) {
			return false
		}

		if remoteRegex.MatchString(location) {
			return true
		}
	}

	text := job.Location + " " + job.Title + " " + job.Description

	if notRemoteRegex.MatchString(text) {
		return false
	}

	return remoteRegex.MatchString(text)
}
