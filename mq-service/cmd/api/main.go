package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
}

func main() {
	connect, err := connectMQ()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer connect.Close()

	log.Println("Listening for MQ messages...")

	consumer, err := NewConsumer(connect)

	if err != nil {
		panic(err)
	}

	err = consumer.Listen([]string{"log.INFO", "log.WARN", "log.ERRROR"})

	if err != nil {
		log.Println(err)
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
