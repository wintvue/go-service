package main

import (
	"log"
	"net/http"
)

type Config struct {
}

const port = "80"

func main() {
	app := Config{}

	log.Println("Starting app...")
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
