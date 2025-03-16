package loan

import (
	"loan-service/model"
	"loan-service/pkg/mail"
)

type Controller struct {
	mailHandler *mail.Mail
}

var Loans = []model.Loan{}

func NewController(mailHandler *mail.Mail) Controller {
	return Controller{
		mailHandler: mailHandler,
	}
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
