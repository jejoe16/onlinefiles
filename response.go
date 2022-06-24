package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type PostResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type PostUid struct {
	UID uuid.UUID `json:"id" validate:"required"`
}

type PostID struct {
	ID int `json:"id" validate:"gte=0,required,numeric"`
}

type PostEmail struct {
	Email string `json:"email" validate:"required,gte=4,lte=320,email,lowercase"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}
