package main

import (
	"log"
	"net/http"
	"rack/pkg"
)

func main() {
	server := &http.Server{
		Handler: pkg.Router(),
		Addr:    ":8080",
	}
	log.Fatal(server.ListenAndServe())
}
