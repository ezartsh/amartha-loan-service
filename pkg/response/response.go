package response

import (
	"encoding/json"
	"net/http"
)

type HttpResponse struct {
	w http.ResponseWriter
}

func NewResponse(w http.ResponseWriter) *HttpResponse {
	return &HttpResponse{
		w: w,
	}
}

func (r HttpResponse) Json(status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		r.w.WriteHeader(http.StatusInternalServerError)
		r.w.Write([]byte(err.Error()))
		return
	}

	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(status)
	r.w.Write(response)
}

func (r HttpResponse) Error(code int, message error) {
	r.Json(code, map[string]string{"error": message.Error()})
}

func (r HttpResponse) ErrorWithData(code int, message string, body interface{}) {
	r.Json(code, map[string]interface{}{"message": message, "errors": body})
}
