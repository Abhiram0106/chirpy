package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Abhiram0106/chirpy/internal/database"
	"github.com/golang-jwt/jwt/v5"
)

type PutUsersHandler struct {
	jwtSecret string
}

type putUserReq struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (h *PutUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	body, readErr := io.ReadAll(r.Body)

	if readErr != nil {
		log.Println(readErr)
		respondWithError(w, http.StatusInternalServerError, readErr.Error())
		return
	}

	putUserReq := putUserReq{}

	marshalErr := json.Unmarshal(body, &putUserReq)

	if marshalErr != nil {
		log.Println(marshalErr)
		respondWithError(w, http.StatusInternalServerError, marshalErr.Error())
		return
	}

	authHeader := strings.Fields(r.Header.Get("Authorization"))

	if len(authHeader) != 2 {
		respondWithError(w, http.StatusUnauthorized, "Invalid JWT")
		return
	}

	jwtString := authHeader[1]
	claims := jwt.MapClaims{}

	jwToken, parseClaimErr := jwt.ParseWithClaims(
		jwtString,
		claims,
		func(t *jwt.Token) (interface{}, error) { return []byte(h.jwtSecret), nil },
	)

	if parseClaimErr != nil {
		log.Println(parseClaimErr)
		respondWithError(w, http.StatusUnauthorized, parseClaimErr.Error())
		return
	}

	userIDString, getSubjectErr := jwToken.Claims.GetSubject()

	if getSubjectErr != nil {
		log.Println(getSubjectErr)
		respondWithError(w, http.StatusUnauthorized, getSubjectErr.Error())
		return
	}

	userID, AtoIErr := strconv.Atoi(userIDString)

	if AtoIErr != nil {
		log.Println(AtoIErr)
		respondWithError(w, http.StatusInternalServerError, AtoIErr.Error())
		return
	}

	newDBConnection, dbConnectionErr := database.NewDB(databasePath)

	if dbConnectionErr != nil {
		log.Println(dbConnectionErr)
		respondWithError(w, http.StatusInternalServerError, dbConnectionErr.Error())
		return
	}

	updatedUser, updateDBErr := newDBConnection.UpdateUser(userID, putUserReq.Email, putUserReq.Password)

	if updateDBErr != nil {
		log.Println(updateDBErr)
		respondWithError(w, http.StatusInternalServerError, updateDBErr.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, updatedUser)
}
