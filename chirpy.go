package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
	jwtSecret      string
}

func startServer(cfg *apiConfig) {

	mux := http.NewServeMux()

	fileHandler := http.FileServer(http.Dir(filepathRoot))
	mux.Handle(appPath, http.StripPrefix(stripAppPath, cfg.middlewareMetricsInc(fileHandler)))

	mux.HandleFunc(http.MethodGet+" "+healthPath, chirpyHandlerFunc)
	mux.HandleFunc(http.MethodGet+" "+metricsPath, cfg.getFileserverHits)
	mux.HandleFunc(resetMetricsPath, cfg.resetFileserverHits)
	mux.HandleFunc(http.MethodGet+" "+chirpsPath, getChirps)
	mux.Handle(http.MethodPost+" "+chirpsPath, &PostChirpHandler{jwtSecret: cfg.jwtSecret})
	mux.HandleFunc(http.MethodGet+" "+chirpByIDPath, getChirps)
	mux.Handle(http.MethodDelete+" "+chirpByIDPath, &DeleteChirpHandler{jwtSecret: cfg.jwtSecret})
	mux.HandleFunc(http.MethodPost+" "+usersPath, postUsers)
	mux.Handle(http.MethodPost+" "+loginPath, &PostLoginHandler{jwtSecret: cfg.jwtSecret})
	mux.Handle(http.MethodPut+" "+usersPath, &PutUsersHandler{jwtSecret: cfg.jwtSecret})
	mux.Handle(http.MethodPost+" "+refreshJWTPath, &RefreshJWTHandler{jwtSecret: cfg.jwtSecret})
	mux.HandleFunc(http.MethodPost+" "+revokeRefreshTokenPath, revokeRefreshToken)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
