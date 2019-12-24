package main

import (
	"context"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createMongoClient() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(os.Getenv("ConnectionStrings__Mongo"))

	return mongo.Connect(context.TODO(), clientOptions)
}

func getMongoDatabase() (*mongo.Database, error) {
	client, err := createMongoClient()
	if err != nil {
		return nil, err
	}

	return client.Database("Settings"), nil
}

func getDefaultSettings(database *mongo.Database, path []string) (string, *bson.Raw, error) {
	defaultsCollection := database.Collection("DefaultSettings")

	filter := bson.D{{"path", path}}

	result := defaultsCollection.FindOne(context.TODO(), filter)

	doc, err := result.DecodeBytes()

	if err != nil {
		return "", nil, err
	}

	docidValue := doc.Lookup("_id")

	docid := docidValue.StringValue()

	dataValue := doc.Lookup("data")

	data := dataValue.Document()

	return docid, &data, nil
}

func getUserSettings(database *mongo.Database, settingID string, userID string) (*bson.Raw, error) {
	userCollection := database.Collection("UserSettings")

	filter := bson.D{{"_id", bson.D{{"settingId", settingID}, {"userId", userID}}}}

	result := userCollection.FindOne(context.TODO(), filter)

	if result.Err() == nil {
		doc, err := result.DecodeBytes()

		if err != nil {
			return nil, err
		}

		dataValue := doc.Lookup("data")

		data := dataValue.Document()

		return &data, nil
	}

	return nil, nil
}

func mergeSettings(docs ...*bson.Raw) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	first := docs[0]

	elements, err := first.Elements()

	if err != nil {
		return nil, err
	}

	for _, element := range elements {
		settingElement := element.Value()
		for i := 1; i < len(docs); i++ {
			if docs[i] != nil {
				nextDocValue := docs[i].Lookup(element.Key())

				settingElement = nextDocValue
			}

			switch element.Value().Type {
			case bsontype.Boolean:
				m[element.Key()] = settingElement.Boolean()
			case bsontype.DateTime:
				m[element.Key()] = settingElement.DateTime()
			case bsontype.Decimal128:
				m[element.Key()] = settingElement.Decimal128()
			case bsontype.Double:
				m[element.Key()] = settingElement.Double()
			case bsontype.Int32:
				m[element.Key()] = settingElement.Int32()
			case bsontype.Int64:
				m[element.Key()] = settingElement.Int64()
			case bsontype.ObjectID:
				m[element.Key()] = settingElement.ObjectID()
			case bsontype.String:
				m[element.Key()] = settingElement.String()
			default:
				return nil, errors.New("Unsupported settings type")
			}

		}
	}

	return m, nil
}
