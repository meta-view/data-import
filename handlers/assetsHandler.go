package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"meta-view-service/assets"
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
		setContentType(w, asset)
		w.Write(data)
	}
}

func setContentType(w http.ResponseWriter, asset string) {
	ext := filepath.Ext(asset)
	switch ext {
	case ".png":
		fallthrough
	case ".gif":
		w.Header().Set("Content-Type", fmt.Sprintf("image/%s", ext))
	case ".woff":
		fallthrough
	case ".woff2":
		fallthrough
	case ".eot":
		fallthrough
	case ".ttf":
		w.Header().Set("Content-Type", fmt.Sprintf("font/%s", ext))
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "text/javascript")
	default:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
}
