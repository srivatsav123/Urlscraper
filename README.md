# Urlscraper
# UrlScraper: A High-Performance URL Content Downloader  

## Overview  
UrlScraper is a command-line application written in Golang (version 1.23 or newer) that processes a CSV file containing URLs, downloads their content, and saves each response as a randomly named `.txt` file. The application is designed for efficiency, scalability, and resilience, making it ideal for bulk content retrieval.  

## Features  
✅ **High Performance:** Uses goroutines and channels for efficient concurrency.  
✅ **Scalability:** Processes large CSV files without excessive memory usage.  
✅ **Resilient Error Handling:** Logs failed downloads and continues execution.  
✅ **Graceful Shutdown:** Allows up to 5 seconds for ongoing downloads before exiting.  
✅ **Logging & Metrics:** Tracks processed URLs, success/failure rates, and download times.  
✅ **Unit Tests:** Covers core functionalities, including file reading, downloading, and persistence.  

## How It Works  

The pipeline consists of three stages:  

1. **Stage 1 – CSV File Reader**  
   - Reads URLs from the CSV file line by line using a goroutine.  
   - Sends URLs through a channel for processing.  

2. **Stage 2 – Concurrent Downloaders**  
   - Spawns up to 50 goroutines to download content.  
   - Handles network failures gracefully, logging errors instead of crashing.  

3. **Stage 3 – File Persister**  
   - A single goroutine writes downloaded content to disk.  
   - Generates random filenames with a `.txt` extension.  

## Installation  

Ensure you have **Golang 1.23 or newer** installed. Then, run:  

```sh
git clone https://github.com/yourusername/urlscraper.git  
cd urlscraper  
go mod tidy  
go build -o urlscraper  
./urlscraper -file urls.csv -workers 5
``` 

## Run Test cases  

```sh
go to respective package say ex:reader 
go test -v
```

