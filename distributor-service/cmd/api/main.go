package main

import (
	"log"
	"net/http"
)

const port = "8000"

type Config struct{}

func main() {
	app := Config{}

	log.Printf("Starting distributor service on port %s\n", port)

	// http server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: app.routes(),
	}

	// start http server
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
