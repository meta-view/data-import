package main

import (
	"fmt"
	"log"
	"meta-view-service/handlers"
	"meta-view-service/services"
	"meta-view-service/tools"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	port = 9000
)

func main() {

	handlers.LoadTemplates()
	db, err := services.CreateDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	plugins, err := tools.LoadPlugins("plugins", db)
	if err != nil {
		log.Fatal(err)
	}
	router := httprouter.New()
	//router.ServeFiles("/assets/*filepath", http.Dir("assets"))
	router.GET("/assets/*filepath", handlers.AssetsHandler())
	router.GET("/", handlers.IndexHandler(plugins, db))
	router.GET("/form", handlers.UploadFormHandler())
	router.POST("/upload", handlers.UploadHandler(plugins))
	router.POST("/import", handlers.ImportHandler(plugins))

	log.Printf("Serving Application on port http://localhost:%d", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
