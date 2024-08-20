package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/Abhiram0106/chirpy/internal/database"
)

func revokeRefreshToken(w http.ResponseWriter, r *http.Request) {

	authHeader := strings.Fields(r.Header.Get("Authorization"))

	if len(authHeader) != 2 {
		respondWithError(w, http.StatusUnauthorized, "Invalid Refresh Token")
		return
	}

	refreshToken := authHeader[1]

	newDBConnection, dbConnectionErr := database.NewDB(databasePath)

	if dbConnectionErr != nil {
		log.Println(dbConnectionErr)
		respondWithError(w, http.StatusInternalServerError, dbConnectionErr.Error())
		return
	}

	revokeErr := newDBConnection.RevokeRefreshToken(refreshToken)

	if revokeErr != nil {
		log.Println(revokeErr)
		respondWithError(w, http.StatusInternalServerError, revokeErr.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
