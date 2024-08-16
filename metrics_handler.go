package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		log.Printf("tracking method %s, at %s hits incremented to %d\n", r.Method, r.URL.Path, cfg.fileserverHits)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getFileserverHits(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add(http.CanonicalHeaderKey("content-type"), "text/html; charset=utf-8")
	htmlText, err := os.ReadFile("metrics.html")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Error reading file"))
	} else {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(fmt.Sprintf(string(htmlText), cfg.fileserverHits)))
	}
}

func (cfg *apiConfig) resetFileserverHits(writer http.ResponseWriter, request *http.Request) {
	cfg.fileserverHits = 0
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Hits reset to 0"))
}
