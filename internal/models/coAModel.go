package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Plano de contas ou em inglês char of acount
type ChartOfAccount struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CodAccount     string             `json:"codAccount" bson:"codAccount"`
	Description    string             `json:"description" bson:"description,omitempty"`
	Year           []int              `json:"year" bson:"year"`
	Type           string             `json:"type" bson:"type,omitempty"`
	Active         bool               `json:"active" bson:"active"`
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt" bson:"updatedAt"`
	IdUserCreated  string             `json:"idUserCreated" bson:"idUserCreated,omitempty"`
	IdUserUpdate   string             `json:"idUserUpdate" bson:"idUserUpdate,omitempty"`
	CostCentersCOA []CostCenterCOA    `json:"costCentersCOA" bson:"costCentersCOA,omitempty"`
}
