package utils

import (
	"context"
	"fmt"
	"time"

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

func CreateRepl(repl *models.Repl) (primitive.ObjectID, error) {
	client := models.Mg.Client
	collection := client.Database(models.Mg.Db.Name()).Collection("repls")
	res, err := collection.InsertOne(context.TODO(), repl)
	insertedID := res.InsertedID.(primitive.ObjectID)
	return insertedID, err
}

func GetRepl(replID primitive.ObjectID) (*models.Repl, error) {
	client := models.Mg.Client
	collection := client.Database(models.Mg.Db.Name()).Collection("repls")
	var repl models.Repl
	filter := bson.D{{"_id", replID}}
	err := collection.FindOne(context.TODO(), filter).Decode(&repl)
	return &repl, err
}

func GetReplsByUserID(userID primitive.ObjectID) ([]models.Repl, error) {
	client := models.Mg.Client
	collection := client.Database(models.Mg.Db.Name()).Collection("repls")
	var repls []models.Repl
	filter := bson.D{{"ownerId", userID}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var repl models.Repl
		if err := cursor.Decode(&repl); err != nil {
			return nil, err
		}
		repls = append(repls, repl)
	}
	return repls, nil
}
func UpdateRepl(replIDHex primitive.ObjectID, repl *models.Repl) error {
	collection := models.Mg.Db.Collection("repls")
	ctx := context.Background()
	update := bson.M{
		"$set": bson.M{
			"name":          repl.Name,
			"description":   repl.Description,
			"collaborators": repl.Collaborators,
			"updatedAt":     time.Now(),
			"language":      repl.Language,
			"framework":     repl.Framework,
			"visibility":    repl.Visibility,
			"status":        repl.Status,
			"containerID":   repl.ContainerID,
			"containerPort": repl.ContainerPort,
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": replIDHex}, update)
	if err != nil {
		return fmt.Errorf("failed to update repl: %v", err)
	}

	return nil
}

func DeleteRepl(replID primitive.ObjectID) error {
	client := models.Mg.Client
	collection := client.Database(models.Mg.Db.Name()).Collection("repls")
	filter := bson.D{{"_id", replID}}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("repl with ID %s does not exist", replID.Hex())
	}
	return nil
}

func ReplExists(replID primitive.ObjectID) (bool, error) {
	client := models.Mg.Client
	collection := client.Database(models.Mg.Db.Name()).Collection("repls")
	filter := bson.D{{"_id", replID}}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
