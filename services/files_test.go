package services

import (
	"log"
	"path"
	"testing"
)

func TestSaveTextFile(t *testing.T) {
	fs := CreateFileStorage(path.Join("..", "data", "files"))

	savedFile, err := fs.SaveFile(path.Join("..", "build.sh"))
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("saved file to %s", savedFile)

	fs.DeleteFile(savedFile)
}

func TestSaveBinarytFile(t *testing.T) {
	fs := CreateFileStorage(path.Join("..", "data", "files"))

	savedFile, err := fs.SaveFile(path.Join("..", "logo.png"))
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("saved file to %s", savedFile)

	fs.DeleteFile(savedFile)
}
