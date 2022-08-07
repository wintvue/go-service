package main

import (
	"net/http"
)

func (app *Config) distribute(write http.ResponseWriter, req *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "hit distributor service",
	}
	_ = app.writeJSON(write, http.StatusOK, payload)
}
