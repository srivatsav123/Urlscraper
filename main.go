package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
	"Urlscraper/fetcher"
	"Urlscraper/reader"
	"Urlscraper/writer"
)

func main() {
	filePath := flag.String("file", "urls.csv", "Path to CSV file containing URLs")
	maxWorkers := flag.Int("workers", 50, "Maximum number of concurrent downloads")
	flag.Parse()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Warn("Received termination signal. Gracefully shutting down...")
		cancel()
	}()

	urls := make(chan string, *maxWorkers)
	results := make(chan []byte, *maxWorkers)
	errors := make(chan error, *maxWorkers)

	var wg sync.WaitGroup
	fetcher := fetcher.NewHTTPFetcher()
	fileWriter := writer.NewDiskWriter()

	// Start reading CSV file
	go func() {
		if err := reader.ReadCSV(*filePath, urls); err != nil {
			log.Errorf("Failed to read CSV: %v", err)
		}
	}()

	// Start worker goroutines
	var workerWG sync.WaitGroup
	for i := 0; i < *maxWorkers; i++ {
		workerWG.Add(1)
		go func() {
			defer workerWG.Done()
			fetcher.Worker(ctx, urls, results, errors)
		}()
	}

	// Close results channel after workers finish
	go func() {
		workerWG.Wait()
		close(results)
	}()

	// Start file writer
	wg.Add(1)
	go writer.ProcessFiles(fileWriter, results, &wg)

	// Error handling
	go func() {
		for err := range errors {
			log.Error(err)
		}
	}()

	wg.Wait()
	close(errors)

	log.Info("All files processed successfully!")
}

