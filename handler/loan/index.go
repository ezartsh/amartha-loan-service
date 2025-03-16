package loan

import (
	"errors"
	"loan-service/model"
	"loan-service/pkg/request"
	"loan-service/pkg/response"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func (c Controller) Index(resp *response.HttpResponse, req *request.HttpRequest) {

	resp.Json(http.StatusOK, Loans)
	return
}

func (c Controller) GetById(resp *response.HttpResponse, req *request.HttpRequest) {

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

	resp.Json(http.StatusOK, existingLoan)
	return
}

func (c Controller) AgreementLetter(resp *response.HttpResponse, req *request.HttpRequest) {

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

	if existingLoan.AgreementLetter == nil {
		resp.Error(http.StatusNotFound, errors.New("agreement letter not found"))
		return
	}

	// Ensure the file exists
	data, err := os.ReadFile(*existingLoan.AgreementLetter)
	if err != nil {
		resp.Error(http.StatusNotFound, errors.New("agreement letter not found"))
		return
	}
	if _, err := os.Stat(*existingLoan.AgreementLetter); os.IsNotExist(err) {
		resp.Error(http.StatusNotFound, errors.New("agreement letter not found"))
		return
	}

	paths := strings.Split(*existingLoan.AgreementLetter, "/")
	fileName := paths[len(paths)-1]

	// Set headers for inline display in the browser
	resp.Writer().Header().Set("Content-Type", "application/pdf")
	resp.Writer().Header().Set("Content-Disposition", "inline; filename="+fileName)
	resp.Writer().Write(data)

	return
}
