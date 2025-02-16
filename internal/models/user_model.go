package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name" bson:"name"`
	LastName       string             `json:"lastName" bson:"lastName,omitempty"`
	Email          string             `json:"email" bson:"email"`
	PassportNumber string             `json:"passportNumber" bson:"passportNumber,omitempty"`
	Password       string             `json:"password" bson:"password"`
	Perfil         []int              `json:"perfil" bson:"perfil"`
	Username       string             `json:"username" bson:"username,omitempty"`
}
