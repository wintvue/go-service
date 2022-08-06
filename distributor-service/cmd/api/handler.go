package main

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) distribute(write http.ResponseWriter, req *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "hit distributor service",
	}

	out, _ := json.MarshalIndent(payload, "", "\t")
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)
	write.Write(out)
}
