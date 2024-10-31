package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repl struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string             `bson:"name"`
	Description   string             `bson:"description"`
	OwnerID       primitive.ObjectID `bson:"ownerId"`
	Collaborators []Collaborator     `bson:"collaborators"`
	CreatedAt     time.Time          `bson:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt"`
	Language      string             `bson:"language"`
	Framework     string             `bson:"framework"`
	Visibility    string             `bson:"visibility"`
	Status        string             `bson:"status"`
	ContainerID   string             `bson:"containerID"`
	ContainerPort string             `bson:"containerPort"`
	ImageName	  string 			 `bson:"imageName"`
}
