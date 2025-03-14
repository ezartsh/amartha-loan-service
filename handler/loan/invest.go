package loan

import (
	"encoding/json"
	"errors"
	"loan-service/model"
	"loan-service/pkg/request"
	"loan-service/pkg/response"
	"loan-service/pkg/validation"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (c Controller) Invest(resp *response.HttpResponse, req *request.HttpRequest) {

	dataRequest := InvestLoanRequest{}
	vars := mux.Vars(req.HttpRequest())

	var loanId uint64

	intId, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		resp.Error(http.StatusBadRequest, errors.New("path identifier is not in valid format"))
		return
	}

	loanId = intId

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

	formValidation := req.Validation(validation.SchemaRules{
		"InvestorId": []string{validation.Required, validation.Numeric},
		"Amount":     []string{validation.Required, validation.Numeric},
	})

	errorBags, err := formValidation.Validate(dataRequest)

	if err != nil {
		resp.ErrorWithData(http.StatusUnprocessableEntity, err.Error(), errorBags)
		return
	}

	approvedLoan := model.Loan{
		ID:                 loanId,
		BorrowerId:         1,
		PrincipalAmount:    100,
		ApprovalEvidence:   "",
		ApprovalEmployeeId: 1,
		Investors: []model.LoanInvestor{
			{
				InvestorId: dataRequest.InvestorId,
				Amount:     dataRequest.Amount,
			},
		},
		Status: model.LoanStateInvested,
	}

	resp.Json(http.StatusOK, approvedLoan)
	return
}
