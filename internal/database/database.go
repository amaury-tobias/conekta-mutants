package database

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MutantCollection *mongo.Collection

func Init() error {
	clientOptions := options.Client().ApplyURI("mongodb://" + getHost() + ":27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	MutantCollection = client.Database("conekta").Collection("humans")
	return nil
}

func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	MutantCollection.Database().Client().Disconnect(ctx)
}

func getHost() string {
	dbHost := os.Getenv("DB_HOST")
	if len(dbHost) == 0 {
		dbHost = "localhost"
	}
	return dbHost
}
