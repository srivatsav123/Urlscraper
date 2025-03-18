package fetcher

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type URLFetcher interface {
	Fetch(url string) ([]byte, error)
}

type HTTPFetcher struct {
	client *http.Client
}

func NewHTTPFetcher() *HTTPFetcher {
	return &HTTPFetcher{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (h *HTTPFetcher) Fetch(url string) ([]byte, error) {
	var data []byte
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		resp, err := h.client.Get(url)
		if err == nil {
			defer resp.Body.Close()
			data, err = io.ReadAll(resp.Body)
			if err == nil {
				return data, nil
			}
		}
		log.Warnf("Retrying (%d/%d) for %s due to error: %v", i+1, maxRetries, url, err)
		time.Sleep(time.Second * 2)
	}
	return nil, errors.New("failed to fetch URL after retries")
}

func (h *HTTPFetcher) Worker(ctx context.Context, urls <-chan string, results chan<- []byte, errors chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		case url, ok := <-urls:
			if !ok {
				return
			}

			data, err := h.Fetch(url)
			if err != nil {
				errors <- err
				continue
			}
			results <- data
		}
	}
}

