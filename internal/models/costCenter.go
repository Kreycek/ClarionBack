package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CostCenter struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CodCostCenter string             `json:"codCostCenter" bson:"codCostCenter,omitempty"`
	Description   string             `json:"description" bson:"description,omitempty"`
	CostCenterSub []CostCenterSub    `json:"costCenterSub" bson:"costCenterSub,omitempty"`
	Active        bool               `json:"active" bson:"active,omitempty"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
	IdUserCreated string             `json:"idUserCreated" bson:"idUserCreated,omitempty"`
	IdUserUpdate  string             `json:"idUserUpdate" bson:"idUserUpdate,omitempty"`
}
