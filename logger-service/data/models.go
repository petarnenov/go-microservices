package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error inserting log entry: ", err)
		return err
	}
	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	ops := options.Find()
	ops.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(context.Background(), bson.D{}, ops)
	if err != nil {
		log.Println("error finding log entries: ", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var entries []*LogEntry
	for cursor.Next(ctx) {
		var entry LogEntry
		err := cursor.Decode(&entry)
		if err != nil {
			log.Println("error decoding log entry: ", err)
			return nil, err
		}
		entries = append(entries, &entry)
	}

	return entries, nil
}

func (l *LogEntry) GetEntryById(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("error converting id to object id: ", err)
		return nil, err
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docId}).Decode(&entry)
	if err != nil {
		log.Println("error finding log entry: ", err)
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) DeleteEntryById(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("error converting id to object id: ", err)
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": docId})
	if err != nil {
		log.Println("error deleting log entry: ", err)
		return err
	}

	return nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	err := collection.Drop(ctx)
	if err != nil {
		log.Println("error dropping collection: ", err)
		return err
	}

	return nil
}

func (l *LogEntry) UpdateEntryById(id string, entry LogEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("error converting id to object id: ", err)
		return err
	}

	_, err = collection.UpdateOne(ctx,
		bson.M{"_id": docId},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: entry.Name},
				{Key: "data", Value: entry.Data},
				{Key: "updated_at", Value: time.Now()},
			}},
		},
	)
	if err != nil {
		log.Println("error updating log entry: ", err)
		return err
	}

	return nil
}
