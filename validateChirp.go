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
		respondWithJSON(w, 200, cleanChirp(c.Body))
	} else {
		respondWithError(w, 400, "Chirp is too long")
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

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	respBody := errorResponse{
		Error: msg,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
    w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	type validChirp struct {
		CleanedBody string `json:"cleaned_body"`
	}
	respBody := validChirp{
		CleanedBody:  payload.(string),
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
    w.Write(dat)
}