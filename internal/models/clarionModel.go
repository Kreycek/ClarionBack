package models

type EmailRequest struct {
	Email string `json:"email"`
}

type ChartOfAccountVerifyExistRequest struct {
	CodAccount string `json:"codAccount"`
}

type DailyVerifyExistRequest struct {
	CodDaily string `json:"codDaily"`
}
