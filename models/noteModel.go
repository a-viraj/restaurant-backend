package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID        primitive.ObjectID `bson:"_id"`
	Text      string             `json:"text"`
	Tittle    string             `json:"tittle"`
	CreatedAt time.Time          `json:"createdat"`
	UpdatedAt time.Time          `json:"updatedat"`
	NoteId    string             `json:"noteid"`
}
