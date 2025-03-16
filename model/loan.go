package model

import (
	"loan-service/utils"
)

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
	ID                 int              `json:"id"`
	BorrowerId         int              `json:"borrower_id"`
	PrincipalAmount    float64          `json:"principal_amount"`
	InterestRate       float64          `json:"interest_rate"`
	Roi                float64          `json:"roi"`
	ApprovalEvidence   string           `json:"approval_evidence"`
	ApprovalEmployeeId int              `json:"approval_employee_id"`
	ApprovalDate       *utils.LocalTime `json:"approval_date"`
	DisburseEvidence   string           `json:"disburse_evidence"`
	DisburseEmployeeId int              `json:"disburse_employee_id"`
	DisburseDate       *utils.LocalTime `json:"disburse_date"`
	Investors          []LoanInvestor   `json:"loan_investors"`
	Status             LoanState        `json:"status"`
	CreatedAt          utils.LocalTime  `json:"created_at"`
	UpdatedAt          utils.LocalTime  `json:"updated_at"`
}
