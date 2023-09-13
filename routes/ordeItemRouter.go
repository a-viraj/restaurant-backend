package routes

import (
	controller "github.com/aviraj/resturant-management/controllers"
	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(r *gin.Engine) {

	r.GET("/orderItems", controller.GetOrderItems())
	r.GET("/orderItems/:orderItemId", controller.GetOrderItem())
	r.POST("/orderItems", controller.CreateOrderItem())
	r.PATCH("/orderItems/:orderItemId", controller.UpdateOrderItem())
	r.GET("/orderItems-order/:orderId", controller.GetOrderItemsByOrder())
}
