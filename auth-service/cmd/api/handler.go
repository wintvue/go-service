package main

import (
	"bytes"
	"encoding/json"
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

	err = app.logRequest("authentication", fmt.Sprintf("%s logged in with email", user.Email))

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("User %s successfully logged in", user.Email),
		Data:    user,
	}
	_ = app.writeJSON(write, http.StatusOK, payload)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	req, err := http.NewRequest("POST", "http://logger-service/v1/writeLog", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(req)

	if err != nil {
		return err
	}

	return nil
}
