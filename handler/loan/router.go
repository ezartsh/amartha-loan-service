package loan

import (
	"loan-service/pkg/mail"
	"loan-service/pkg/router"
)

func RegisterRoutes(httpRouter *router.Router, mailHandler *mail.Mail) {
	controller := NewController(mailHandler)
	httpRouter.Group("/loans", func(route router.GroupRoute) {
		route.Get("/", controller.Index)
		route.Get("/{id}", controller.GetById)
		route.Post("/create", controller.Create)
		route.Put("/{id}/approve", controller.Approve)
		route.Put("/{id}/invest", controller.Invest)
		route.Put("/{id}/disburse", controller.Disburse)
		route.Get("/{id}/agreement-letter", controller.AgreementLetter)
		route.Post("/{id}/resend-email", controller.ResendEmail)
	})
}
