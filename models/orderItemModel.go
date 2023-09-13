package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID          primitive.ObjectID `bson:"_id"`
	Quantity    *string            `json:"quantity" validate:"required,eq=S|eq=M|eq=L"`
	UnitPrice   *float64           `json:"unitprice" validate:"required"`
	CreatedAt   time.Time          `json:"createdat"`
	UpdatedAt   time.Time          `json:"updatedat"`
	FoodId      *string            `json:"foodid" validate:"required"`
	OrderItemId string             `json:"orderitemid"`
	OrderId     string             `json:"orderid" validate:"required"`
}
