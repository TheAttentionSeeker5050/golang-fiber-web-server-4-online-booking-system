package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// get configuration for MongoDB from the .env file
func GetMongoClient() *mongo.Client {
	SetEnvironmentFromFile()

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI is not set in .env file")
		return nil
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		return nil
	}

	return client
}

// a singleton function to return the database client instance using the database name from the .env file and the collection name as function parameter, return collection and error
func GetMongoCollection(mongoClient *mongo.Client, collectionName string) (*mongo.Collection, error) {

	// Get the database name from the .env file
	dbName := os.Getenv("MONGODB_DB_NAME")
	if dbName == "" {
		// log.Fatal("MONGODB_DB_NAME is not set in .env file")
		return nil, fmt.Errorf("MONGODB_DB_NAME is not set in .env file")
	}

	// Get the collection name from the function parameter
	collection := mongoClient.Database(dbName).Collection(collectionName)
	return collection, nil
}

func CloseMongoClientConnection(client *mongo.Client) {
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			fmt.Println("Error closing MongoDB connection")
		}
	}()
}
