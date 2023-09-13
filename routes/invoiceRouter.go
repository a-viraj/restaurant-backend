package routes

import (
	controller "github.com/aviraj/resturant-management/controllers"
	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(r *gin.Engine) {

	r.GET("/invoice", controller.GetInvoices())
	r.GET("/invoice/:invoiceId", controller.GetInvoice())
	r.POST("/invoice", controller.CreateInvoice())
	r.PATCH("/invoice/:invoiceId", controller.UpdateInvoice())
}
