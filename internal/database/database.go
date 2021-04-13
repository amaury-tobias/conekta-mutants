package database

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

var DBClient MongoClient

func NewMongoDatabase(databaseSession Session) (MongoDatabase, error) {
	clientOptions := options.Client().ApplyURI("mongodb://" + getHost() + ":27017")
	client, err := databaseSession.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return client.Database("conekta"), nil
}

func getHost() string {
	dbHost := os.Getenv("DB_HOST")
	if len(dbHost) == 0 {
		dbHost = "localhost"
	}
	return dbHost
}
