package main

import (
	"github.com/go-chi/chi/v5"
)

func (app *Config) routers(router chi.Router) {
	router.Post("/", app.distribute)
	router.Post("/handler", app.Handle)
}
