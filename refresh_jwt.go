package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Abhiram0106/chirpy/internal/database"
	"github.com/golang-jwt/jwt/v5"
)

type RefreshJWTHandler struct {
	jwtSecret string
}

func (h *RefreshJWTHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

	userID, tokenValidityErr := newDBConnection.IsRefreshTokenValid(refreshToken)

	if tokenValidityErr != nil {
		log.Println(tokenValidityErr)
		respondWithError(w, http.StatusUnauthorized, tokenValidityErr.Error())
		return
	}

	jwtoken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
			Subject:   strconv.Itoa(userID),
		})

	signedJWTString, signingErr := jwtoken.SignedString([]byte(h.jwtSecret))

	if signingErr != nil {
		log.Println(signingErr)
		respondWithError(w, http.StatusInternalServerError, signingErr.Error())
		return
	}

	response := struct {
		NewJWT string `json:"token"`
	}{
		NewJWT: signedJWTString,
	}

	respondWithJSON(w, http.StatusOK, response)
}
