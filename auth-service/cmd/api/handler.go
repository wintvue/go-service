package main

import (
	"fmt"
	"net/http"
)

// import (
// 	"database/sql"
// 	"log"
// 	"os"
// 	"time"
// )

func (app *Config) Authenticate(write http.ResponseWriter, req *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(write, req, &requestPayload)
	if err != nil {
		app.errorJSON(write, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		app.errorJSON(write, err, http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(write, err, http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("User %s successfully logged in", user.Email),
		Data:    user,
	}
	_ = app.writeJSON(write, http.StatusOK, payload)
}
