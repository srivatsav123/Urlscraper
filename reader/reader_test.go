package reader

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadCSV_ValidFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testfile-*.csv")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("url\nhttps://example.com\nhttps://test.com\n")
	assert.NoError(t, err)
	tempFile.Close()

	urls := make(chan string, 2)
	err = ReadCSV(tempFile.Name(), urls) // Updated function name
	assert.NoError(t, err)

	expected := []string{"https://example.com", "https://test.com"}
	var result []string
	for url := range urls {
		result = append(result, url)
	}

	assert.Equal(t, expected, result)
}

func TestReadCSV_EmptyFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testfile-*.csv")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	urls := make(chan string, 1)
	err = ReadCSV(tempFile.Name(), urls) // Updated function name
	assert.NoError(t, err)
}

func TestReadCSV_FileNotFound(t *testing.T) {
	urls := make(chan string, 1)
	err := ReadCSV("nonexistent.csv", urls) // Updated function name
	assert.Error(t, err)
}

