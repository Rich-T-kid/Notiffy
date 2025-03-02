package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
The only purpose of this files is to properly Load data from disk [remote or local]
into memeory for this application to run properly this should
all be handling in this file alone. As of right now this is looking like a MongoDB implementation. just because its quick and easy.

Goes without saying but everything should rely on interfaces and not implementaion details
*/
// This is used to store application state, Each Service should have their own implemntaion of this
type Storage interface {
	Save() error // Save data to disk
	Load() error // Load data from disk into ram
}

// Example Template below
func NewStorage() Storage {
	return &StorageSolution{}
}

type StorageSolution struct {
}

func (s *StorageSolution) Save() error { return nil }
func (s *StorageSolution) Load() error { return nil }

func EstablishMongoConnection() *mongo.Database {
	mongoURI := os.Getenv("MONGO_URI")
	fmt.Printf("MongoURI:%s\n", mongoURI)

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Establish connection
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the database to check if it's reachable
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB connection test failed: %v", err)
	}
	db := client.Database("Notification-Service")
	return db
}
