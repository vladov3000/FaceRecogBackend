package database

import "go.mongodb.org/mongo-driver/bson"

type Database interface {
	Disconnect()
	AddPerson(person Person) error
	GetPerson(filter bson.M) (Person, bool, error)
}
