package handlers

import (
	"log"
	"meta-view-service/services"
	"meta-view-service/tools"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// IndexHandler - handler for the basic index file
func IndexHandler(plugins map[string]*tools.Plugin, db *services.Database) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		queryValues := r.URL.Query()
		query := make(map[string]interface{})
		log.Printf("queryValues: %s\n", queryValues)
		for k, v := range queryValues {
			query[k] = v
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
		for provider := range results {
			plugin := plugins[provider]
			pluginResults := results[provider].(map[string]interface{})
			renders, err := plugin.Present(pluginResults, "")
			if err == nil {
				elements = append(elements, renders...)
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
