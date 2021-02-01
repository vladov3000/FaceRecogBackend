package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

interface Database {
	Disconnect()
}
