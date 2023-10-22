package main

import (
	"authentication/data"
	"database/sql"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting Authentication service on port", webPort)

	//TODO: connect to database

	//set up config
	app := Config{}

	srv := &http.Server{
		Addr: ":" + webPort,
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if(err != nil) {
		log.Panic(err)
	}
}
