package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func UserCollection(client *mongo.Client) *mongo.Collection {
	return client.Database("user").Collection("users")
}
