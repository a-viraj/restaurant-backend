package routes

import (
	controller "github.com/aviraj/resturant-management/controllers"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(r *gin.Engine) {

	r.GET("/orders", controller.GetOrders())
	r.GET("/orders/:orderId", controller.GetOrder())
	r.POST("/orders", controller.CreateOrder())
	r.PATCH("/orders/:orderId", controller.UpdateOrder())
}
