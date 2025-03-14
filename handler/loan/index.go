package loan

import (
	"loan-service/model"
	"loan-service/pkg/request"
	"loan-service/pkg/response"
	"net/http"
)

func (c Controller) Index(resp *response.HttpResponse, req *request.HttpRequest) {

	loan := model.Loan{}

	if err := req.GetFormRequest(&loan); err != nil {
		resp.Error(http.StatusBadRequest, err)
		return
	}

	resp.Json(http.StatusOK, loan)
	return
}
