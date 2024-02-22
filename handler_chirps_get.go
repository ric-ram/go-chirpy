package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/ric-ram/go-chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	defaultSortingOrder := "asc"

	// Get the chirp author_id from the query parameters if exists
	authorIDString := r.URL.Query().Get("author_id")

	// Get sorting order from query paramenters if exists
	sortingOrder := r.URL.Query().Get("sort")
	if sortingOrder == "" {
		sortingOrder = defaultSortingOrder
	}

	// Retrieve the chirps written by the author_id
	chirps := []Chirp{}
	dbChirps, err := cfg.manageGetChirps(authorIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:       chirp.ID,
			Body:     chirp.Body,
			AuthorID: chirp.AuthorID,
		})
	}

	sortedChirps := sortChirpsById(chirps, sortingOrder)

	respondWithJSON(w, http.StatusOK, sortedChirps)

}

func (cfg *apiConfig) handlerChirpsGetById(w http.ResponseWriter, r *http.Request) {
	paramID := chi.URLParam(r, "chirpID")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.DB.GetChirpsById(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:       chirp.ID,
		Body:     chirp.Body,
		AuthorID: chirp.AuthorID,
	})
}

func sortChirpsById(chirps []Chirp, order string) []Chirp {
	sort.Slice(chirps, func(i, j int) bool {
		if order == "asc" {
			return chirps[i].ID < chirps[j].ID
		} else if order == "desc" {
			return chirps[i].ID > chirps[j].ID
		} else {
			return false
		}
	})

	return chirps
}

// manageGetChirps returns list of chirps wheter there is or not an author_id
func (cfg *apiConfig) manageGetChirps(authorID string) ([]database.Chirp, error) {
	if authorID != "" {
		authorId, err := strconv.Atoi(authorID)
		if err != nil {
			return []database.Chirp{}, fmt.Errorf("error parsing author id")
		}

		dbChirps, err := cfg.DB.GetChirpsByAuthorId(authorId)
		if err != nil {
			return []database.Chirp{}, fmt.Errorf("couldn't retrieve chirps for the selected author id")
		}
		return dbChirps, nil

	} else {
		dbChirps, err := cfg.DB.GetChirps()
		if err != nil {
			return []database.Chirp{}, fmt.Errorf("couldn't retrieve chirps")
		}
		return dbChirps, nil
	}
}
