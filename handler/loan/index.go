package loan

import (
	"errors"
	"loan-service/model"
	"loan-service/pkg/request"
	"loan-service/pkg/response"
	"net/http"
	"strconv"

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
