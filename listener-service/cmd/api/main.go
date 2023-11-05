package main

import (
	"fmt"
	"listener/event"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
}

func main() {
	fmt.Println("Starting listener service...")

	//connect to rabbitmq
	connection, err := connect()
	if err != nil {
		log.Println("Failed to connect to rabbitmq")
		os.Exit(1)
	}
	defer connection.Close()

	//start the consumer
	consumer, err := event.NewConsumer(connection)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	//watch for events
	err = consumer.Listen([]string{"log.INFO", "log.ERROR", "log.WARNING"})
	if err != nil {
		log.Println(err)
	}

	log.Println("Connected to rabbitmq", connection)
}

func connect() (*amqp.Connection, error) {
	//connect to rabbitmq
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	//do not continue until rabbit is ready

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("Failed to connect to rabbitmq")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			return nil, err
		}

		log.Printf("back off:%v, counts:%v", backoff, counts)
		time.Sleep(backoff)
		backoff *= 2
	}

	return connection, nil
}
