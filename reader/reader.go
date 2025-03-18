package reader

import (
	"bufio"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

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
		url := strings.TrimSpace(scanner.Text())
		if url != "" {
			urls <- url
		}
	}
	close(urls)
	if err := scanner.Err(); err != nil {
		log.Errorf("Error reading CSV: %v", err)
		return err
	}
	log.Info("CSV file read successfully")
	return nil
}

