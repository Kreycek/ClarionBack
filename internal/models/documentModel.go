package models

import (
	"time"
)

type Document struct {
	CodDocument string    `json:"codDocument" bson:"codDocument"`
	Description string    `json:"description" bson:"description"`
	DtAdd       time.Time `json:"dtAdd" bson:"dtAdd"`
}

type DocumentCompany struct {
	NameDocument      string    `json:"nameDocument" bson:"nameDocument"`
	DocumentNumber    string    `json:"documentNumber" bson:"documentNumber"`
	DtAdd             time.Time `json:"dtAdd" bson:"dtAdd"`
	Country           string    `json:"country" bson:"country,omitempty"`
	City              string    `json:"city" bson:"city,omitempty"`
	CodPostal         string    `json:"codPostal" bson:"codPostal,omitempty"`
	Address           string    `json:"address" bson:"address,omitempty"`
	AddressNumber     string    `json:"addressNumber" bson:"addressNumber,omitempty"`
	AddressComplement string    `json:"addressComplement" bson:"addressComplement,omitempty"`
	WorkNow           bool      `json:"workNow" bson:"workNow,omitempty"`
}
