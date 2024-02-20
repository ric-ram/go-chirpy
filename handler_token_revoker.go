package main

import (
	"net/http"

	"github.com/ric-ram/go-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerTokenRevoke(w http.ResponseWriter, r *http.Request) {
	// get header token
	headerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	// check if is valid refresh token
	_, err = auth.ValidateRefreshJwtToken(headerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't validate JWT")
		return
	}

	// Revoke the token in the database
	// refreshToken : date of revoken
	// use refreshToken as id
	err = cfg.DB.AddRevokeToken(headerToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke the token")
		return
	}

	respondWithJSON(w, http.StatusOK, "Token succesfully revoked")
}
