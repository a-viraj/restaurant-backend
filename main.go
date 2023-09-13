package main

import (
	"os"

	database "github.com/aviraj/resturant-management/database"
	middleware "github.com/aviraj/resturant-management/middleware"
	routes "github.com/aviraj/resturant-management/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	r := gin.New()
	r.Use(gin.Logger())
	routes.UserRoutes(r)
	r.Use(middleware.Authenticate())

	routes.FoodRoutes(r)
	routes.MenuRoutes(r)
	routes.TableRoutes(r)
	routes.OrderRoutes(r)
	routes.OrderItemRoutes(r)
	routes.InvoiceRoutes(r)

	r.Run(":" + port)
}
