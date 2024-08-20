package main

import (
	"crypto/rand"
	"encoding/hex"
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
	Email    string `json:"email"`
	Password string `json:"password"`
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

	jwtoken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
			Subject:   strconv.Itoa(user.ID),
		})

	signedJWTString, signingErr := jwtoken.SignedString([]byte(h.jwtSecret))

	if signingErr != nil {
		log.Println(signingErr)
		respondWithError(w, http.StatusInternalServerError, signingErr.Error())
		return
	}

	user.JWT = signedJWTString
	random32bytes := make([]byte, 32)
	_, randErr := rand.Read(random32bytes)

	if randErr != nil {
		log.Println(randErr)
		respondWithError(w, http.StatusInternalServerError, randErr.Error())
		return
	}

	refreshToken := hex.EncodeToString(random32bytes)

	addTokenToDBErr := dbConnection.AddRefreshToken(refreshToken, time.Now().AddDate(0, 2, 0), user.ID)

	if addTokenToDBErr != nil {
		log.Println(addTokenToDBErr)
		respondWithError(w, http.StatusInternalServerError, addTokenToDBErr.Error())
		return
	}

	user.RefreshToken = refreshToken
	respondWithJSON(w, http.StatusOK, user)
}
