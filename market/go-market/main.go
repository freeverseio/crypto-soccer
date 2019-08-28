package main

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func Index(w http.ResponseWriter, r *http.Request) {
	log.Debug("ciao")
}

func main() {
	log.Info("Starting ...")

	router := mux.NewRouter()

	router.HandleFunc("/calcs", Index)

	log.Fatal(http.ListenAndServe(":8080", router))

	log.Info("... exiting")
}
