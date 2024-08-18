package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func startServer() {

	mux := http.NewServeMux()
	cfg := apiConfig{
		fileserverHits: 0,
	}

	fileHandler := http.FileServer(http.Dir(filepathRoot))
	mux.Handle(appPath, http.StripPrefix(stripAppPath, cfg.middlewareMetricsInc(fileHandler)))

	mux.HandleFunc(http.MethodGet+" "+healthPath, chirpyHandlerFunc)
	mux.HandleFunc(http.MethodGet+" "+metricsPath, cfg.getFileserverHits)
	mux.HandleFunc(resetMetricsPath, cfg.resetFileserverHits)
	mux.HandleFunc(http.MethodGet+" "+chirpsPath, getChirps)
	mux.HandleFunc(http.MethodPost+" "+chirpsPath, postChirp)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
