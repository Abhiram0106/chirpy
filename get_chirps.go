package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Abhiram0106/chirpy/internal/database"
)

func getChirps(w http.ResponseWriter, r *http.Request) {

	DBConnection, dbConnectionErr := database.NewDB(databasePath)

	if dbConnectionErr != nil {
		log.Println(dbConnectionErr.Error())
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	chirps, getChirpsErr := DBConnection.GetChirps()

	if getChirpsErr != nil {
		log.Println(getChirpsErr.Error())
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if chirpIDStr := r.PathValue(chirpIDWildCard); len(chirpIDStr) != 0 {
		chirpIDIndex, AtoiErr := strconv.Atoi(chirpIDStr)

		if AtoiErr != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid chirp id")
			return
		}

		chirpIDIndex--

		if chirpIDIndex < 0 || chirpIDIndex >= len(chirps) {
			respondWithError(w, http.StatusNotFound, "No chirp found")
			return
		}

		chirpOfID := chirps[chirpIDIndex]

		respondWithJSON(w, http.StatusFound, chirpOfID)
		return
	}

	respondWithJSON(w, http.StatusFound, chirps)
}
