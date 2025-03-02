package models

import "time"

type Exercise struct {
	Year       int       `json:"year" bson:"year"`
	StartMonth string    `json:"startMonth" bson:"startMonth"`
	EndMonth   string    `json:"endMonth" bson:"endMonth"`
	DtAdd      time.Time `json:"dtAdd" bson:"dtAdd"`
}

type MovementAccount struct {
	DtMovement  int       `json:"dtMovement" bson:"dtMovement"`
	CodAccount  string    `json:"codAccount" bson:"codAccount"`
	DebitValue  string    `json:"debitValue" bson:"debitValue"`
	CreditValue time.Time `json:"creditValue" bson:"creditValue"`
}

type CostCenterSecondary struct {
	CodCostCenterSecondary string `json:"codCostCenterSecondary" bson:"codCostCenterSecondary"`
	Description            string `json:"description" bson:"description"`
}

type CostCenterCOA struct {
	IdCostCenter         string   `json:"idCostCenter" bson:"idCostCenter"`
	CostCentersSecondary []string `json:"costCentersSecondary" bson:"costCentersSecondary"`
}
