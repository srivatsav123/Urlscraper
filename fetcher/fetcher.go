package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// URLFetcher defines an interface for fetching data from a URL.
type URLFetcher interface {
	Fetch(url string) ([]byte, error)
}

// HTTPFetcher implements URLFetcher using the net/http package.
type HTTPFetcher struct {
	client *http.Client
}

// NewHTTPFetcher creates an HTTPFetcher with a default timeout.
func NewHTTPFetcher(timeout time.Duration) *HTTPFetcher {
	return &HTTPFetcher{
		client: &http.Client{Timeout: timeout},
	}
}

// Fetch retrieves data from the given URL with retry logic.
func (h *HTTPFetcher) Fetch(url string) ([]byte, error) {
	const maxRetries = 3
	var data []byte
	var err error

	for i := 0; i < maxRetries; i++ {
		resp, err := h.client.Get(url)
		if err == nil {
			defer resp.Body.Close()
			data, err = io.ReadAll(resp.Body)
			if err == nil {
				return data, nil
			}
		}
		fmt.Printf("[WARN] Retrying (%d/%d) for %s due to error: %v\n", i+1, maxRetries, url, err)
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, err)
}

// Worker fetches URLs concurrently and sends results to a channel.
func Worker(fetcher URLFetcher, urls <-chan string, results chan<- []byte, errors chan<- error) {
	for url := range urls {
		data, err := fetcher.Fetch(url)
		if err != nil {
			errors <- fmt.Errorf("failed to download %s: %v", url, err)
			continue
		}
		results <- data
	}
}

