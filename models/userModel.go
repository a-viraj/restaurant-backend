package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    *string            `json:"firstname" validate:"required"`
	LastName     *string            `json:"lastname"`
	Password     *string            `json:"password"`
	Email        *string            `json:"email"`
	Avatar       *string            `json:"avatar"`
	Phone        *string            `json:"phone"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refreshtoken"`
	CreatedAt    time.Time          `json:"createdat"`
	UpdatedAt    time.Time          `json:"updatedat"`
	UserType     *string            `json:"usertype"`
	UserId       string             `json:"userid"`
}
