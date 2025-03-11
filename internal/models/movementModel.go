package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movement struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Date            time.Time          `json:"date" bson:"date,omitempt"`
	CompanyFullData string             `json:"companyFullData" bson:"companyFullData"`
	CompanyId       primitive.ObjectID `json:"companyId,omitempty" bson:"companyId,omitempty"`
	CompanyDocument string             `json:"companyDocument" bson:"companyDocument"`
	CodDaily        string             `json:"codDaily" bson:"codDaily"`
	CodDocument     string             `json:"codDocument" bson:"codDocument"`
	Movements       []Movements        `json:"movements" bson:"movements,omitempty"`
	Active          bool               `json:"active" bson:"active,omitempty"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
	IdUserCreated   string             `json:"idUserCreated" bson:"idUserCreated,omitempty"`
	IdUserUpdate    string             `json:"idUserUpdate" bson:"idUserUpdate,omitempty"`
}
