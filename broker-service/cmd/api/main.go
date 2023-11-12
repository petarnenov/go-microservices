package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "8080"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	connection, err := connect()
	if err != nil {
		log.Println("Failed to connect to rabbitmq")
		os.Exit(1)
	}
	defer connection.Close()

	app := Config{
		Rabbit: connection,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
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
