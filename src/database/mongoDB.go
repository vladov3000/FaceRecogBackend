package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	dbName     string
	collection *mongo.Collection
	client     *mongo.Client
	ctx        context.Context
	cancel     context.CancelFunc
}

func getCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func NewMongoDB(dbName string) MongoDB {
	ctx, cancel := getCtx()
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017/"))
	if err != nil {
		log.Fatalf("Error while connecting to database: %s", err)
	}

	for err = errors.New(""); err != nil; err = client.Ping(ctx, readpref.Primary()) {
		log.Print("Searching for Mongo database...")
		time.Sleep(time.Second)
	}
	log.Print("Found Mongo database")

	collection := client.Database(dbName).Collection("person")

	return MongoDB{
		dbName:     dbName,
		collection: collection,
		client:     client,
		cancel:     cancel,
	}
}

func (db MongoDB) Disconnect() {
	ctx, cancel := getCtx()
	defer cancel()

	if err := db.client.Disconnect(ctx); err != nil {
		log.Fatalf("Error while disconnecting: %s", err)
	}
}

func (db MongoDB) AddPerson(person Person) error {
	ctx, cancel := getCtx()
	defer cancel()

	data, err := bson.Marshal(person)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to marshal %+v into bson: %s", person, err))
	}

	db.collection.InsertOne(ctx, data)
	return nil
}

func (db MongoDB) GetPerson(filter bson.M) (Person, bool, error) {
	var res Person

	ctx, cancel := getCtx()
	defer cancel()

	cur := db.collection.FindOne(ctx, filter)
	if err := cur.Err(); err == mongo.ErrNoDocuments {
		return res, false, nil
	}

	if err := cur.Decode(&res); err != nil {
		return res, true, errors.New(fmt.Sprintf("Couldn't decode %+v: %s", res, err))
	}

	return res, true, nil
}
