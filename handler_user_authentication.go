package main

import (
	"encoding/json"
	"net/http"

	"github.com/ric-ram/go-chirpy/internal/auth"
)

type AuthenticatedUser struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
		IsChirpyRed:  existingUser.IsChirpRed,
		Token:        accessJwtToken,
		RefreshToken: refreshJwtToken,
	})
}
