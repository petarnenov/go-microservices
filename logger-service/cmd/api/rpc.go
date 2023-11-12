package main

import (
	"context"
	"log"
	"logger/data"
	"time"
)

type RPCServer struct {
}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogItem(payload RPCPayload, result *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error inserting log entry: ", err)
		return err
	}

	*result = "Processed payload via RPC: " + payload.Name 

	return nil
}