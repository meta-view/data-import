package main

import (
	"fmt"
	"log"
	"meta-view-service/handlers"
	"meta-view-service/tools"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	port = 9000
)

func main() {

	handlers.LoadTemplates()
	plugins, err := tools.LoadPlugins("plugins")
	if err != nil {
		log.Fatal(err)
	}
	router := httprouter.New()
	router.ServeFiles("/assets/*filepath", http.Dir("assets"))
	router.GET("/", handlers.IndexHandler())
	router.POST("/import", handlers.ImportHandler(plugins))
	log.Printf("Serving Application on port %d", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
