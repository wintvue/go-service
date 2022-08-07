package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJSON(write http.ResponseWriter, req *http.Request, data any) error {
	maxBytes := 1048567
	req.Body = http.MaxBytesReader(write, req.Body, int64(maxBytes))

	dec := json.NewDecoder(req.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})

	if err != io.EOF {
		return err
	}

	return nil
}

func (app *Config) writeJSON(write http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)

	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			write.Header()[k] = v
		}
	}
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(status)

	_, err = write.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) errorJSON(write http.ResponseWriter, err error, status ...int) error {
	statCode := http.StatusBadRequest

	if len(status) > 0 {
		statCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(write, statCode, payload)
}
