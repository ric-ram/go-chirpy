package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/ric-ram/go-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpDelete(w http.ResponseWriter, r *http.Request) {
	// Get header token
	headerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	// Validate if is access token
	validToken, err := auth.ValidateAccessJwtToken(headerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	// Get current user ID
	currentUserIDString, err := auth.GetUserID(validToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting id from token")
	}

	currentUserID, err := strconv.Atoi(currentUserIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing id")
		return
	}

	// Get chirp ID from request params
	paramID := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(paramID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	// Get chirp by ID
	chirp, err := cfg.DB.GetChirpsById(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	// Compare chirp author ID with Token current user ID
	if chirpID != currentUserID {
		respondWithError(w, http.StatusForbidden, "Incorrect user")
		return
	}

	// Delete Chirp from database
	err = cfg.DB.DeleteChirp(chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, "")
}
