package services

import (
	"log"
	"os"
	"path"
)

var rootDataDirectory = path.Join("data", "files")

// FileStorage - basic struct for the files
type FileStorage struct {
	Root string
}

// CreateFileStorage - creates a basic handle for storring files
func CreateFileStorage() (*FileStorage, error) {
	checkFolder(rootDataDirectory)
	return &FileStorage{
		Root: rootDataDirectory,
	}, nil
}

func checkFolder(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, 0700)
		log.Printf("created directory %s", folder)
	}
}

func (fs *FileStorage) saveFile(folder string, file string) (string, error) {
	return "", nil
}
