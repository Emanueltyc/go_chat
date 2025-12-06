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

func Connect() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("MONGO_URI")
	db := os.Getenv("MONGO_DB")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("There was an error trying to connect to the database: ", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("Could not ping the database: ", err)
	}

	fmt.Println("Connected to Database succesfully!")
	
	return client.Database(db)
}