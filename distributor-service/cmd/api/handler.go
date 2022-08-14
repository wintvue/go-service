package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type requestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) distribute(write http.ResponseWriter, req *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "hit distributor service",
	}
	_ = app.writeJSON(write, http.StatusOK, payload)
}

func (app *Config) Handle(write http.ResponseWriter, req *http.Request) {
	var payload requestPayload
	err := app.readJSON(write, req, &payload)
	if err != nil {
		app.errorJSON(write, err)
		return
	}

	if payload.Action == "auth" {
		app.Authenticate(write, payload.Auth)

	} else {
		app.errorJSON(write, errors.New("unknown action "+payload.Action))
	}
}

func (app *Config) Authenticate(write http.ResponseWriter, auth AuthPayload) {
	jsonData, _ := json.MarshalIndent(auth, "", "\t")
	request, err := http.NewRequest("POST", "http://auth-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(write, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(write, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusUnauthorized {
		app.errorJSON(write, errors.New("Invalid Credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(write, errors.New("Error when calling auth-service"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(write, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(write, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(write, http.StatusAccepted, payload)
}
