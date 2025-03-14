package loan

import "mime/multipart"

type CreateLoanRequest struct {
	BorrowerId      uint    `json:"borrower_id"`
	PrincipalAmount float64 `json:"principal_amount"`
	Rate            float64 `json:"rate"`
}

type ApproveLoanRequest struct {
	EmployeeId      uint                 `json:"employee_id"`
	EvidencePicture multipart.FileHeader `json:"evidence_picture"`
}

type InvestLoanRequest struct {
	InvestorId uint    `json:"investor_id"`
	Amount     float64 `json:"amount"`
}

type DisburseLoanRequest struct {
	EmployeeId      uint                 `json:"employee_id"`
	EvidencePicture multipart.FileHeader `json:"evidence_picture"`
}
