package main

import (
	"encoding/json"
	"net/http"

	"github.com/ric-ram/go-chirpy/internal/auth"
)

type AuthenticatedUser struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	existingUser, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = auth.ValidatePassword(existingUser.Password, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect password")
		return
	}

	defaultExpiration := 60 * 60 * 24
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = defaultExpiration
	} else if params.ExpiresInSeconds > defaultExpiration {
		params.ExpiresInSeconds = defaultExpiration
	}

	accessJwtToken, err := auth.CreateJwtToken(existingUser.ID, cfg.jwtSecret, "access")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create Access JWT")
		return
	}

	refreshJwtToken, err := auth.CreateJwtToken(existingUser.ID, cfg.jwtSecret, "refresh")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create Refresh JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, AuthenticatedUser{
		ID:           existingUser.ID,
		Email:        existingUser.Email,
		Token:        accessJwtToken,
		RefreshToken: refreshJwtToken,
	})
}
