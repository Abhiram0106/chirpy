package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/Abhiram0106/chirpy/internal/database"
)

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func postUsers(w http.ResponseWriter, r *http.Request) {

	userReq := userRequest{}

	defer r.Body.Close()
	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		log.Println(readErr)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	unmarshalErr := json.Unmarshal(body, &userReq)
	if unmarshalErr != nil {
		log.Println(unmarshalErr)
		if errors.Is(unmarshalErr, &json.SyntaxError{}) {
			respondWithError(w, http.StatusBadRequest, unmarshalErr.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		}
		return
	}

	if validationErr := userReq.validate(); validationErr != nil {
		respondWithError(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	newDBConnection, dbConnectionErr := database.NewDB(databasePath)

	if dbConnectionErr != nil {
		log.Println(dbConnectionErr)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	newUser, createUserErr := newDBConnection.CreateUser(userReq.Email, userReq.Password)

	if createUserErr != nil {
		log.Println(createUserErr)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusCreated, newUser)
}

func (u *userRequest) validate() error {

	if u.Email == "" || u.Password == "" {
		return errors.New("Both email and password are required")
	}

	return nil
}
