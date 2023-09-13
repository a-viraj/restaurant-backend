package routes

import (
	controller "github.com/aviraj/resturant-management/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.GET("/users", controller.GetUsers())
	r.GET("/users/:userId", controller.GetUser())
	r.POST("/users/signup", controller.Signup())
	r.POST("/users/login", controller.Login())
}
