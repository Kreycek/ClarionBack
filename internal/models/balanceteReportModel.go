package models

import (
	"time"
)

type Balancete struct {
	DateGeration       time.Time `json:"dateGeration" bson:"dateGeration,omitempt"`
	CodAccount         string    `json:"codAccount" bson:"codAccount,omitempty"`
	Description        string    `json:"description" bson:"description,omitempty"`
	DebitValue         float64   `json:"debitValue" bson:"debitValue,omitempty"`
	CreditValue        float64   `json:"creditValue" bson:"creditValue,omitempty"`
	BalanceDebitValue  float64   `json:"balanceDebitValue" bson:"balanceDebitValue,omitempty"`
	BalandeCreditValue float64   `json:"balanceCredittValue" bson:"balanceCredittValue,omitempty"`
	FatherCod          string    `json:"fatherCod" bson:"fatherCod,omitempty"`
	Class              string    `json:"class" bson:"class,omitempty"`
	Sum                bool      `json:"sum" bson:"sum,omitempty"`
}
