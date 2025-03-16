package loan

import (
	"encoding/json"
	"errors"
	"fmt"
	"loan-service/config"
	"loan-service/model"
	"loan-service/pkg/logger"
	"loan-service/pkg/request"
	"loan-service/pkg/response"
	"loan-service/pkg/validation"
	"loan-service/utils"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
)

func (c Controller) Invest(resp *response.HttpResponse, req *request.HttpRequest) {

	eventTime := time.Now().In(utils.JakartaTimeLocation)
	dataRequest := InvestLoanRequest{}
	vars := mux.Vars(req.HttpRequest())

	var loanId int

	loanId, err := strconv.Atoi(vars["id"])

	if err != nil {
		resp.Error(http.StatusBadRequest, errors.New("path identifier is not in valid format"))
		return
	}

	var existingLoan *model.Loan = searchLoanById(loanId)

	if existingLoan == nil {
		resp.Error(http.StatusNotFound, errors.New("loan record not found"))
		return
	}

	if existingLoan.Status != model.LoanStateApproved {
		resp.Error(http.StatusNotFound, errors.New("invest can only available if status is approved."))
		return
	}

	var totalInvested float64
	var existingInvestors = []int{}

	for _, investor := range existingLoan.Investors {
		totalInvested += investor.Amount
		existingInvestors = append(existingInvestors, investor.InvestorId)
	}

	maximumAvailableAmount := existingLoan.PrincipalAmount - totalInvested

	if maximumAvailableAmount == 0 {
		resp.Error(http.StatusBadRequest, errors.New("cannot invest, already reach the limit."))
		return
	}

	if err := req.GetFormRequest(&dataRequest); err != nil {
		switch err.(type) {
		case *json.UnmarshalTypeError:
			break
		default:
			resp.Error(http.StatusBadRequest, errors.New("data request is not in valid format"))
			return
		}

		resp.Error(http.StatusBadRequest, err)
		return
	}

	var maximumInvestAmount string = "lte=" + strconv.FormatFloat(maximumAvailableAmount, 'f', -1, 64)

	formValidation := req.Validation(validation.SchemaRules{
		"InvestorId": []string{validation.Required, validation.Numeric},
		"Amount":     []string{validation.Required, validation.Numeric, maximumInvestAmount},
	})

	errorBags, err := formValidation.Validate(dataRequest)

	if slices.Contains(existingInvestors, dataRequest.InvestorId) {
		errorBags.Add("investor_id", "investor cannot invest more than one")
	}

	if err != nil || len(errorBags) > 0 {
		resp.ErrorWithData(http.StatusUnprocessableEntity, "validation input error", errorBags)
		return
	}

	existingLoan.Investors = append(existingLoan.Investors, model.LoanInvestor{
		Id:         len(existingLoan.Investors) + 1,
		InvestorId: dataRequest.InvestorId,
		Amount:     dataRequest.Amount,
		CreatedAt:  utils.LocalTime{Time: eventTime},
		UpdatedAt:  utils.LocalTime{Time: eventTime},
	})

	dir, err := os.Getwd()
	if err != nil {
		logger.AppLog.Error(err, "failed to get os directory")
		resp.Error(http.StatusInternalServerError, err)
		return
	}

	if dataRequest.Amount == maximumAvailableAmount {
		existingLoan.Status = model.LoanStateInvested

		fileLocation := filepath.Join(dir, "upload", utils.HashFileName(fmt.Sprintf("agreement_letter_loan_%d.pdf", existingLoan.ID)))

		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 16)
		pdf.Text(40, 20, fmt.Sprintf("Agreement Letter Loan id: %d", existingLoan.ID))
		pdf.Text(40, 40, "Content of Agreement Letter.")

		existingLoan.AgreementLetter = &fileLocation

		err := pdf.OutputFileAndClose(fileLocation)
		if err != nil {
			resp.Error(http.StatusInternalServerError, err)
			return
		}

		for _, investor := range existingLoan.Investors {
			receipient := fmt.Sprintf("investor_%d@example.com", investor.InvestorId)

			go func() {
				if err := c.mailHandler.Send(
					receipient,
					fmt.Sprintf("Agreement Letter Loan id : %d", existingLoan.ID),
					fmt.Sprintf(`
					<html>
						<body>
							<h1>Dear Investor %d,</h1>
							<p>Here is the <a href="%s" target="_blank">link</a> for the agreement letter.</p>
							<p>Thanks,<br>Company</p>
						</body>
					</html>
				`, investor.InvestorId, fmt.Sprintf("http://localhost:%d/loans/%d/agreement-letter", config.Env.HTTPPort, existingLoan.ID)),
				); err != nil {
					logger.AppLog.Errorf(err, "failed to send email to %s", receipient)
					return
				}
				logger.AppLog.Debugf("email successfuly sent to %s", receipient)
			}()
		}
	}

	Loans[existingLoan.ID-1] = *existingLoan

	resp.Json(http.StatusOK, existingLoan)
	return
}

func (c Controller) ResendEmail(resp *response.HttpResponse, req *request.HttpRequest) {

	vars := mux.Vars(req.HttpRequest())

	var loanId int

	loanId, err := strconv.Atoi(vars["id"])

	if err != nil {
		resp.Error(http.StatusBadRequest, errors.New("path identifier is not in valid format"))
		return
	}

	var existingLoan *model.Loan = searchLoanById(loanId)

	if existingLoan == nil {
		resp.Error(http.StatusNotFound, errors.New("loan record not found"))
		return
	}

	if existingLoan.Status != model.LoanStateInvested {
		resp.Error(http.StatusNotFound, errors.New("resent can only available if status is invested."))
		return
	}

	for _, investor := range existingLoan.Investors {
		receipient := fmt.Sprintf("investor_%d@example.com", investor.InvestorId)

		go func() {
			if err := c.mailHandler.Send(
				receipient,
				fmt.Sprintf("Agreement Letter Loan id : %d", existingLoan.ID),
				fmt.Sprintf(`
					<html>
						<body>
							<h1>Dear Investor %d,</h1>
							<p>Here is the <a href="%s" target="_blank">link</a> for the agreement letter.</p>
							<p>Thanks,<br>Company</p>
						</body>
					</html>
				`, investor.InvestorId, fmt.Sprintf("http://localhost:%d/loans/%d/agreement-letter", config.Env.HTTPPort, existingLoan.ID)),
			); err != nil {
				logger.AppLog.Errorf(err, "failed to send email to %s", receipient)
				return
			}

			logger.AppLog.Debugf("email successfuly sent to %s", receipient)
		}()
	}

	resp.Json(http.StatusOK, existingLoan)
	return
}
