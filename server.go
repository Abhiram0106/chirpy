package main

import (
	"fmt"
	"net/http"
)

func startServer() {

	serveMuxer := http.NewServeMux()
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: serveMuxer,
	}

	err := server.ListenAndServe()

	if err != nil {
		fmt.Println(err)
	}
}
