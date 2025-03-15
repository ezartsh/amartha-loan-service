package loan

import (
	"loan-service/pkg/router"
)

func RegisterRoutes(httpRouter *router.Router) {
	controller := NewController()
	httpRouter.Group("/loan", func(route router.GroupRoute) {
		route.Get("/", controller.Index)
		route.Get("/{id}", controller.GetById)
		route.Post("/create", controller.Create)
		route.Put("/{id}/approve", controller.Approve)
		route.Put("/{id}/invest", controller.Invest)
		route.Put("/{id}/disburse", controller.Disburse)
	})
}
