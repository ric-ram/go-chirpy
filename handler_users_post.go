package main

import (
	"encoding/json"
	"net/http"

	"github.com/ric-ram/go-chirpy/internal/auth"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerUsersPost(w http.ResponseWriter, r *http.Request) {
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

	encryptedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid password")
		return
	}

	user, err := cfg.DB.CreateUSer(params.Email, string(encryptedPassword))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:    user.ID,
		Email: params.Email,
	})
}
