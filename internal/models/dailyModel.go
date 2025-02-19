package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Daily struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CodDaily      string             `json:"codDaily" bson:"codDaily,omitempty"`
	Description   string             `json:"description" bson:"description,omitempty"`
	Documents     []Document         `json:"documents" bson:"documents,omitempty"`
	Active        bool               `json:"active" bson:"active,omitempty"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
	IdUserCreated string             `json:"idUserCreated" bson:"idUserCreated,omitempty"`
	IdUserUpdate  string             `json:"idUserUpdate" bson:"idUserUpdate,omitempty"`
}
