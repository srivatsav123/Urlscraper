package reader

import (
	"bufio"
	"os"
)

// ReadCSV reads a CSV file and sends URLs to a channel.
func ReadCSV(filePath string, urls chan<- string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	firstLine := true
	for scanner.Scan() {
		if firstLine {
			firstLine = false
			continue
		}
		urls <- scanner.Text()
	}
	close(urls)

	return scanner.Err()
}

