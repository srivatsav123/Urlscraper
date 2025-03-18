package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
         
	"Urlscraper/fetcher"
        "Urlscraper/reader"
        "Urlscraper/writer"

)

func main() {
	filePath := flag.String("file", "urls.csv", "Path to CSV file containing URLs")
	maxWorkers := flag.Int("workers", 50, "Maximum number of concurrent downloads")
	timeout := flag.Int("timeout", 10, "HTTP request timeout in seconds")
	flag.Parse()

	urls := make(chan string, *maxWorkers)
	results := make(chan []byte, *maxWorkers)
	errors := make(chan error, *maxWorkers)

	// Initialize fetcher & writer
	fetcherInstance := fetcher.NewHTTPFetcher(time.Duration(*timeout) * time.Second)
	fileWriter := &writer.DiskWriter{}

	// Start reading CSV file
	go func() {
		if err := reader.ReadCSV(*filePath, urls); err != nil {
			fmt.Printf("[ERROR] Failed to read CSV: %v\n", err)
		}
	}()

	// Start workers
	var workerWG sync.WaitGroup
	for i := 0; i < *maxWorkers; i++ {
		workerWG.Add(1)
		go func() {
			fetcher.Worker(fetcherInstance, urls, results, errors)
			workerWG.Done()
		}()
	}

	// Close results channel after workers finish
	go func() {
		workerWG.Wait()
		close(results)
	}()

	// Start writer
	var writerWG sync.WaitGroup
	writerWG.Add(1)
	go writer. ProcessFiles(fileWriter, results, &writerWG)

	// Error handling
	go func() {
		for err := range errors {
			fmt.Println("[ERROR]", err)
		}
	}()

	writerWG.Wait()
	close(errors)

	fmt.Println("[INFO] All files processed successfully!")
}

