package models

import (
	"time"
)

type Phone struct {
	CodCountry  string    `json:"codCountry" bson:"codCountry"`
	CodState    string    `json:"codState" bson:"codState"`
	PhoneNumber string    `json:"phoneNumber" bson:"phoneNumber"`
	DtAdd       time.Time `json:"dtAdd" bson:"dtAdd"`
}
