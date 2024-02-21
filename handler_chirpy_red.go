package main

import (
	"encoding/json"
	"net/http"

	"github.com/ric-ram/go-chirpy/internal/auth"
)

type UpgradedUser struct {
	Email       string `json:"email"`
	ID          int    `json:"id"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerChirpyRed(w http.ResponseWriter, r *http.Request) {
	// Get Api Key from header
	apiKey, err := auth.GetPolkaApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find Polka Api Key")
		return
	}

	err = auth.ValidatePolkaApiKey(apiKey, cfg.polkaApiSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Polka Api Key")
		return
	}

	// Define body struct
	type data struct {
		UserID int `json:"user_id"`
	}
	type parameters struct {
		Data  data   `json:"data"`
		Event string `json:"event"`
	}

	// Get params from body
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// Compare event to user.upgrade
	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, struct{}{})
		return
	}

	// Search user_id in database
	user, err := cfg.DB.GetUserByID(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User was not found")
		return
	}

	// Update user with user_id
	_, err = cfg.DB.UpgradeUser(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error upgrading user")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
