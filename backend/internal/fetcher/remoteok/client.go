package remoteok

import (
    "context"
    "encoding/xml"
    "fmt"
    "io"
    "log"
    "jobstream/internal/domain"
    "net/http"
    "time"
)

const (
    platformName = "RemoteOK"
    baseURL      = "https://remoteok.com/rss"
)

type Client struct {
    httpClient *http.Client
}

func NewClient() *Client {
    return &Client{
        httpClient: &http.Client{
            Timeout: 40 * time.Second,
        },
    }
}

func (c *Client) Name() string { return platformName }

func (c *Client) Fetch(ctx context.Context) ([]domain.Job, error) {
    var lastErr error

    for attempt := 0; attempt < 4; attempt++ {

        req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
        if err != nil {
            return nil, fmt.Errorf("RemoteOK: create request: %w", err)
        }

        // Standard feed reader headers
        req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
        req.Header.Set("Accept", "application/rss+xml, application/xml;q=0.9, */*;q=0.8")

        resp, err := c.httpClient.Do(req)
        if err != nil {
            lastErr = err
            time.Sleep(time.Duration(600+attempt*800) * time.Millisecond)
            continue
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            lastErr = fmt.Errorf("RemoteOK status %d", resp.StatusCode)
            time.Sleep(time.Duration(600+attempt*800) * time.Millisecond)
            continue
        }

        bodyBytes, _ := io.ReadAll(resp.Body)

        // Decode XML
        var feed RSSFeed
        if err := xml.Unmarshal(bodyBytes, &feed); err != nil {
            log.Println("REMOTEOK XML DECODE ERROR:", err)
            // log.Println("REMOTEOK RAW BODY:", string(bodyBytes))
            lastErr = err
            time.Sleep(time.Duration(600+attempt*800) * time.Millisecond)
            continue
        }

        if feed.Channel == nil {
            lastErr = fmt.Errorf("RemoteOK RSS channel is missing")
            time.Sleep(time.Duration(600+attempt*800) * time.Millisecond)
            continue
        }

        jobs := make([]domain.Job, 0, len(feed.Channel.Items))
        for _, item := range feed.Channel.Items {
            jobs = append(jobs, item.toDomain())
        }

        return jobs, nil
    }

    return nil, fmt.Errorf("RemoteOK failed after retries: %w", lastErr)
}
