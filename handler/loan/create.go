package loan

import (
	"encoding/json"
	"errors"
	"loan-service/config"
	"loan-service/model"
	"loan-service/pkg/request"
	"loan-service/pkg/response"
	"loan-service/pkg/validation"
	"loan-service/utils"
	"net/http"
	"strconv"
	"time"
)

func (c Controller) Create(resp *response.HttpResponse, req *request.HttpRequest) {

	eventTime := time.Now().In(utils.JakartaTimeLocation)
	dataRequest := CreateLoanRequest{}

	if err := req.GetFormRequest(&dataRequest); err != nil {
		switch err.(type) {
		case *json.UnmarshalTypeError:
			break
		default:
			resp.Error(http.StatusBadRequest, errors.New("data request is not in valid format"))
			return
		}
	}

	// assuming that the roi >= company margin profit
	// the interest rate value must be (margin profit * 2) at minimum

	var minimumInterestRate string = "gte=" + strconv.FormatFloat((config.Env.CompanyMarginInPercent*2), 'f', -1, 64)
	var minimumLoanAmount string = "gte=" + strconv.FormatFloat((config.Env.MinimumLoanAmout), 'f', -1, 64)
	var maximumLoanAmount string = "lte=" + strconv.FormatFloat((config.Env.MaximumLoanAmout), 'f', -1, 64)

	formValidation := req.Validation(validation.SchemaRules{
		"BorrowerId":      []string{validation.Required, validation.Numeric},
		"PrincipalAmount": []string{validation.Required, validation.Numeric, minimumLoanAmount, maximumLoanAmount},
		"InterestRate":    []string{validation.Required, validation.Numeric, minimumInterestRate},
	})

	errorBags, err := formValidation.Validate(dataRequest)

	if err != nil {
		resp.ErrorWithData(http.StatusUnprocessableEntity, err.Error(), errorBags)
		return
	}

	newLoan := model.Loan{
		ID:              len(Loans) + 1,
		BorrowerId:      dataRequest.BorrowerId,
		PrincipalAmount: dataRequest.PrincipalAmount,
		InterestRate:    dataRequest.InterestRate,
		Roi:             dataRequest.InterestRate - config.Env.CompanyMarginInPercent,
		Status:          model.LoanStateProposed,
		Investors:       []model.LoanInvestor{},
		CreatedAt:       utils.LocalTime{Time: eventTime},
		UpdatedAt:       utils.LocalTime{Time: eventTime},
	}

	Loans = append([]model.Loan{newLoan}, Loans...)

	resp.Json(http.StatusCreated, newLoan)
	return
}
