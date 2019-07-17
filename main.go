package main

import (
	"fmt"
	"log"
	"meta-view-service/handlers"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	port = 9000
)

func main() {

	handlers.LoadTemplates()

	router := httprouter.New()
	router.ServeFiles("/assets/*filepath", http.Dir("assets"))
	router.GET("/", handlers.IndexHandler)
	router.POST("/import", handlers.ImportHandler)
	log.Printf("Serving Application on port %d", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
