package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aviraj/resturant-management/database"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	jwt.StandardClaims
}

var userCollection mongo.Collection = *database.OpenCollection(database.Client, "user")

var SECRETKEY = os.Getenv("secretKey")

func GenerateAllTokens(email string, firstName string, lastName string, uid string) (string, string, error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Uid:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRETKEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRETKEY))
	if err != nil {
		log.Panic(err)
	}
	return token, refreshToken, err

}

func UpdateAllToken(token string, refreshToken string, userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: token})
	updateObj = append(updateObj, bson.E{Key: "refreshToken", Value: refreshToken})
	UpdatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: UpdatedAt})

	upsert := true
	filter := bson.M{"userId": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)
	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}

}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRETKEY), nil
		},
	)

	//the token is invalid

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	//the token is expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprint("token is expired")
		msg = err.Error()
		return
	}

	return claims, msg

}