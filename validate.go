package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func validate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type Error struct {
		Error string `json:"error"`
	}

	type valid struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	errResponse := Error{
		Error: "Something went wrong",
	}

	errResponse2 := Error{
		Error: "Chirp is too long",
	}
	
	validResponse := valid{
		Valid: true,
	}

	if err != nil {
		dat, _ := json.Marshal(errResponse)
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(dat)
		return
	}
	if len(params.Body) > 140 {
		dat, _ := json.Marshal(errResponse2)
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(dat)
		return
	}

	dat, _ := json.Marshal(validResponse)
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}
