package handlers

import (
	"log"
	"meta-view-service/services"
	"meta-view-service/tools"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// IndexHandler - handler for the basic index file
func IndexHandler(plugins map[string]*tools.Plugin, db *services.Database) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		queryValues := r.URL.Query()
		query := make(map[string]interface{})
		log.Printf("queryValues: %s\n", queryValues)
		for k, v := range queryValues {
			query[k] = strings.Join(v, "")
		}
		if query["table"] == nil {
			query["table"] = "testDB"
		}
		log.Printf("render results for %s\n", query)
		results, err := db.ReadEntries(query)
		if err != nil {
			log.Println(err)
			return
		}

		elements := make([]string, 0)
		var element map[string]interface{}
		for id := range results {
			element = results[id].(map[string]interface{})
			provider := element["provider"].(string)
			plugin := plugins[provider]
			render, err := plugin.Present(element, "")
			if err == nil {
				elements = append(elements, render)
			} else {
				log.Println(err)
			}
		}
		renderTemplate(w, "index.html", elements)
	}
}

// UploadFormHandler - renders the upload Form
func UploadFormHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		renderTemplate(w, "uploadForm.html", nil)
	}
}
