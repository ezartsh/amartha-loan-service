package router

import (
	"loan-service/pkg/request"
	"loan-service/pkg/response"
	"net/http"
)

type Request func(w *response.HttpResponse, r *request.HttpRequest)

func requestHandler(handler Request) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpRequest := request.NewRequest(r)
		httpResponse := response.NewResponse(w)
		// appContext := ctxa.App{}
		// if err := appContext.Store(r); err != nil {
		// 	util.RespondError(w, http.StatusUnauthorized, err)
		// 	return
		// }
		handler(httpResponse, httpRequest)
	}
}
