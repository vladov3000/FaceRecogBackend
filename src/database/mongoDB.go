package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client mongo.Client
}

func NewMongoDB() MongoDB {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Error while connecting to database: %s", err)
	}

	for err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Print("Searching for Mongo database...")
		time.Sleep(time.Second) 
	}
	log.Print("Found Mongo database")

	return MongoDB{
		client: client,
	}
}

func (db MongoDB) Disconnect() {
	if err := client.Disconnect(ctx); err != nil {
		log.Fatalf("Error while disconnecting: %s", err)
	}
}
