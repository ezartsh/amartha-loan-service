package loan

import "mime/multipart"

type CreateLoanRequest struct {
	BorrowerId      int     `json:"borrower_id"`
	PrincipalAmount float64 `json:"principal_amount"`
	InterestRate    float64 `json:"interest_rate"`
}

type ApproveLoanRequest struct {
	EmployeeId      int                   `json:"employee_id"`
	EvidencePicture *multipart.FileHeader `json:"evidence_picture"`
}

type InvestLoanRequest struct {
	InvestorId int     `json:"investor_id"`
	Amount     float64 `json:"amount"`
}

type DisburseLoanRequest struct {
	EmployeeId            int                   `json:"employee_id"`
	SignedAgreementLetter *multipart.FileHeader `json:"signed_agreement_letter"`
}
