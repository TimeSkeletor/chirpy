package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const chirpSizeLimit = 140

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	type chirp struct {
		Body string `json:"body"`
	}

	type validChirp struct {
		Valid bool `json:"valid"`
	}

	type errorResponse struct {
		Error string `json:"error"`
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
	var status int
	var respBody interface{}

	if valid {
		status = 200

		respBody = validChirp{
			Valid: valid,
		}
	} else {
		status = 400
		respBody = errorResponse{
			Error: "Chirp is too long",
		}
	} 

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write(dat)
}
