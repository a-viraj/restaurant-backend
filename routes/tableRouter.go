package routes

import (
	controller "github.com/aviraj/resturant-management/controllers"
	"github.com/gin-gonic/gin"
)

func TableRoutes(r *gin.Engine) {

	r.GET("/tables", controller.GetTables())
	r.GET("/tables/:tableId", controller.GetTable())
	r.POST("/tables", controller.CreateTable())
	r.PATCH("/tables/:tableId", controller.UpdateTable())
}
