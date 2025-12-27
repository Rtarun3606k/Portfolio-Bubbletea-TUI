package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"portfolioTUI/config"
	"time"
)

// 1. Declare Global Client Variable
// usage: database.Client
var Client *mongo.Client

func ConnectToDataBase() {
	fmt.Println("Connecting to the database...")

	// 2. Validate Config
	dbURL := config.DATABASEURL
	if dbURL == "" {
		log.Fatal("Database URL is empty in .env")
	}

	// 3. Set Client Options
	opts := options.Client().ApplyURI(dbURL)

	// 4. Connect
	// We create a temporary context just for the connection handshake
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	Client, err = mongo.Connect(opts)
	if err != nil {
		log.Fatal("Error creating client:", err)
	}

	// 5. Ping to verify connection
	// Ping is the safest way to check.
	if err := Client.Ping(ctx, nil); err != nil {
		log.Fatal("Could not ping MongoDB:", err)
	}

	log.Println("✅ Connected to the database successfully!")
}

// Helper function to get a collection easily
func GetCollection(databaseName, collectionName string) *mongo.Collection {
	// Replace "Cluster0" with your actual Database Name if different
	return Client.Database(databaseName).Collection(collectionName)
}

func GetALLFromCollection(dataBasename, CollectionName string) ([]bson.M, error) {
	coll := GetCollection(dataBasename, CollectionName)

	// In V2, you can pass context.TODO() if you don't have a specific context
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Println(err)
	}

	// fmt.Println(results[0:3])
	return results, nil
}
