package models

import (
	"errors"
	"example/web-server/config"
	"example/web-server/data"
	"example/web-server/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID           string `json:"id" bson:"_id"`
	Email        string `json:"email" bson:"email"`
	FirstName    string `json:"firstName" bson:"firstName"`
	LastName     string `json:"lastName" bson:"lastName"`
	Phone        string `json:"phone" bson:"phone"`
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

	// if claims.name is not empty, then split the name into first name and last name, if it is 1, dont fill, if the lenght is 2, then split it, if it is 3 then choose the first two as first name and the last one as last name, and if it is 4 then choose the first two as first name and the last two as last name. otherwise, dont fill
	firstName, lastName := utils.SplitNameStr(claims.Name)

	// insert the user into the database
	_, err = userCollection.InsertOne(c.Context(), User{
		ID:           claims.ID,
		Email:        claims.Email,
		PasswordHash: "",
		AuthProvider: data.AUTH_PROVIDER_GOOGLE,
		FirstName:    firstName,
		LastName:     lastName,
		Phone:        "",
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

// update the user data, only we can update the data: name, phone
func UpdateUser(c *fiber.Ctx, userCollection *mongo.Collection, user *User) error {
	var oldUserData User

	// first get the user by id
	err := userCollection.FindOne(c.Context(), bson.M{"_id": user.ID}).Decode(&oldUserData)
	if err != nil {
		return errors.New("User not found")
	}

	// update the user, only we can update the data: name, phone
	_, err = userCollection.UpdateOne(c.Context(), bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"firstName": user.FirstName, "lastName": user.LastName, "phone": user.Phone}})
	if err != nil {
		return errors.New("Error updating user")
	}

	return nil
}

func UpdateUserPassword(c *fiber.Ctx, userCollection *mongo.Collection, userID, passwordHash string) error {

	// if the provider is not Local
	var user User
	err := userCollection.FindOne(c.Context(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return errors.New("User not found")
	}

	if user.AuthProvider != data.AUTH_PROVIDER_LOCAL {
		return errors.New("Cannot update password because it is using an account provider")
	}

	// update the user password
	_, err = userCollection.UpdateOne(c.Context(), bson.M{"_id": userID}, bson.M{"$set": bson.M{"passwordHash": passwordHash}})
	if err != nil {
		return errors.New("Error updating user password")
	}

	return nil
}

func GetUserData(c *fiber.Ctx, userCollection *mongo.Collection, userID string) (*User, error) {
	var user User

	// get the user by id
	err := userCollection.FindOne(c.Context(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, errors.New("User not found")
	}

	// remove the password hash from the return value
	user.PasswordHash = ""

	return &user, nil
}

func GetUserPasswordHash(c *fiber.Ctx, userCollection *mongo.Collection, userID string) (string, error) {
	var user User

	// get the user by email
	err := userCollection.FindOne(c.Context(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return "", errors.New("User not found")
	}

	return user.PasswordHash, nil
}
