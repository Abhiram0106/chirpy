package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Abhiram0106/chirpy/internal/database"
	"github.com/golang-jwt/jwt/v5"
)

type DeleteChirpHandler struct {
	jwtSecret string
}

func (h *DeleteChirpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	authHeader := strings.Fields(r.Header.Get("Authorization"))

	if len(authHeader) != 2 {
		respondWithError(w, http.StatusUnauthorized, "Invalid JWT")
		return
	}

	claims := jwt.MapClaims{}

	jwToken, parseClaimErr := jwt.ParseWithClaims(
		authHeader[1],
		claims,
		func(t *jwt.Token) (interface{}, error) { return []byte(h.jwtSecret), nil },
	)

	if parseClaimErr != nil {
		respondWithError(w, http.StatusUnauthorized, parseClaimErr.Error())
		return
	}

	authorIDString, getClaimSubjectErr := jwToken.Claims.GetSubject()

	if getClaimSubjectErr != nil {
		respondWithError(w, http.StatusUnauthorized, getClaimSubjectErr.Error())
		return
	}

	authorID, strconvAuthorIDErr := strconv.Atoi(authorIDString)

	if strconvAuthorIDErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid authorID")
		return
	}

	chirpIDString := r.PathValue(chirpIDWildCard)

	if chirpIDString == "" {
		respondWithError(w, http.StatusUnauthorized, "Specify chirp ID")
		return
	}

	chirpID, strconvErr := strconv.Atoi(chirpIDString)

	if strconvErr != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	newDBConnection, dbConnectionErr := database.NewDB(databasePath)

	if dbConnectionErr != nil {
		respondWithError(w, http.StatusInternalServerError, dbConnectionErr.Error())
		return
	}

	if deleteChirpErr := newDBConnection.DeleteChirpByID(chirpID, authorID); deleteChirpErr != nil {
		respondWithError(w, http.StatusUnauthorized, deleteChirpErr.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
