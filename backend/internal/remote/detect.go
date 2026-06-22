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
	// Strong evidence that the role is NOT remote.
	strongNegativeTextRegex = regexp.MustCompile(
		`(?i)\b(not\s+(?:a\s+)?remote|no\s+remote|non[-\s]?remote|remote\s+not\s+available|does\s+not\s+offer.{0,30}remote|doesn't\s+offer.{0,30}remote|not\s+offer.{0,30}remote|no\s+hybrid\s+or\s+remote|onsite\s+role|on-site\s+role|must\s+be\s+onsite|must\s+work\s+onsite|strictly\s+onsite|100%\s+in[-\s]?office)\b`,
	)

	// If these appear in location/title, it's likely not fully remote.
	negativeLocationTitleRegex = regexp.MustCompile(
		`(?i)\b(onsite|on-site|in[-\s]?office|office[-\s]?based|hybrid|partly\s+remote|partially\s+remote)\b`,
	)

	// Positive remote signals.
	positiveRemoteRegex = regexp.MustCompile(
		`(?i)\b(remote|remotely|remote-first|fully\s+remote|100%\s+remote|work\s+from\s+home|wfh|anywhere|worldwide|distributed)\b`,
	)
)

func Detect(job domain.Job) bool {
	// 1. Remote-only providers
	if remoteOnlyPlatforms[job.Platform] {
		return true
	}

	location := strings.ToLower(strings.TrimSpace(job.Location))
	title := strings.ToLower(strings.TrimSpace(job.Title))
	description := strings.ToLower(strings.TrimSpace(job.Description))

	text := location + " " + title + " " + description

	// 2. Strong negative signals always win
	if strongNegativeTextRegex.MatchString(text) {
		return false
	}

	// 3. Explicit remote location is highly trustworthy
	if positiveRemoteRegex.MatchString(location) {
		return true
	}

	// 4. Explicit onsite/hybrid location or title
	if negativeLocationTitleRegex.MatchString(location) ||
		negativeLocationTitleRegex.MatchString(title) {
		return false
	}

	// 5. Fallback to scanning all text
	return positiveRemoteRegex.MatchString(text)
}
