package main

import (
	"log"
	"net/http"
)

func chirpyHandlerFunc(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add(http.CanonicalHeaderKey("content-type"), "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	_, err := writer.Write([]byte(http.StatusText(http.StatusOK)))
	if err != nil {
		log.Println(err.Error())
	}
}
