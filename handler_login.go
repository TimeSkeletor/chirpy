package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/timeskeletor/chirpy/internal/auth"
)

type UserLogin struct {
	ID        			uuid.UUID `json:"id"`
	CreatedAt 			time.Time `json:"created_at"`
	UpdatedAt 			time.Time `json:"updated_at"`
	Email     			string    `json:"email"`
	Token     			string    `json:"token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 				string	`json:"email"`
		Password 			string 	`json:"password"`
		Expires_in_seconds 	int		`json:"expires_in_seconds,omitempty"`
	}

	type response struct {
		UserLogin
	}

	decorder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decorder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decodere parameters", err)
		return
	}

	expiresInSeconds := 3600 // Default: 1 hour in seconds
	if params.Expires_in_seconds > 0 {
		// Client specified a value
		if params.Expires_in_seconds > 3600 {
			// If greater than 1 hour, cap at 1 hour
			expiresInSeconds = 3600
		} else {
			// Otherwise use the specified value
			expiresInSeconds = params.Expires_in_seconds
		}
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Incorrect email or password", err)
		return
	}

	if err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(expiresInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error generating token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		UserLogin: UserLogin{
			ID: 		user.ID,
			CreatedAt: 	user.CreatedAt,
			UpdatedAt: 	user.UpdatedAt,
			Email: 		user.Email,
			Token: 		token,
		},
	})
}
