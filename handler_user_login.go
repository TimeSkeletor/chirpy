package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/timeskeletor/chirpy/internal/auth"
)

type UserWithoutPassword struct {
	ID        			uuid.UUID `json:"id"`
	CreatedAt 			time.Time `json:"created_at"`
	UpdatedAt 			time.Time `json:"updated_at"`
	Email     			string    `json:"email"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		UserWithoutPassword
	}

	decorder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decorder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decodere parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Incorrect email or password", err)
		return
	}

	if err := auth.CheckPassword(user.HashedPassword, params.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		UserWithoutPassword: UserWithoutPassword{
			ID: 		user.ID,
			CreatedAt: 	user.CreatedAt,
			UpdatedAt: 	user.UpdatedAt,
			Email: 		user.Email,
		},
	})
}