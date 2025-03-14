package request

import (
	"loan-service/pkg/validation"
	"net/http"

	"github.com/ezartsh/inrequest"
)

type HttpRequest struct {
	http *http.Request
}

func NewRequest(r *http.Request) *HttpRequest {
	return &HttpRequest{
		http: r,
	}
}

func (r HttpRequest) GetFormRequest(data interface{}) error {
	req := inrequest.FormData(r.http)
	return req.ToBind(data)
}

func (r HttpRequest) Validation(schemaRules validation.SchemaRules) *validation.Valid {
	return validation.NewStructValidation(schemaRules)
}
func (r HttpRequest) HttpRequest() *http.Request {
	return r.http
}
