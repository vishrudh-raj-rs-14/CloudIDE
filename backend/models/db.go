package models

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var Mg MongoInstance;

func Connect(dbName string, mongoURI string) error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))

	if err != nil {
		return err
	}


	db := client.Database(dbName)

	if err != nil {
		return err
	}
	Mg = MongoInstance{
		Client: client,
		Db:     db,
	}
	return nil
}