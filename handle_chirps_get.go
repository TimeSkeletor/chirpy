package main

import (
	"net/http"
	"sort"

	"github.com/timeskeletor/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	authorIDParam := r.URL.Query().Get("author_id")
	sortParam := r.URL.Query().Get("sort")
	
	var dbChirps []database.Chirp
	var err error

	if authorIDParam != "" {
		authorID, parseErr := uuid.Parse(authorIDParam)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id UUID", parseErr)
			return
		}

		dbChirps, err = cfg.db.GetChirpsFromAuthor(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}
	} else {
		dbChirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}

	if sortParam != "" {
		switch sortParam {
		case "asc":
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
			})
		case "desc":
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
			})
		default:
		}
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
