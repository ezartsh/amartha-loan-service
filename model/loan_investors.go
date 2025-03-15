package model

type LoanInvestor struct {
	Id         int     `json:"id"`
	InvestorId int     `json:"investor_id"`
	Amount     float64 `json:"amount"`
}
