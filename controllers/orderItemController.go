package controllers

import (
	"context"
	"github.com/aviraj/resturant-management/database"
	"github.com/aviraj/resturant-management/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderItemPack struct {
	Table_id    *string
	Order_items []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := orderItemCollection.Find(context.TODO(), bson.M{})

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing ordered items"})
			return
		}
		var allOrderItems []bson.M
		if err = result.All(ctx, &allOrderItems); err != nil {
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, allOrderItems)
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("orderid")

		allOrderItems, err := ItemsByOrder(orderId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order items by order ID"})
			return
		}
		c.JSON(http.StatusOK, allOrderItems)
	}
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	matchStage := bson.D{{Key:"$match", Value:bson.D{{Key:"orderid",Value: id}}}}
	lookupStage := bson.D{{Key:"$lookup", Value:bson.D{{Key:"from",Value: "food"}, {Key:"localField", Value:"foodid"}, {Key:"foreignField", Value:"foodid"}, {Key:"as",Value: "food"}}}}
	unwindStage := bson.D{{Key:"$unwind", Value:bson.D{{Key:"path",Value: "$food"}, {Key:"preserveNullAndEmptyArrays",Value: true}}}}

	lookupOrderStage := bson.D{{Key:"$lookup",Value: bson.D{{Key:"from", Value:"order"}, {Key:"localField",Value: "orderid"}, {Key:"foreignField", Value:"orderid"}, {Key:"as", Value:"order"}}}}
	unwindOrderStage := bson.D{{Key:"$unwind",Value: bson.D{{Key:"path",Value: "$order"}, {Key:"preserveNullAndEmptyArrays",Value: true}}}}

	lookupTableStage := bson.D{{Key:"$lookup", Value:bson.D{{Key:"from", Value:"table"}, {Key:"localField",Value: "order.tableid"}, {Key:"foreignField", Value:"tableid"}, {Key:"as",Value: "table"}}}}
	unwindTableStage := bson.D{{Key:"$unwind", Value:bson.D{{Key:"path",Value: "$table"}, {Key:"preserveNullAndEmptyArrays",Value: true}}}}

	projectStage := bson.D{
		{Key:"$project",Value: bson.D{
			{Key:"id",Value: 0},
			{Key:"amount",Value: "$food.price"},
			{Key:"total_count",Value: 1},
			{Key:"foodname",Value: "$food.name"},
			{Key:"food_image", Value:"$food.food_image"},
			{Key:"tablenumber", Value:"$table.tablenumber"},
			{Key:"tableid",Value: "$table.trableid"},
			{Key:"orderid",Value: "$order.orderid"},
			{Key:"price", Value:"$food.price"},
			{Key:"quantity",Value: 1},
		}}}

	groupStage := bson.D{{Key:"$group",Value: bson.D{{Key:"id", Value:bson.D{{Key:"orderid",Value: "$orderid"}, {Key:"tableid",Value: "$tableid"}, {Key:"tablenumber", Value:"$tablenumber"}}}, {Key:"paymentdue",Value: bson.D{{Key:"$sum",Value: "$amount"}}}, {Key:"totalcount",Value: bson.D{{Key:"$sum",Value: 1}}}, {Key:"orderitems",Value: bson.D{{Key:"$push",Value: "$$ROOT"}}}}}}

	projectStage2 := bson.D{
		{Key:"$project",Value: bson.D{

			{Key:"id",Value: 0},
			{Key:"payment_due",Value: 1},
			{Key:"total_count",Value: 1},
			{Key:"table_number",Value: "$id.tablenumber"},
			{Key:"orderitems",Value: 1},
		}}}

	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2})

	if err != nil {
		panic(err)
	}

	if err = result.All(ctx, &OrderItems); err != nil {
		panic(err)
	}

	defer cancel()

	return OrderItems, err

}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		orderItemId := c.Param("orderitemid")
		var orderItem models.OrderItem

		err := orderItemCollection.FindOne(ctx, bson.M{"orderitemid": orderItemId}).Decode(&orderItem)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing ordered item"})
			return
		}
		c.JSON(http.StatusOK, orderItem)
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var orderItem models.OrderItem

		orderItemId := c.Param("orderitemid")

		filter := bson.M{"orderitemid": orderItemId}

		var updateObj primitive.D

		if orderItem.UnitPrice != nil {
			updateObj = append(updateObj, bson.E{Key:"unitprice",Value: *&orderItem.UnitPrice})
		}

		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{Key:"quantity", Value:*orderItem.Quantity})
		}

		if orderItem.FoodId != nil {
			updateObj = append(updateObj, bson.E{Key:"foodid", Value:*orderItem.FoodId})
		}

		orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key:"updated_at",Value: orderItem.UpdatedAt})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderItemCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key:"$set", Value:updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := "Order item update failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var orderItemPack OrderItemPack
		var order models.Order

		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order.OrderDate, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		orderItemsToBeInserted := []interface{}{}
		order.TableId = orderItemPack.Table_id
		order_id := OrderItemOrderCreator(order)

		for _, orderItem := range orderItemPack.Order_items {
			orderItem.OrderId = order_id

			validationErr := validate.Struct(orderItem)

			if validationErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				return
			}
			orderItem.ID = primitive.NewObjectID()
			orderItem.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.OrderItemId = orderItem.ID.Hex()
			var num = toFixed(*orderItem.UnitPrice, 2)
			orderItem.UnitPrice = &num
			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}

		insertedOrderItems, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)

		if err != nil {
			log.Fatal(err)
		}
		defer cancel()

		c.JSON(http.StatusOK, insertedOrderItems)
	}
}
