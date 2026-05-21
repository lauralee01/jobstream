package weworkremotely

import (
	"encoding/xml"
	"fmt"
	"jobstream/internal/domain"
	"strings"
	"time"
)

// RSS represents the top-level RSS feed structure
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// Channel represents the channel element in the RSS feed
type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Items   []Item   `xml:"item"`
}

// Item represents a single job posting item in the RSS feed
type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	Guid        string   `xml:"guid"`
	Region      string   `xml:"region"`
	Category    string   `xml:"category"`
}

// toDomain converts the XML RSS Item into our core domain Job model
func (item *Item) toDomain() domain.Job {
	// Parse company name and job title from the Title element ("Company: Title")
	parts := strings.SplitN(item.Title, ": ", 2)
	company := "Unknown"
	title := item.Title
	if len(parts) == 2 {
		company = strings.TrimSpace(parts[0])
		title = strings.TrimSpace(parts[1])
	}

	// Parse PostedAt publication date
	postedAt := time.Now()
	if item.PubDate != "" {
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			postedAt = t
		} else if t, err := time.Parse(time.RFC1123, item.PubDate); err == nil {
			postedAt = t
		}
	}

	location := strings.TrimSpace(item.Region)
	if location == "" {
		location = "Remote"
	}

	// Get unique SourceID from link or guid
	sourceID := item.Guid
	if sourceID == "" {
		sourceID = item.Link
	}
	
	// Clean up the URL to extract a clean slug if possible
	if strings.Contains(sourceID, "weworkremotely.com/remote-jobs/") {
		urlParts := strings.Split(sourceID, "weworkremotely.com/remote-jobs/")
		if len(urlParts) == 2 {
			sourceID = urlParts[1]
		}
	}

	return domain.Job{
		ID:          fmt.Sprintf("weworkremotely-%s", sourceID),
		SourceID:    sourceID,
		Title:       title,
		Company:     company,
		Location:    location,
		Category:    strings.TrimSpace(item.Category),
		Description: item.Description,
		URL:         item.Link,
		PostedAt:    postedAt,
		CreatedAt:   time.Now(),
	}
}
