package main

import (
	"encoding/json"
	"net/http"

	"github.com/markphelps/optional"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string        `json:"email"`
		Password         string        `json:"password"`
		ExpiresInSeconds *optional.Int `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	existingUser, exists, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve user from database")
		return
	}
	if !exists {
		respondWithError(w, http.StatusNotFound, "No user with that email")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect password")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:    existingUser.ID,
		Email: existingUser.Email,
	})
}
