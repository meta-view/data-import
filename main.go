package main

import (
	"fmt"
	"log"
	"meta-view-service/handlers"
	"meta-view-service/services"
	"meta-view-service/tools"
	"net/http"
	"time"

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
	testDB(db)

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

func testDB(db *services.Database) {
	data := make(map[string]interface{})
	data["table"] = "testDB"
	data["date"] = time.Now()
	err := db.InsertEntry(data)
	if err != nil {
		log.Printf("error %v", err)
	}
}
