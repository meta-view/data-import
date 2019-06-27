package handlers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	renderTemplate(w, "index.tmpl", nil)
}

