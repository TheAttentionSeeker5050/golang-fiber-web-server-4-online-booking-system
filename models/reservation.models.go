package models

import (
	"context"
	"errors"
	"example/web-server/config"
	"example/web-server/data"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Reservation is a struct that represents the reservation of a booking resource
// mongodb collection name: reservations
// The properties of these are
// - ID: string // this is the primary key
// - BookingResourceID: string // this is the booking resource id of the reservation
// - AssistantEmail: string // this is the assistant email of the reservation
// - AssistantPhone: string // this is the assistant phone of the reservation
// - AssistantFirstName: string // this is the assistant first name of the reservation
// - AssistantLastName: string // this is the assistant last name of the reservation
// - StartDate: time.Time // this is the start date of the reservation
// - EndDate: time.Time // this is the end date of the reservation
// - CreatedAt: time.Time
// - UpdatedAt: time.Time
type Reservation struct {
	ID                 string           `json:"id" bson:"_id"`
	OwnerID            string           `json:"ownerID" bson:"ownerID"`
	BookingResourceID  string           `json:"bookingResourceID" bson:"bookingResourceID"`
	BookingResource    *BookingResource `json:"bookingResource" bson:"bookingResource"`
	AssistantEmail     string           `json:"assistantEmail" bson:"assistantEmail"`
	AssistantPhone     string           `json:"assistantPhone" bson:"assistantPhone"`
	AssistantFirstName string           `json:"assistantFirstName" bson:"assistantFirstName"`
	AssistantLastName  string           `json:"assistantLastName" bson:"assistantLastName"`
	StartDate          time.Time        `json:"startDate" bson:"startDate"`
	EndDate            time.Time        `json:"endDate" bson:"endDate"`
	CreatedAt          time.Time        `json:"createdAt" bson:"createdAt"`
	UpdatedAt          time.Time        `json:"updatedAt" bson:"updatedAt"`
}

func GetReservationCollection() (*mongo.Client, *mongo.Collection, error) {
	// Assuming you use a database connection named `db`
	// get the mongo client
	mongoClient := config.GetMongoClient()
	if mongoClient == nil {
		return nil, nil, errors.New("Mongo client is nil")
	}

	// now we will save the reservation to the database
	// get the database and collection
	clientCollection, err := config.GetMongoCollection(mongoClient, "reservations")
	if err != nil {
		return nil, nil, errors.New("Error getting collection")
	}

	return mongoClient, clientCollection, nil
}

// The crud operations for the reservation model
// - CreateReservation
func CreateReservation(res *Reservation) error {
	return nil
}

// - GetReservation
func GetReservation(id string) (*Reservation, error) {
	return nil, nil
}

// - UpdateReservation
func UpdateReservation(res *Reservation) error {
	return nil
}

// - DeleteReservation
func DeleteReservation(id string) error {
	return nil
}

// - GetReservations with filters, pagination, and sorting for the user
func GetReservations(c *fiber.Ctx, reservationCollection *mongo.Collection, ownerID string, filters map[string]string, offset int64, limit int64, sortStr string, sortOrder int) ([]Reservation, error) {
	if filters == nil {
		filters = map[string]string{}
	}

	// create a filter for the ownerID
	filters["ownerID"] = ownerID

	// make the *options.FindOptions object
	var getOptions *options.FindOptions = options.Find()

	// make sort order from ascending or descending or asc or desc to either 1 or -1
	if sortOrder != data.SORT_ORDER_ASC && sortOrder != data.SORT_ORDER_DESC {
		sortOrder = data.SORT_ORDER_DESC
	}

	if sortStr == "" {
		sortStr = "startDate"
	}

	// if the page is nil, set it to 1 and reference it to the page argument
	getOptions.SetSkip(offset)
	getOptions.SetLimit(limit)
	getOptions.SetSort(bson.D{{sortStr, sortOrder}})

	// get the reservations, filter, sort, and paginate
	cursor, err := reservationCollection.Find(context.TODO(), filters, getOptions)
	if err != nil {
		return nil, errors.New("Error getting reservations")
	}

	// get the reservations from the database
	var reservations []Reservation
	err = cursor.All(context.TODO(), &reservations)
	if err != nil {
		return nil, errors.New("Error getting reservations")
	}

	cursor.Close(context.TODO())

	return reservations, nil
}
