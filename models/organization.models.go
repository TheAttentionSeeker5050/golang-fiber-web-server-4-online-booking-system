package models

import (
	"context"
	"errors"
	"example/web-server/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Organization struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	IDString         string             `json:"idString" bson:"idString"` // this just serves the purpose of converting the ObjectID to a string and its optional
	Name             string             `json:"name" bson:"name"`
	CreatedAt        time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time          `json:"updatedAt" bson:"updatedAt"`
	OwnerID          string             `json:"ownerID" bson:"ownerID"`
	Locations        *[]Location        `json:"locations" bson:"locations"`
	Reservations     *[]Reservation     `json:"reservations" bson:"reservations"`
	BookingResources *[]BookingResource `json:"bookingResources" bson:"bookingResources"`
}

// function to get the organization collection
func GetOrganizationCollection() (*mongo.Client, *mongo.Collection, error) {
	// Assuming you use a database connection named `db`
	// get the mongo client
	mongoClient := config.GetMongoClient()
	if mongoClient == nil {
		return nil, nil, errors.New("Mongo client is nil")
	}

	// now we will save the user to the database
	// get the database and collection
	clientCollection, err := config.GetMongoCollection(mongoClient, "organizations")
	if err != nil {
		return nil, nil, errors.New("Error getting collection")
	}

	return mongoClient, clientCollection, nil
}

// The crud operations for the organization model
// - CreateOrganization
func CreateOrganization(c *fiber.Ctx, organizationCollection *mongo.Collection, org *Organization) error {
	// only one organization can be created by a user with the same name
	// so we are checking if the organization already exists for the user
	err := organizationCollection.FindOne(context.TODO(), bson.M{"name": org.Name, "ownerID": org.OwnerID}).Err()
	if err == nil {
		return errors.New("Organization already exists")
	}

	// insert the organization into the database
	_, err = organizationCollection.InsertOne(context.TODO(), org)
	if err != nil {
		return errors.New("Error inserting organization")
	}

	return nil
}

// - GetOrganization
func GetOrganization(c *fiber.Ctx, organizationCollection *mongo.Collection, organizationID string) (*Organization, error) {
	// find the organization by id
	var org Organization
	// parse the id string to a primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(organizationID)
	if err != nil {
		return nil, errors.New("Error parsing object id")
	}
	err = organizationCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&org)
	if err != nil {
		return nil, errors.New("Organization not found")
	}

	// add the id string to the organization
	org.IDString = org.ID.Hex()

	return &org, nil
}

// - UpdateOrganization
func UpdateOrganization(c *fiber.Ctx, organizationCollection *mongo.Collection, org *Organization) error {

	// update the organization
	_, err := organizationCollection.UpdateOne(c.Context(), bson.M{"_id": org.ID}, bson.M{"$set": org})
	if err != nil {
		return errors.New("Error updating organization")
	}

	return nil
}

// - DeleteOrganization
func DeleteOrganization(c *fiber.Ctx, organizationCollection *mongo.Collection, org *Organization) error {
	// delete the organization
	_, err := organizationCollection.DeleteOne(context.TODO(), bson.M{"_id": org.ID})
	if err != nil {
		return errors.New("Error deleting organization")
	}

	return nil
}

// - GetOrganizations with filters, pagination, and sorting for the user
func GetOrganizations(c *fiber.Ctx, organizationCollection *mongo.Collection, ownerID string, filters map[string]string, offset int64, limit int64, sort string, sortOrder string) ([]Organization, error) {

	if filters == nil {
		filters = map[string]string{}
	}

	// create a filter for the ownerID
	filters["ownerID"] = ownerID

	// make the *options.FindOptions object
	var getOptions *options.FindOptions = options.Find()

	// if the page is nil, set it to 1 and reference it to the page argument

	getOptions.SetSkip(offset)
	getOptions.SetLimit(limit)

	// get the organizations, filter, sort, and paginate
	cursor, err := organizationCollection.Find(context.TODO(), filters, getOptions)
	if err != nil {
		return nil, errors.New("Error getting organizations")
	}

	// get the organizations
	var orgs []Organization

	// Iterate through the returned cursor and add the id string to the organization
	for cursor.Next(context.Background()) {
		var org Organization
		err := cursor.Decode(&org)
		if err != nil {
			return nil, errors.New("Error decoding organization")
		}
		org.IDString = org.ID.Hex()
		orgs = append(orgs, org)
	}

	cursor.Close(context.TODO())

	return orgs, nil
}
