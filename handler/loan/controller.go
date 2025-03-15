package loan

import (
	"loan-service/model"
)

type Controller struct {
}

var Loans = []model.Loan{}

func NewController() Controller {
	return Controller{}
}

func searchLoanById(id int) *model.Loan {

	var existingLoan *model.Loan

	for _, loan := range Loans {
		if loan.ID == id {
			existingLoan = &loan
			break
		}
	}
	return existingLoan
}
