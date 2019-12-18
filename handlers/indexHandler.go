package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// IndexHandler - handler for the basic index file
func IndexHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		renderTemplate(w, "index.tmpl", nil)
	}
}
