package services

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

// FileStorage - basic struct for the files
type FileStorage struct {
	Root string
}

// CreateFileStorage - creates a basic handle for storring files
func CreateFileStorage(folder string) *FileStorage {
	checkFolder(folder)
	return &FileStorage{
		Root: folder,
	}
}

func checkFolder(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, 0700)
		log.Printf("created directory %s", folder)
	}
}

// SaveFile - saves a given file, returns the filepath
func (fs *FileStorage) SaveFile(file string) (string, error) {
	sourceFileStat, err := os.Stat(file)
	if err != nil {
		return "", err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return "", fmt.Errorf("%s is not a regular file", file)
	}

	source, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer source.Close()

	checksum, err := GetSha1ChecksumOfFile(file)
	if err != nil {
		return "", err
	}

	contentType, err := GetFileContentType(file)
	if err != nil {
		return "", err
	}

	dstFolder := path.Join(fs.Root, contentType)
	checkFolder(dstFolder)
	dstFilename := fmt.Sprintf("%s%s", checksum, path.Ext(file))
	dst := path.Join(dstFolder, dstFilename)

	destination, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	log.Printf("%d bytes written to %s", nBytes, dst)
	return fmt.Sprintf("%s/%s", contentType, dstFilename), nil
}

// DeleteFile - delete a file
func (fs *FileStorage) DeleteFile(file string) error {
	path := path.Join(fs.Root, file)
	var err = os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

// ReadFile - reads the content of a file
func (fs *FileStorage) ReadFile(file string) ([]byte, error) {
	path := path.Join(fs.Root, file)
	log.Printf("loading file %s", path)
	return ioutil.ReadFile(path)
}

// GetContentType - reads the content of a file
func (fs *FileStorage) GetContentType(file string) (string, error) {
	path := path.Join(fs.Root, file)
	return GetFileContentType(path)
}

/*
 * Tools
 */

// GetSha1Checksum - returns the Sha1Checksum of a string
func GetSha1Checksum(content string) string {
	bv := []byte(content)
	h := sha1.New()
	h.Write(bv)
	return hex.EncodeToString(h.Sum(nil))
}

// GetSha1ChecksumOfFile - returns the Sha1Checksum of a file
func GetSha1ChecksumOfFile(file string) (string, error) {

	f, err := os.Open(file)
	if err != nil {
		log.Printf("error opening %s\n", file)
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// GetFileContent - returns the content of a file as string
func GetFileContent(file string) (string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	output := string([]byte(content))
	return output, nil
}

// GetFileBase64 - returns the content of a file as base64
func GetFileBase64(file string) (string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(content))
	return encoded, nil
}

// GetFileContentType - returns the content type of a file
func GetFileContentType(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		log.Printf("error opening %s\n", file)
		return "", err
	}
	defer f.Close()

	buffer := make([]byte, 512)

	_, err = f.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
