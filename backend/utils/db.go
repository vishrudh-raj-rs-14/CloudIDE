package utils

import (
	"context"

	"github.com/vishrudh-raj-rs-14/cloudIDEbackend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Function to find a user by email
func FindUserByEmail(email string) (*models.User, error) {
	client := models.Mg.Client
	collection := client.Database(models.Mg.Db.Name()).Collection("users")

	var user models.User
	filter := bson.D{{"email", email}} // Create a filter to find the user by email

	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		return nil, err // Return any other errors
	}

	return &user, nil // Return the found user
}

func FindUserByUserName(username string) (*models.User, error) {
	client := models.Mg.Client
	collection := client.Database(models.Mg.Db.Name()).Collection("users")

	var user models.User
	filter := bson.D{{"username", username}} // Create a filter to find the user by email

	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		return nil, err // Return any other errors
	}

	return &user, nil // Return the found user
}

func CreateUser(user *models.User) (primitive.ObjectID,  error) {
	client := models.Mg.Client
	collection := client.Database(models.Mg.Db.Name()).Collection("users")
	res, err := collection.InsertOne(context.TODO(), user)
	insertedID := res.InsertedID.(primitive.ObjectID)
	return insertedID, err
}