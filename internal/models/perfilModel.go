package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Perfil struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	ShortName string             `json:"shortName" bson:"shortName"`
}
