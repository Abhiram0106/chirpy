package main

import (
	"log"
	"net/http"
)

func startServer() {

	filepathRoot := "."
	port := "8080"
	serveMuxer := http.NewServeMux()

	server := http.Server{
		Addr:    ":" + port,
		Handler: serveMuxer,
	}

	fileHandler := http.FileServer(http.Dir(filepathRoot))
	serveMuxer.Handle("/", fileHandler)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
