package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient abstracts HTTP GET requests for testing and customization.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// DefaultHTTPClient creates an http.Client configured with reasonable timeouts.
func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

// FetchHTML executes a GET request using the provided client, sets standard headers,
// and returns the read closer response body.
func FetchHTML(ctx context.Context, client HTTPClient, targetURL string) (io.ReadCloser, error) {
	if client == nil {
		client = DefaultHTTPClient()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", targetURL, err)
	}

	// Browser-like user agent to prevent bot rejection on public listings
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 EventBeep/1.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL %s: %w", targetURL, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected HTTP status %d fetching %s", resp.StatusCode, targetURL)
	}

	return resp.Body, nil
}
