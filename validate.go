package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func validate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Valid string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Valid: filter(params.Body),
	})
}

func filter(s string) string {
	profane := map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}
	words := strings.Split(s, " ")
	filtered := []string{}
	for i := 0; i < len(words); i++ {
		if profane[strings.ToLower(words[i])] {
			filtered = append(filtered, "****")
		} else {
			filtered = append(filtered, words[i])
		}
	}
	return strings.Join(filtered, " ")
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}