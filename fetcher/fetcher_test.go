package fetcher

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFetchSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	}))
	defer ts.Close()

	fetcher := NewHTTPFetcher()
	data, err := fetcher.Fetch(ts.URL)

	assert.NoError(t, err)
	assert.Equal(t, "Hello, world!", string(data))
}

func TestFetchFailure(t *testing.T) {
	fetcher := NewHTTPFetcher()
	_, err := fetcher.Fetch("http://invalid.url")

	assert.Error(t, err)
}

func TestFetchTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer ts.Close()

	fetcher := NewHTTPFetcher()
	fetcher.client.Timeout = 1 * time.Second

	_, err := fetcher.Fetch(ts.URL)
	assert.Error(t, err)
}

