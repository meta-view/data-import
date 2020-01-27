package handlers

import (
	"fmt"
	"log"
	"meta-view-service/assets"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// AssetsHandler - handler for the basic index file
func AssetsHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		asset := fmt.Sprintf("assets%s", ps.ByName("filepath"))
		log.Printf("loading asset: %s", asset)
		data, err := assets.Asset(asset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		SetContentType(w, asset)
		w.Write(data)
	}
}
