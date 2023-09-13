package routes

import (
	controller "github.com/aviraj/resturant-management/controllers"
	"github.com/gin-gonic/gin"
)

func FoodRoutes(r *gin.Engine) {
	r.GET("/foods", controller.GetFoods())
	r.GET("/foods/:foodId", controller.GetFood())
	r.POST("/foods", controller.CreateFood())
	r.PATCH("/foods/:foodId", controller.UpdateFood())
}
