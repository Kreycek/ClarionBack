package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Company struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CodCompany      string             `json:"codCompany" bson:"codCompany,omitempty"`
	Name            string             `json:"name" bson:"name,omitempty"`
	CAE             string             `json:"cae" bson:"cae,omitempty"`
	Documents       []DocumentCompany  `json:"documents" bson:"documents,omitempty"`
	MainActivity    string             `json:"mainActivity" bson:"mainActivity,omitempty"`
	OtherActivities string             `json:"otherActivities" bson:"otherActivities,omitempty"`
	LegalNature     string             `json:"legalNature" bson:"legalNature,omitempty"`
	SocialCapital   string             `json:"socialCapital" bson:"socialCapital,omitempty"`
	NationalCapital string             `json:"nationalCapital" bson:"nationalCapital,omitempty"`
	ExtraCapital    string             `json:"extraCapital" bson:"extraCapital,omitempty"`
	PublicCapital   string             `json:"publicCapital" bson:"publicCapital,omitempty"`
	VATRegime       string             `json:"vatRegime" bson:"vatRegime,omitempty"`
	Email           string             `json:"email" bson:"email,omitempty"`
	WebSite         string             `json:"webSite" bson:"webSite,omitempty"`
	Active          bool               `json:"active" bson:"active,omitempty"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
	IdUserCreated   string             `json:"idUserCreated" bson:"idUserCreated,omitempty"`
	IdUserUpdate    string             `json:"idUserUpdate" bson:"idUserUpdate,omitempty"`
	Phone           []Phone            `json:"phone" bson:"phone,omitempty"`
	Exercise        []Exercise         `json:"exercise" bson:"exercise,omitempty"`
}
