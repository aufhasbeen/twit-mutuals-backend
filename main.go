package main

import (
	// system
	"log"
	"net/http"
	"time"

	"github.com/aufheben/mutuals-server/local/routing"
	"github.com/aufheben/mutuals-server/local/twitterapi"
)

func response(w http.ResponseWriter, r *http.Request) {

}

func main() {
	twitterapi.Init()

	router := routing.Init()

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
