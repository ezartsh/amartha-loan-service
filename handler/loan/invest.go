package loan

import (
	"encoding/json"
	"errors"
	"loan-service/model"
	"loan-service/pkg/request"
	"loan-service/pkg/response"
	"loan-service/pkg/validation"
	"loan-service/utils"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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

	if dataRequest.Amount == maximumAvailableAmount {
		existingLoan.Status = model.LoanStateInvested
	}

	Loans[existingLoan.ID-1] = *existingLoan

	resp.Json(http.StatusOK, existingLoan)
	return
}
