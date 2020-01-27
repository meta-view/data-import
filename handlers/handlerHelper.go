package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
)

// SetContentType - set a specific contentType
func SetContentType(w http.ResponseWriter, asset string) {
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
