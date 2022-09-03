package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const port = "80"

type Config struct {
	conn *amqp.Connection
}

func main() {
	connect, errs := connectMQ()
	if errs != nil {
		log.Println(errs)
		os.Exit(1)
	}
	defer connect.Close()

	app := Config{
		conn: connect,
	}

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

func connectMQ() (*amqp.Connection, error) {
	var counts int64
	var sleep = 1 * time.Second
	var connection *amqp.Connection
	var err error

	for {
		connection, err = amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("MQ not ready")
			counts += 1
		} else {
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		sleep = time.Duration(math.Pow(float64(counts), 2)) * time.Second

		log.Println("sleep")
		time.Sleep(sleep)
		continue
	}

	return connection, nil
}
