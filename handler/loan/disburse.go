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
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (c Controller) Disburse(resp *response.HttpResponse, req *request.HttpRequest) {

	eventTime := time.Now().In(utils.JakartaTimeLocation)
	dataRequest := DisburseLoanRequest{}
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
		resp.Error(http.StatusNotFound, errors.New("disburse can only available when status is invested."))
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

	formValidation := req.Validation(validation.SchemaRules{
		"EmployeeId":      []string{validation.Required, validation.Numeric},
		"EvidencePicture": []string{validation.Required, validation.File},
	})

	errorBags, err := formValidation.Validate(dataRequest)

	if err != nil {
		resp.ErrorWithData(http.StatusUnprocessableEntity, err.Error(), errorBags)
		return
	}

	existingLoan.Status = model.LoanStateDisbursed
	existingLoan.DisburseDate = &utils.LocalTime{Time: eventTime}
	existingLoan.DisburseEmployeeId = 1

	Loans[existingLoan.ID-1] = *existingLoan

	resp.Json(http.StatusOK, existingLoan)
	return
}
