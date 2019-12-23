package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/kataras/iris"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// You will be using this Trainer type later in the program
type Trainer struct {
	Name string
	Age  int
	City string
}

func main() {
	settingName := "MainPagePagination"
	fmt.Println("Settings name is ", settingName)

	path := []string{settingName}

	database, err := getMongoDatabase()
	if err != nil {
		log.Fatal(err)
	}

	docid, defaultData, err := getDefaultSettings(database, path)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(docid)
	fmt.Println(defaultData)

	user1id := "7b803e2d-ee0e-4213-a025-9db732bcbb2e"
	user2id := "ad2ea197-310a-4832-940c-2935bd6fa511"

	user1Data, err := getUserSettings(database, docid, user1id)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(user1Data)

	user2Data, err := getUserSettings(database, docid, user2id)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(user2Data)

	m1, err := mergeSettings(defaultData, user1Data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m1)

	m2, err := mergeSettings(defaultData, user2Data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m2)

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://root:9fP30ErG0fBv5R@localhost:52540")

	var client *mongo.Client
	// Connect to MongoDB
	client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("test").Collection("trainers")

	ash := Trainer{"Ash", 10, "Pallet Town"}
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 15, "Pewter City"}

	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	trainers := []interface{}{misty, brock}

	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

	filter := bson.D{{"name", "Ash"}}

	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// create a value into which the result can be decoded
	var result Trainer

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(2)

	// Here's an array in which you can store the decoded documents
	var results []*Trainer

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Trainer
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)

	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

	app := iris.Default()

	api := app.Party("/api")
	{
		v1 := api.Party("/1.0")
		{
			v1.Get("/{path:path}", func(ctx iris.Context) {
				pathPatam := ctx.Params().Get("path")

				path := strings.Split(pathPatam, "/")

				m := make(map[string]interface{})

				m["path"] = path[0]

				ctx.JSON(m)
			})
		}
	}

	app.Run(iris.Addr(":8080"))
}
