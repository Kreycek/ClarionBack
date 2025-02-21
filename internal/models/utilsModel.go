package models

import "time"

type Exercise struct {
	Year       int       `json:"ano" bson:"ano"`
	StartMonth string    `json:"startMonth" bson:"startMonth"`
	EndMonth   string    `json:"endMonth" bson:"endMonth"`
	DtAdd      time.Time `json:"dtAdd" bson:"dtAdd"`
}
