package main

import (
	"log"
	"net/http"
)

func startServer() {

	filepathRoot := "."
	healthPath := "/healthz"
	port := "8080"
	mux := http.NewServeMux()

	fileHandler := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/*", http.StripPrefix("/app", fileHandler))
	mux.HandleFunc(healthPath, chirpyHandlerFunc)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func chirpyHandlerFunc(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set(http.CanonicalHeaderKey("content-type"), "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	huh, yeah := writer.Write([]byte("OK"))
	log.Println(huh)
	if yeah != nil {
		log.Println("error = " + yeah.Error())
	}
}
