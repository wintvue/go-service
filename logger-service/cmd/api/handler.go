package main

import (
	"logger/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(write http.ResponseWriter, req *http.Request) {
	var requestPayload JSONPayload

	errs := app.readJSON(write, req, &requestPayload)
	if errs != nil {
		return
	}

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(write, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(write, http.StatusAccepted, resp)
}
