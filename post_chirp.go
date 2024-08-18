package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Abhiram0106/chirpy/internal/database"
)

func postChirp(w http.ResponseWriter, r *http.Request) {

	type chirp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirpVal := chirp{}
	err := decoder.Decode(&chirpVal)

	if err != nil {
		log.Printf("PostChirp Error decoding %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(chirpVal.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	DBConnection, DBErr := database.NewDB(databasePath)

	if DBErr != nil {
		log.Println(DBErr.Error())
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	newChirp, DBCreateChirpErr := DBConnection.CreateChirp(chirpProfanityFilter(chirpVal.Body))

	if DBCreateChirpErr != nil {
		log.Println(DBCreateChirpErr.Error())
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusAccepted, newChirp)
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
