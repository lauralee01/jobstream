package remoteok

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "jobstream/internal/domain"
    "net/http"
    "time"
)

const (
    platformName = "RemoteOK"
    baseURL      = "https://remoteok.com/api"
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

        // FULL browser headers (required)
        req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
        req.Header.Set("Accept", "application/json, text/plain, */*")
        req.Header.Set("Accept-Language", "en-US,en;q=0.9")
        req.Header.Set("Referer", "https://remoteok.com/")
        req.Header.Set("Origin", "https://remoteok.com")
        req.Header.Set("X-Requested-With", "XMLHttpRequest")
        req.Header.Set("Sec-Fetch-Site", "same-origin")
        req.Header.Set("Sec-Fetch-Mode", "cors")
        req.Header.Set("Sec-Fetch-Dest", "empty")
        req.Header.Set("Cache-Control", "no-cache")

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

        // Detect HTML fallback (Cloudflare or error page)
        if bytes.Contains(bodyBytes, []byte("<html")) {
            log.Println("REMOTEOK ERROR: Received HTML instead of JSON")
            log.Println("REMOTEOK RAW BODY:", string(bodyBytes))
            lastErr = fmt.Errorf("RemoteOK returned HTML instead of JSON")
            time.Sleep(time.Duration(600+attempt*800) * time.Millisecond)
            continue
        }

        // Decode JSON
        var apiJobs []RemoteOKJob
        if err := json.Unmarshal(bodyBytes, &apiJobs); err != nil {
            log.Println("REMOTEOK JSON DECODE ERROR:", err)
            log.Println("REMOTEOK RAW BODY:", string(bodyBytes))
            lastErr = err
            time.Sleep(time.Duration(600+attempt*800) * time.Millisecond)
            continue
        }

        if len(apiJobs) > 0 && apiJobs[0].ID == "" {
            apiJobs = apiJobs[1:]
        }

        jobs := make([]domain.Job, 0, len(apiJobs))
        for _, j := range apiJobs {
            jobs = append(jobs, j.toDomain())
        }

        return jobs, nil
    }

    return nil, fmt.Errorf("RemoteOK failed after retries: %w", lastErr)
}
