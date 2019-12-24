package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"meta-view-service/tools"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var zipDataDirectory = path.Join("data", "zip")
var rawDataDirectory = path.Join("data", "raw")

func init() {
	folders := []string{zipDataDirectory, rawDataDirectory}
	for _, folder := range folders {
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			os.MkdirAll(folder, 0700)
			log.Printf("created directory %s", folder)
		}
	}
}

// ImportHandler - deals with the import of dumps
func ImportHandler(plugins map[string]*tools.Plugin) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		handleUpload(w, r, plugins)
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request, plugins map[string]*tools.Plugin) {

	r.ParseMultipartForm(512 << 20)

	fhs := r.MultipartForm.File["files[]"]
	for _, fh := range fhs {
		file, err := fh.Open()
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()
		log.Printf("uploaded: %v\n", fh.Header)
		filename := path.Join(zipDataDirectory, fh.Filename)
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		dest := strings.Replace(strings.Replace(filename, zipDataDirectory, rawDataDirectory, 1), ".zip", "", 1)
		files, err := importData(filename, dest)

		for _, file := range files {
			log.Printf("extracting file %s\n", file)
		}

		markers := make(map[string]float64)
		for _, plugin := range plugins {
			value, err := plugin.Detect(dest)
			if err != nil {
				log.Printf("Error: %s\n", err)
			} else {
				markers[plugin.Provider.Name] = value
			}
		}

		for _, plugin := range plugins {
			err := plugin.Import(dest)
			if err != nil {
				log.Printf("Error: %s\n", err)
			}
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func importData(src string, dest string) ([]string, error) {
	log.Printf("importing %s to %s\n", src, dest)
	var filenames []string

	if strings.HasSuffix(src, ".zip") {

		r, err := zip.OpenReader(src)
		if err != nil {
			return filenames, err
		}
		defer r.Close()
		for _, f := range r.File {

			// Store filename/path for returning and using later on
			fpath := filepath.Join(dest, f.Name)

			// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
			if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
				return filenames, fmt.Errorf("%s: illegal file path", fpath)
			}

			filenames = append(filenames, fpath)

			if f.FileInfo().IsDir() {
				// Make Folder
				os.MkdirAll(fpath, os.ModePerm)
				continue
			}

			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return filenames, err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return filenames, err
			}

			rc, err := f.Open()
			if err != nil {
				return filenames, err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()
			rc.Close()

			if err != nil {
				return filenames, err
			}
		}
		return filenames, nil
	}
	return filenames, nil
}

// ImportDoneHandler - showing some stats about the imported data
func ImportDoneHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}
