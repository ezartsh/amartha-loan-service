package routes

import (
	"loan-service/handler/loan"
	"loan-service/pkg/mail"
	"loan-service/pkg/router"
	"net/http"
)

func WebInit(mailHandler *mail.Mail) http.Handler {
	httpRouter := router.NewRouter()

	loan.RegisterRoutes(httpRouter, mailHandler)

	return httpRouter.GetContainer()
}
