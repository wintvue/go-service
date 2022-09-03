package main

import "github.com/go-chi/chi/v5"

func (app *Config) routers(router chi.Router) {
	router.Post("/sendMail", app.SendMail)
}
