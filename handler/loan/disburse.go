package loan

import (
	"encoding/json"
	"errors"
	"io"
	"loan-service/model"
	"loan-service/pkg/logger"
	"loan-service/pkg/request"
	"loan-service/pkg/response"
	"loan-service/pkg/validation"
	"loan-service/utils"
	"net/http"
	"os"
	"path/filepath"
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
		"EmployeeId":            []string{validation.Required, validation.Numeric},
		"SignedAgreementLetter": []string{validation.Required},
	})

	errorBags, err := formValidation.Validate(dataRequest)

	if err != nil {
		resp.ErrorWithData(http.StatusUnprocessableEntity, err.Error(), errorBags)
		return
	}

	file, err := req.HttpRequest().MultipartForm.File["signed_agreement_letter"][0].Open()
	if err != nil {
		logger.AppLog.Error(err, "failed to open file from request")
		resp.Error(http.StatusInternalServerError, err)
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		logger.AppLog.Error(err, "failed to get os directory")
		resp.Error(http.StatusInternalServerError, err)
		return
	}

	fileLocation := filepath.Join(dir, "upload", utils.HashFileName(dataRequest.SignedAgreementLetter.Filename))
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.AppLog.Error(err, "failed to open target directory")
		resp.Error(http.StatusInternalServerError, err)
		return
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, file); err != nil {
		logger.AppLog.Error(err, "failed to store the file")
		resp.Error(http.StatusInternalServerError, err)
		return
	}

	existingLoan.DisburseEmployeeId = &dataRequest.EmployeeId
	existingLoan.DisburseDate = &utils.LocalTime{Time: eventTime}
	existingLoan.SignedAgreementLetter = &fileLocation
	existingLoan.UpdatedAt = utils.LocalTime{Time: eventTime}
	existingLoan.Status = model.LoanStateDisbursed

	Loans[existingLoan.ID-1] = *existingLoan

	resp.Json(http.StatusOK, existingLoan)
	return
}
