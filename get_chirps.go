package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Abhiram0106/chirpy/internal/database"
)

func getChirps(w http.ResponseWriter, r *http.Request) {

	authorIDString := r.URL.Query().Get("author_id")

	DBConnection, dbConnectionErr := database.NewDB(databasePath)

	if dbConnectionErr != nil {
		log.Println(dbConnectionErr.Error())
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if chirpIDStr := r.PathValue(chirpIDWildCard); len(chirpIDStr) != 0 {
		requestedChirpID, AtoiErr := strconv.Atoi(chirpIDStr)

		if AtoiErr != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid chirp id")
			return
		}

		chirp, getChirpErr := DBConnection.GetChirpByID(requestedChirpID)

		if getChirpErr != nil {
			respondWithError(w, http.StatusNotFound, "No chirp found")
			return
		}

		respondWithJSON(w, http.StatusOK, chirp)
		return
	}

	authorID, _ := strconv.Atoi(authorIDString)

	chirps, getChirpsErr := DBConnection.GetChirps(authorID)

	if getChirpsErr != nil {
		log.Println(getChirpsErr.Error())
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
