package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents the structure of a MongoDB document
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Username     string             `bson:"username"`
	Email        string             `bson:"email"`
	Password     string             `bson:"password"`
	CreatedAt    time.Time          `bson:"createdAt"`
	LastLogin    time.Time          `bson:"lastLogin"`
	Preferences  UserPreferences    `bson:"preferences"`
	GithubID     string             `bson:"githubId,omitempty"`
}

// UserPreferences represents user's IDE preferences
type UserPreferences struct {
	Theme    string `bson:"theme"`
	FontSize int    `bson:"fontSize"`
	// Add other IDE preferences as needed
}