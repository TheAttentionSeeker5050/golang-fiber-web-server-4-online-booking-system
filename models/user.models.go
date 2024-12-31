package models

import (
	"errors"
	"example/web-server/config"
	"example/web-server/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID           string `json:"id" bson:"_id"`
	Email        string `json:"email" bson:"email"`
	PasswordHash string `json:"passwordHash" bson:"passwordHash"`
	AuthProvider string `json:"authProvider" bson:"authProvider"`
	Picture      string `json:"picture"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func GetUserCollection() (*mongo.Client, *mongo.Collection, error) {
	// Assuming you use a database connection named `db`
	// get the mongo client
	mongoClient := config.GetMongoClient()
	if mongoClient == nil {
		return nil, nil, errors.New("Mongo client is nil")
	}

	// now we will save the user to the database
	// get the database and collection
	clientCollection, err := config.GetMongoCollection(mongoClient, "users")
	if err != nil {
		return nil, nil, errors.New("Error getting collection")
	}

	return mongoClient, clientCollection, nil
}

func SaveUserToDBUsingGoogleProvider(c *fiber.Ctx, userCollection *mongo.Collection, claims *utils.GoogleClaims) error {
	// check if the user already exists with the same email
	err := userCollection.FindOne(c.Context(), claims.Email).Err()
	if err == nil {
		return errors.New("User already exists")
	}

	// insert the user into the database
	_, err = userCollection.InsertOne(c.Context(), User{
		ID:           claims.ID,
		Email:        claims.Email,
		PasswordHash: "",
		AuthProvider: "Google",
		Picture:      claims.Picture,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		return errors.New("Error inserting user")
	}

	return nil
}

func SaveUserToDBUsingLocalAuthProvider(c *fiber.Ctx, userCollection *mongo.Collection, user User) error {
	// check if the user already exists with the same email
	err := userCollection.FindOne(c.Context(), user.Email).Err()
	if err == nil {
		return errors.New("User already exists")
	}

	// insert the user into the database
	_, err = userCollection.InsertOne(c.Context(), user)
	if err != nil {
		return errors.New("Error inserting user")
	}

	return nil
}
