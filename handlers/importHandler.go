package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/julienschmidt/httprouter"
)

var rawDataDirectory = path.Join("data", "raw")

func init() {
	if _, err := os.Stat(rawDataDirectory); os.IsNotExist(err) {
		os.MkdirAll(rawDataDirectory, 0700)
		log.Printf("created directory %s", rawDataDirectory)
	}
}

// ImportHandler - deals with the import of dumps
func ImportHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	handleUpload(w, r)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(512 << 20)

	fhs := r.MultipartForm.File["files[]"]
	for _, fh := range fhs {
		file, err := fh.Open()
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()
		log.Printf("uploaded: %v", fh.Header)
		filename := path.Join(rawDataDirectory, fh.Filename)
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ImportDoneHandler - showing some stats about the imported data
func ImportDoneHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}
