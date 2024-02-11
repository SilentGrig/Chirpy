package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type request struct {
	Body string `json:"body"`
}

type response struct {
	CleanedBody string `json:"cleaned_body"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	request := request{}
	err := decoder.Decode(&request)
	if err != nil {
		log.Printf("Error decoding request: %s", err)
		responseWithError(w, http.StatusBadRequest, "Couldn't decode request")
		return
	}
	const maxChirpLength = 140
	if len(request.Body) > maxChirpLength {
		responseWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	responseWithJson(w, http.StatusOK, response{
		CleanedBody: cleanChirp(request.Body),
	})
}

func responseWithError(w http.ResponseWriter, status int, msg string) {
	if status >= 500 {
		log.Printf("Responding with %d error: %s", status, msg)
	}
	responseWithJson(w, status, errorResponse{
		Error: msg,
	})
}

func responseWithJson(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON response: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(status)
	w.Write(data)
}

func cleanChirp(s string) string {
	profaneWords := [3]string{"kerfuffle", "sharbert", "fornax"}
	splitWords := strings.Split(s, " ")
	for i, word := range splitWords {
		lowerCasedWord := strings.ToLower(word)
		for _, badWord := range profaneWords {
			if lowerCasedWord == badWord {
				splitWords[i] = "****"
			}
		}
	}
	return strings.Join(splitWords, " ")
}
