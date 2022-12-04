package mongodb

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDatabase *mongo.Database
var MongoCtx context.Context
var MongoClient *mongo.Client

func Connect() {
	println("Connecting to MongoDB Atlas")

	MONGO_URI := os.Getenv("MONGO_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	// assign globally
	MongoClient = client
	MongoDatabase = client.Database("blog")

	println("Successfully connected to MongoDB Atlas")
}

func Disconnect() {
	err := MongoClient.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}

func GetMongoConnection() (*mongo.Database, context.Context) {
	return MongoDatabase, MongoCtx
}

func Collection(name string) (*mongo.Collection, context.Context) {
	return MongoDatabase.Collection(name), MongoCtx
}
