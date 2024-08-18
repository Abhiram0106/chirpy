package main

import (
	"log"
	"net/http"

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

	respondWithJSON(w, http.StatusOK, chirps)
}
