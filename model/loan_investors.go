package model

type LoanInvestor struct {
	Id         uint    `json:"id"`
	InvestorId uint    `json:"investor_id"`
	Amount     float64 `json:"amount"`
}
