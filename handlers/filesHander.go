package handlers

import (
	"fmt"
	"meta-view-service/services"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// FilesHandler - handler for serving files
func FilesHandler(fs *services.FileStorage) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		file := ps.ByName("filepath")
		contentType, err := fs.GetContentType(file)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("error getting contentType %s", err.Error()),
				http.StatusInternalServerError)
		}

		data, err := fs.ReadFile(file)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("error loading image %s", err.Error()),
				http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", contentType)
		w.Write(data)
	}
}
