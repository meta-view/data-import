package handlers

import (
	"log"
	"meta-view-service/tools"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// IndexHandler - handler for the basic index file
func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	testPlugin()
	renderTemplate(w, "index.tmpl", nil)
}

func testPlugin() {
	testPlugin, err := tools.LoadPlugin("test/plugin")
	if err != nil {
		log.Fatal(err)
	}
	accountName, err := testPlugin.GetAccountName()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("account: %s", accountName)
}
