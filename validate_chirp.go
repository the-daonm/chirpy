package main

import (
	"net/http"
	"encoding/json"
	"log"
	"strings"
)

func cleanProfanity(body string) string {
	profanes := map[string]bool{
			"kerfuffle": true,
			"sharbert":  true,
			"fornax":    true,
	}
	words := strings.Split(body, " ")
	for i, word := range words {
		if profanes[strings.ToLower(word)] {
				words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	errResponse := errorResponse{
		Error: msg,
	}

	respondWithJSON(w, code, errResponse)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}		

	dec := json.NewDecoder(r.Body)
	params := parameters{}
	err := dec.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	type successResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}
	sucResponse := successResponse{
		CleanedBody: cleanProfanity(params.Body),
	}

	respondWithJSON(w, 200, sucResponse)
}
