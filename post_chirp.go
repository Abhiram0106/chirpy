package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Abhiram0106/chirpy/internal/database"
	"github.com/golang-jwt/jwt/v5"
)

type PostChirpHandler struct {
	jwtSecret string
}

type chirp struct {
	Body string `json:"body"`
}

func (h *PostChirpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

	authorID, AtoIErr := strconv.Atoi(userIDString)

	if AtoIErr != nil {
		log.Println(AtoIErr)
		respondWithError(w, http.StatusInternalServerError, AtoIErr.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	chirpVal := chirp{}
	err := decoder.Decode(&chirpVal)

	if err != nil {
		log.Printf("PostChirp Error decoding %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(chirpVal.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	if validationErr := chirpVal.validate(); validationErr != nil {
		respondWithError(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	DBConnection, DBErr := database.NewDB(databasePath)

	if DBErr != nil {
		log.Println(DBErr.Error())
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	newChirp, DBCreateChirpErr := DBConnection.CreateChirp(chirpProfanityFilter(chirpVal.Body), authorID)

	if DBCreateChirpErr != nil {
		log.Println(DBCreateChirpErr.Error())
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusCreated, newChirp)
}

func chirpProfanityFilter(chirp string) (cleanedChirp string) {

	words := strings.Fields(chirp)

	profaneList := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	for i, word := range words {
		if _, ok := profaneList[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}

	cleanedChirp = strings.Join(words, " ")
	log.Println(cleanedChirp)

	return
}

func (c *chirp) validate() error {
	if c.Body == "" {
		return errors.New("body can't be empty")
	}
	return nil
}
