package main

import (
	"encoding/json"
	"net/http"
	"log"
)

type request struct {
	Body string `json:"body"`
}

type response struct {
	Valid bool `json:"valid"`
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
		w.WriteHeader(400)
		data, err := json.Marshal(errorResponse{
			Error: "Something went wrong",
		})
		if err != nil {
			return
		}
		w.Write(data)
		return
	}
	is_valid := len(request.Body) <= 140
	response, err := json.Marshal(response{
		Valid: is_valid,
	})
	if err != nil {
		log.Printf("Error marshalling response JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	status := 200
	if !is_valid {
		status = 400
	}
	w.WriteHeader(status)
	w.Write(response)
}
