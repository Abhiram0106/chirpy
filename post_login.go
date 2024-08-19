package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Abhiram0106/chirpy/internal/database"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func postLogin(w http.ResponseWriter, r *http.Request) {

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

	respondWithJSON(w, http.StatusOK, user)
}
