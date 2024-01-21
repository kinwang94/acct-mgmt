package db

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDB() (*mongo.Collection, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
	)
	log.Println("Connecting to database: ", uri)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a connection to the database")
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, errors.Wrap(err, "failed to verify that the client can connect to the database")
	}
	log.Println("Connected successfully.")

	collection := client.Database("management").Collection("account")
	return collection, nil
}
