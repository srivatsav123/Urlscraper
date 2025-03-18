package writer

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
)

// FileWriter defines an interface for writing data to a file.
type FileWriter interface {
	Write(data []byte, filename string) error
}

// DiskWriter implements FileWriter by writing data to disk.
type DiskWriter struct{}

// Write writes data to a specified file.
func (d *DiskWriter) Write(data []byte, filename string) error {
	return os.WriteFile(filename, data, 0644)
}

// FileWriter processes results and writes them to files.
func ProcessFiles(fileWriter FileWriter, results <-chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range results {
		filename := fmt.Sprintf("output_%d.txt", rand.Int())
		if err := fileWriter.Write(data, filename); err != nil {
			fmt.Printf("[ERROR] Failed to write file: %v\n", err)
		}
	}
}

