package writer

import (
	"fmt"
	"math/rand"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

type FileWriter interface {
	Write(data []byte, filename string) error
}

type DiskWriter struct{}

func NewDiskWriter() *DiskWriter {
	return &DiskWriter{}
}

func (d *DiskWriter) Write(data []byte, filename string) error {
	return os.WriteFile(filename, data, 0644)
}

func ProcessFiles(fileWriter FileWriter, results <-chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range results {
		filename := fmt.Sprintf("output_%d.txt", rand.Intn(1000000))
		err := fileWriter.Write(data, filename)
		if err != nil {
			log.Errorf("Failed to write file %s: %v", filename, err)
		} else {
			log.Infof("Successfully wrote to %s", filename)
		}
	}
}

