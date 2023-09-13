package routes

import (
	controller "github.com/aviraj/resturant-management/controllers"
	"github.com/gin-gonic/gin"
)

func MenuRoutes(r *gin.Engine) {

	r.GET("/menus", controller.GetMenus())
	r.GET("/menus/:menusId", controller.GetMenu())
	r.POST("/menus", controller.CreateMenu())
	r.PATCH("/menus/:menusId", controller.UpdateMenu())
}
