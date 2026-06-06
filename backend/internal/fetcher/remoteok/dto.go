package remoteok

import (
	"encoding/xml"
	"fmt"
	"jobstream/internal/category"
	"jobstream/internal/domain"
	"strings"
	"time"
)

type RSSFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel *Channel `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Company     string `xml:"company"`
	Description string `xml:"description"`
	Tags        string `xml:"tags"`
	Location    string `xml:"location"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Link        string `xml:"link"`
}

func (i *Item) toDomain() domain.Job {
	postedAt := time.Now()
	if t, err := time.Parse(time.RFC1123Z, i.PubDate); err == nil {
		postedAt = t
	} else if t, err := time.Parse(time.RFC3339, i.PubDate); err == nil {
		postedAt = t
	}

	tagsList := []string{}
	if i.Tags != "" {
		parts := strings.Split(i.Tags, ",")
		for _, p := range parts {
			tagsList = append(tagsList, strings.TrimSpace(p))
		}
	}

	cat := ""
	if len(tagsList) > 0 {
		cat = tagsList[0]
	}

	title := strings.TrimSpace(i.Title)
	company := strings.TrimSpace(i.Company)

	jobCategory := category.Normalize(cat, title)

	return domain.Job{
		ID:          fmt.Sprintf("remoteok-%s", i.GUID),
		SourceID:    i.GUID,
		Title:       title,
		Company:     company,
		Location:    strings.TrimSpace(i.Location),
		Category:    jobCategory,
		Description: strings.TrimSpace(i.Description),
		URL:         strings.TrimSpace(i.Link),
		PostedAt:    postedAt,
		CreatedAt:   time.Now(),
	}
}