package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Table struct {
	ID            primitive.ObjectID `bson:"_id"`
	NumberOfGuest *int64             `json:"numberofguest"`
	TableNumber   *int64             `json:"tablenumber"`
	CreatedAt     time.Time          `json:"createdat"`
	UpdatedAt     time.Time          `json:"updatedat"`
	Table_id      string             `json:"table_id"`
}
