package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID             primitive.ObjectID `bson:"_id"`
	InvoiceId      string             `json:"invoiceid"`
	OrderId        string             `json:"orderid"`
	PaymentMethod  *string             `json:"paymentmethod" validate:"eq=CASH|eq=CARD"`
	PaymentStatus  *string             `json:"paymentstatus" validate:"required"`
	PaymentDueDate time.Time          `json:"paymentduedate"`
	CreatedAt      time.Time          `json:"createdat"`
	UpdatedAt      time.Time          `json:"updatedat"`
}
