package main

import (
	"net/http"
	"strconv"

	"github.com/ric-ram/go-chirpy/internal/auth"
	"github.com/ric-ram/go-chirpy/internal/database"
)

func (cfg *apiConfig) handlerTokenRefresh(w http.ResponseWriter, r *http.Request) {
	// define response type
	type response struct {
		Token string `json:"token"`
	}

	// get header token - if not 401
	headerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	// check if is valid refresh token - if not 401
	validToken, err := auth.ValidateRefreshJwtToken(headerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}
	if validToken == nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh Token")
		return
	}

	// check if token is not revoked - if not 401
	_, err = cfg.DB.GetRevokedTokenById(headerToken)
	if err != database.ErrNotExist {
		respondWithError(w, http.StatusUnauthorized, "Couldn't read database")
		return
	}

	// create new access tokem - 200
	userIDString, err := auth.GetUserID(validToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting id from token")
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing id")
		return
	}

	newToken, err := auth.CreateJwtToken(userID, cfg.jwtSecret, "access")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create Access JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: newToken,
	})
}
