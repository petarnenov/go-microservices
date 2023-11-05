package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	chanel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(chanel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		err = ch.QueueBind(
			q.Name,
			topic,
			"logs_topic",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			payload := Payload{}
			err := json.Unmarshal(d.Body, &payload)
			if err != nil {
				log.Println("Failed to unmarshal message")
				continue
			}

			go handlePayload(payload)
		}
	}()

	log.Printf(" [*] Waiting for message from [Exchange, Queue] [logs_topic,%s]", q.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	log.Printf("Received message: %v", payload)
	switch payload.Name {
	case "log", "event":
		log.Println("Log or event message")
		err := handleLog(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		log.Println("Auth message")
	default:
		log.Println("Unknown event")
	}
}

func handleLog(payload Payload) error {
	// create some json we'll send to the logger microservice
	jsonData, _ := json.MarshalIndent(payload, "", "\t")

	logServiceURL := "http://logger-service/log"

	// call the service
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {

		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {

		return err
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		return errors.New("error calling log service")
	}

	return nil
}
