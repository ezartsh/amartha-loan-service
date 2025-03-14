package loan

import (
	"encoding/json"
	"errors"
	"loan-service/model"
	"loan-service/pkg/request"
	"loan-service/pkg/response"
	"loan-service/pkg/validation"
	"net/http"
)

func (c Controller) Create(resp *response.HttpResponse, req *request.HttpRequest) {

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

	formValidation := req.Validation(validation.SchemaRules{
		"BorrowerId":      []string{validation.Required, validation.Numeric},
		"PrincipalAmount": []string{validation.Required, validation.Numeric},
		"Rate":            []string{validation.Required, validation.Numeric},
	})

	errorBags, err := formValidation.Validate(dataRequest)

	if err != nil {
		resp.ErrorWithData(http.StatusUnprocessableEntity, err.Error(), errorBags)
		return
	}

	newLoan := model.Loan{
		BorrowerId:      dataRequest.BorrowerId,
		PrincipalAmount: dataRequest.PrincipalAmount,
		Status:          model.LoanStateProposed,
	}

	resp.Json(http.StatusOK, newLoan)
	return
}
