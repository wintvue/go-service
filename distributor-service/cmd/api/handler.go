package main

import (
	"bytes"
	"distributor/event"
	"encoding/json"
	"errors"
	"net/http"
)

type requestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
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
	} else if payload.Action == "log" {
		app.logEventWithMQ(write, payload.Log)
	} else if payload.Action == "mail" {
		app.sendMail(write, payload.Mail)
	} else {
		app.errorJSON(write, errors.New("unknown action "+payload.Action))
	}
}

func (app *Config) logItem(write http.ResponseWriter, log LogPayload) {
	jsonData, _ := json.MarshalIndent(log, "", "\t")

	request, err := http.NewRequest("POST", "http://logger-service/v1/writeLog", bytes.NewBuffer(jsonData))

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(write, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(write, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged!"

	app.writeJSON(write, http.StatusAccepted, payload)
}

func (app *Config) Authenticate(write http.ResponseWriter, auth AuthPayload) {
	jsonData, _ := json.MarshalIndent(auth, "", "\t")
	request, err := http.NewRequest("POST", "http://auth-service/v1/authenticate", bytes.NewBuffer(jsonData))
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

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(write, errors.New("Invalid Credentials"))
		return
	} else if response.StatusCode == http.StatusAccepted {
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

func (app *Config) sendMail(write http.ResponseWriter, mail MailPayload) {
	jsonData, _ := json.MarshalIndent(mail, "", "\t")

	request, err := http.NewRequest("POST", "http://mail-service/v1/sendMail", bytes.NewBuffer(jsonData))

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(write, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(write, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Email sent successfully!!"

	app.writeJSON(write, http.StatusAccepted, payload)
}

func (app *Config) logEventWithMQ(w http.ResponseWriter, log LogPayload) {
	err := app.pushToQueue(log.Name, log.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via queue"

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.conn)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")

	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}
