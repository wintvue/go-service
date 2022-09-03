package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	connection *amqp.Connection
	queueName  string
}

func NewConsumer(connection *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		connection: connection,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.connection.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type jsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQuee(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload jsonPayload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("waiting for message on %s", q.Name)
	<-forever

	return nil
}

func handlePayload(payload jsonPayload) error {
	switch payload.Name {
	case "log", "event":
		err := callLog(payload)
		if err != nil {
			log.Println(err)
		}
	// case "auth":

	default:
		err := callLog(payload)
		if err != nil {
			log.Println(err)
		}

		return nil
	}
	return nil
}

func callLog(data jsonPayload) error {
	jsonData, _ := json.MarshalIndent(data, "", "\t")

	request, err := http.NewRequest("POST", "http://logger-service/v1/writeLog", bytes.NewBuffer(jsonData))

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}

func declareExchange(channel *amqp.Channel) error {
	return channel.ExchangeDeclare(
		"logs_topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}

func declareRandomQuee(channel *amqp.Channel) (amqp.Queue, error) {
	return channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
}
