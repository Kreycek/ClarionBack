package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Daily struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	CodDaily      string             `json:"codDaily" bson:"codDaily"`
	Description   string             `json:"description" bson:"description"`
	Documents     []Document         `json:"documents" bson:"documents"`
	Active        bool               `json:"active" bson:"active"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"`
	IdUserCreated string             `json:"idUserCreated" bson:"idUserCreated,omitempty"`
	IdUserUpdate  string             `json:"idUserUpdate" bson:"idUserUpdate,omitempty"`
}
