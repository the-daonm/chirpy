package api

import (
	"encoding/json"
	"log"
	"net/http"
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
		loweredWord := strings.ToLower(word)
		trimmedWord := strings.Trim(loweredWord, ".,!?;:")
		if profanes[trimmedWord] {
			prefix := ""
			for j := 0; j < len(word); j++ {
				if !strings.ContainsAny(string(word[j]), ".,!?;:") {
					prefix = word[:j]
					break
				}
			}

			suffix := ""
			for j := len(word) - 1; j >= 0; j-- {
				if !strings.ContainsAny(string(word[j]), ".,!?;:") {
					suffix = word[j+1:]
					break
				}
			}

			words[i] = prefix + "****" + suffix
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
