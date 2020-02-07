package main

import (
	"fmt"
	"log"
	"meta-view-service/handlers"
	"meta-view-service/services"
	"meta-view-service/tools"
	"net/http"
	"os"
	"path"

	"github.com/julienschmidt/httprouter"
)

// VersionString - the version of the application
var VersionString string

const (
	port = 9000
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	handlers.LoadTemplates(VersionString)

	fs := services.CreateFileStorage(path.Join("data", "files"))

	db, err := services.CreateDatabase(path.Join("data", "database"))
	if err != nil {
		return err
	}
	defer db.Close()

	plugins, err := tools.LoadPlugins("plugins", db, fs)
	if err != nil {
		return err
	}

	router := httprouter.New()

	router.GET("/version", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintln(w, VersionString)
	})
	//router.ServeFiles("/assets/*filepath", http.Dir("assets"))
	router.GET("/", handlers.IndexHandler(plugins, db))
	router.GET("/assets/*filepath", handlers.AssetsHandler())
	router.GET("/files/*filepath", handlers.FilesHandler(fs))
	router.GET("/form", handlers.UploadFormHandler())
	router.POST("/upload", handlers.UploadHandler(plugins))
	router.POST("/import", handlers.ImportHandler(plugins))

	log.Printf("Serving Application on port http://localhost:%d", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
