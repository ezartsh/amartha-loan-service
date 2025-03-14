package model

import "loan-service/utils"

type LoanState int

const (
	LoanStateProposed LoanState = iota
	LoanStateApproved
	LoanStateInvested
	LoanStateDisbursed
)

func (s LoanState) ToString() string {
	return []string{"PROPOSED", "APPROVED", "INVESTED", "DISBURSED"}[s]
}

type Borrower struct {
	Id uint
}

type Loan struct {
	ID                 uint64           `json:"id"`
	BorrowerId         uint             `json:"borrower_id"`
	PrincipalAmount    float64          `json:"principal_amount"`
	Rate               float64          `json:"rate"`
	ApprovalEvidence   string           `json:"approval_evidence"`
	ApprovalEmployeeId uint             `json:"approval_employee_id"`
	ApprovalDate       *utils.LocalTime `json:"approval_date"`
	DisburseEvidence   string           `json:"disburse_evidence"`
	DisburseEmployeeId uint             `json:"disburse_employee_id"`
	DisburseDate       *utils.LocalTime `json:"disburse_date"`
	Investors          []LoanInvestor   `json:"loan_investors"`
	Status             LoanState        `json:"status"`
}
