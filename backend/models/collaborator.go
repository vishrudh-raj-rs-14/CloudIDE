package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Collaborator struct {
	UserID primitive.ObjectID `bson:"userId"`
	Role   string             `bson:"role"`
}