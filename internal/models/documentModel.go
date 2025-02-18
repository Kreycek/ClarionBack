package models

import (
	"time"
)

type Document struct {
	CodDocument string    `json:"codDocument" bson:"codDocument"`
	Description string    `json:"description" bson:"description"`
	DtAdd       time.Time `json:"dtAdd" bson:"dtAdd"`
}
