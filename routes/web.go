package routes

import (
	"loan-service/handler/loan"
	"loan-service/pkg/router"
	"net/http"
)

func WebInit() http.Handler {
	httpRouter := router.NewRouter()

	loan.RegisterRoutes(httpRouter)

	return httpRouter.GetContainer()
}
