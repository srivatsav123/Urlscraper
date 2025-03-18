package writer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockWriter struct{}

func (m *MockWriter) Write(data []byte, filename string) error {
	return os.WriteFile(filename, data, 0644)
}

func TestDiskWriter_Write(t *testing.T) {
	fileWriter := &DiskWriter{}
	filename := "test_output.txt"
	data := []byte("test data")

	err := fileWriter.Write(data, filename)
	assert.NoError(t, err)

	writtenData, err := os.ReadFile(filename)
	assert.NoError(t, err)
	assert.Equal(t, data, writtenData)

	os.Remove(filename)
}
