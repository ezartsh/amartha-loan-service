package model

import "loan-service/utils"

type LoanInvestor struct {
	Id         int             `json:"id"`
	InvestorId int             `json:"investor_id"`
	Amount     float64         `json:"amount"`
	CreatedAt  utils.LocalTime `json:"created_at"`
	UpdatedAt  utils.LocalTime `json:"updated_at"`
}
