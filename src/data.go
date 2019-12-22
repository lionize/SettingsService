package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func create_mongo_client() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://root:9fP30ErG0fBv5R@localhost:52540")

	return mongo.Connect(context.TODO(), clientOptions)
}

func get_mongo_database() (*mongo.Database, error) {
	client, err := create_mongo_client()
	if err != nil {
		return nil, err
	}

	return client.Database("Settings"), nil
}
