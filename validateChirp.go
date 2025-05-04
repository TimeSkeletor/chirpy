package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

const (
	chirpSizeLimit = 140
	cleanWord = "****"
)
var badWords = []string{"kerfuffle", "sharbert", "fornax"}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	type chirp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	c := chirp{}
	err := decoder.Decode(&c)
	if err != nil {
		log.Printf("Something went wrong: %s", err)
		w.WriteHeader(500)
		return 
	}

	valid := len(c.Body) <= chirpSizeLimit

	if valid {
		respondWithJSON(w, http.StatusOK, cleanChirp(c.Body))
		return
	} else {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	} 

}

func cleanChirp(msg string) string {
	words := strings.Fields(msg)
	for i, _ := range words {
		word := strings.ToLower(words[i])
		for j, _ := range badWords{
			if word == badWords[j] {
				words[i] = cleanWord
			}
		}
	}
	return strings.Join(words, " ")
}
