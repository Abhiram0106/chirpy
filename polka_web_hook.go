package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/Abhiram0106/chirpy/internal/database"
)

type PostPolkaWebhook struct {
	polkaApiKey string
}

func (h *PostPolkaWebhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	headerAuth := strings.Fields(r.Header.Get("Authorization"))

	if len(headerAuth) != 2 {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	if headerAuth[1] != h.polkaApiKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	type webhookRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	defer r.Body.Close()
	body, readErr := io.ReadAll(r.Body)

	if readErr != nil {
		respondWithError(w, http.StatusInternalServerError, readErr.Error())
		return
	}

	whReq := webhookRequest{}

	unmarhsalErr := json.Unmarshal(body, &whReq)

	if unmarhsalErr != nil {
		respondWithError(w, http.StatusInternalServerError, unmarhsalErr.Error())
		return
	}

	if whReq.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, struct{}{})
		return
	}

	newDBConnection, dbConnectionErr := database.NewDB(databasePath)

	if dbConnectionErr != nil {
		respondWithError(w, http.StatusInternalServerError, dbConnectionErr.Error())
		return
	}

	dbUpgradeUserErr := newDBConnection.UpgradeUserToChirpyRed(whReq.Data.UserID)

	if dbUpgradeUserErr != nil {
		respondWithError(w, http.StatusNotFound, dbUpgradeUserErr.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
