package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Abhiram0106/chirpy/internal/database"
	"github.com/golang-jwt/jwt/v5"
)

type PostLoginHandler struct {
	jwtSecret string
}

type loginRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

func (h *PostLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	body, readErr := io.ReadAll(r.Body)

	if readErr != nil {
		log.Println(readErr)
		respondWithError(w, http.StatusInternalServerError, readErr.Error())
		return
	}

	loginReq := loginRequest{}

	marshalErr := json.Unmarshal(body, &loginReq)

	if marshalErr != nil {
		log.Println(marshalErr)
		respondWithError(w, http.StatusInternalServerError, marshalErr.Error())
		return
	}

	dbConnection, dbConnectionErr := database.NewDB(databasePath)

	if dbConnectionErr != nil {
		log.Println(dbConnection)
		respondWithError(w, http.StatusInternalServerError, dbConnectionErr.Error())
		return
	}

	user, getUsersErr := dbConnection.GetUserByEmailAndPassword(loginReq.Email, loginReq.Password)

	if getUsersErr != nil {
		log.Println(getUsersErr)
		respondWithError(w, http.StatusUnauthorized, getUsersErr.Error())
		return
	}

	if loginReq.ExpiresInSeconds == 0 {
		loginReq.ExpiresInSeconds = 24 * 60 * 60
	}

	jwtoken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(loginReq.ExpiresInSeconds) * time.Second)),
			Subject:   strconv.Itoa(user.ID),
		})

	signedJWTString, signingErr := jwtoken.SignedString([]byte(h.jwtSecret))

	if signingErr != nil {
		log.Println(signingErr)
		respondWithError(w, http.StatusInternalServerError, signingErr.Error())
		return
	}

	user.JWT = signedJWTString

	respondWithJSON(w, http.StatusOK, user)
}
