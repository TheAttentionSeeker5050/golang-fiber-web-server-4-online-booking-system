package models

import (
	"context"
	"errors"
	"example/web-server/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BookingResource is a struct that represents the booking resource of an organization's location. In general terms, it can be a room, a vehicle, a piece of equipment, a table, time with a professional, etc.
// mongodb collection name: booking_resources
// The properties of these are
// - ID: string // this is the primary key
// - Name: string // this is the name of the booking resource
// - Type: string // this is the type of the booking resource, how we call it, for example, room, vehicle, equipment, table, etc. Default is table
// - LocationID: string // this is the location id of the booking resource
// - OwnerID: string // this is the owner id of the booking resource
// - Description: string // this is the description of the booking resource
// - CreatedAt: time.Time
// - UpdatedAt: time.Time
type BookingResource struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Name           string             `json:"name" bson:"name"`
	Type           string             `json:"type" bson:"type"`
	LocationID     string             `json:"locationID" bson:"locationID"`
	Location       *Location          `json:"location" bson:"location"`
	Reservations   []Reservation      `json:"reservations" bson:"reservations"`
	OwnerID        string             `json:"ownerID" bson:"ownerID"`
	OrganizationID string             `json:"organizationID" bson:"organizationID"`
	Organization   *Organization      `json:"organization" bson:"organization"`
	Description    string             `json:"description" bson:"description"`
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func GetBookingResourceCollection() (*mongo.Client, *mongo.Collection, error) {
	mongoClient := config.GetMongoClient()
	if mongoClient == nil {
		return nil, nil, errors.New("Mongo client is nil")
	}

	// get the collection
	clientCollection, err := config.GetMongoCollection(mongoClient, "booking_resources")
	if err != nil {
		return nil, nil, errors.New("Error getting collection")
	}

	return mongoClient, clientCollection, nil
}

// The crud operations for the booking resource model
// - CreateBookingResource
func CreateBookingResource(br *BookingResource) error {
	return nil
}

// - GetBookingResource
func GetBookingResource(id string) (*BookingResource, error) {
	return nil, nil
}

// - UpdateBookingResource
func UpdateBookingResource(br *BookingResource) error {
	return nil
}

// - DeleteBookingResource
func DeleteBookingResource(id string) error {
	return nil
}

// - GetBookingResources with filters, pagination, and sorting for the user
func GetBookingResources(c *fiber.Ctx, bookingResourceCollection *mongo.Collection, ownerID string, filters map[string]string, offset int64, limit int64, sort string, sortOrder string) ([]BookingResource, error) {
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

	// get the booking resources, filter, sort, and paginate
	cursor, err := bookingResourceCollection.Find(context.TODO(), filters, getOptions)
	if err != nil {
		return nil, errors.New("Error getting organizations")
	}

	// get the organizations
	var bookingResources []BookingResource

	err = cursor.All(context.TODO(), &bookingResources)
	if err != nil {
		return nil, errors.New("Error getting organizations")
	}

	cursor.Close(context.TODO())

	return bookingResources, nil
}
