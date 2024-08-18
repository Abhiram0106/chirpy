package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {

	type chirp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirpVal := chirp{}
	err := decoder.Decode(&chirpVal)

	if err != nil {
		log.Printf("ValidChirp Error decoding %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(chirpVal.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedChirp := struct {
		CleanedChirp string `json:"cleaned_body"`
	}{
		CleanedChirp: chirpProfanityFilter(chirpVal.Body),
	}

	respondWithJSON(w, http.StatusOK, cleanedChirp)
}

func chirpProfanityFilter(chirp string) (cleanedChirp string) {

	words := strings.Fields(chirp)

	profaneList := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	for i, word := range words {
		if _, ok := profaneList[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}

	cleanedChirp = strings.Join(words, " ")
	log.Println(cleanedChirp)

	return
}
