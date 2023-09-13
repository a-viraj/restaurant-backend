package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aviraj/resturant-management/database"
	"github.com/aviraj/resturant-management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InvoiceViewFormat struct {
	Invoice_id       string
	Payment_method   string
	Order_id         string
	Payment_status   string
	Payment_due      interface{}
	Table_number     interface{}
	Payment_due_date time.Time
	Order_details    interface{}
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		res, err := invoiceCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing invoice items"})
		}
		var allinvoices []bson.M
		if err := res.All(ctx, allinvoices); err != nil {
			log.Fatal(err)
		}
		c.JSON(200, allinvoices)
	}
}
func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		invoiceId := c.Param("invoiceId")
		var invoice models.Invoice
		err := invoiceCollection.FindOne(ctx, bson.M{"invoiceId": invoiceId}).Decode(&invoice)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		var invoiceView InvoiceViewFormat
		allOrderItems, err := ItemsByOrder(invoice.OrderId)
		invoiceView.Order_id = invoice.InvoiceId
		invoiceView.Payment_due_date = invoice.PaymentDueDate
		invoiceView.Payment_method = "null"
		if invoice.PaymentMethod != nil {
			invoiceView.Payment_method = *invoice.PaymentMethod
		}
		invoiceView.Invoice_id = invoice.InvoiceId
		invoice.PaymentStatus = *&invoice.PaymentStatus
		invoiceView.Payment_due = allOrderItems[0]["payment_due"]
		invoiceView.Table_number = allOrderItems[0]["table_number"]
		invoiceView.Order_details = allOrderItems[0]["order_items"]
		c.JSON(200, invoiceView)
	}
}
func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var invoice models.Invoice
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"orderId": invoice.OrderId}).Decode(&order)
		if err != nil {
			msg := fmt.Sprintf("message:order not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		status := "PENDING"
		if invoice.PaymentStatus == nil {
			invoice.PaymentStatus = &status
		}
		invoice.PaymentDueDate, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
		invoice.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.ID=primitive.NewObjectID()
		invoice.InvoiceId=invoice.ID.Hex()
		res,err:=invoiceCollection.InsertOne(ctx,invoice)
		if err!=nil{
			msg := fmt.Sprintf("message:invoice item not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(200,res)
	}
}
func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var invoice models.Invoice
		invoiceId := c.Param("invoiceId")

		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		filter := bson.M{"invoiceId": invoiceId}
		var updateObj primitive.D
		if invoice.PaymentMethod != nil {
			updateObj = append(updateObj, bson.E{Key: "PaymentMethod", Value: invoice.PaymentMethod})
		}
		if invoice.PaymentStatus != nil {
			updateObj = append(updateObj, bson.E{Key: "PaymentStatus", Value: invoice.PaymentStatus})

		}
		invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: invoice.UpdatedAt})
		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		status := "PENDING"
		if invoice.PaymentStatus == nil {
			invoice.PaymentStatus = &status
		}

		res, err := invoiceCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key:"$set", Value:updateObj},
			},
			&opt,
		)
		if err != nil {
			msg := fmt.Sprintf("invoice item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(200, res)
	}
}
