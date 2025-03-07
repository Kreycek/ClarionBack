package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movement struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	DtMovimento   string             `json:"dtMovimento" bson:"dtMovimento"`
	CodDaily      string             `json:"codDaily" bson:"codDaily"`
	CodDocument   string             `json:"codDocument" bson:"codDocument"`
	Movements     []Movements        `json:"movements" bson:"movements"`
	IVA           string             `json:"iva" bson:"iva"`
	Active        bool               `json:"active" bson:"active,omitempty"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
	IdUserCreated string             `json:"idUserCreated" bson:"idUserCreated,omitempty"`
	IdUserUpdate  string             `json:"idUserUpdate" bson:"idUserUpdate,omitempty"`
}
