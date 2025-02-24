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

type Location struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	IDString         string             `json:"idString" bson:"idString"`
	Name             string             `json:"name" bson:"name"`
	Address          string             `json:"address" bson:"address"`
	City             string             `json:"city" bson:"city"`
	State            string             `json:"state" bson:"state"`
	Zip              string             `json:"zip" bson:"zip"`
	Country          string             `json:"country" bson:"country"`
	CreatedAt        time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time          `json:"updatedAt" bson:"updatedAt"`
	Email            string             `json:"email" bson:"email"`
	Phone            string             `json:"phone" bson:"phone"`
	OwnerID          string             `json:"ownerID" bson:"ownerID"`
	OrganizationID   string             `json:"organizationID" bson:"organizationID"`
	Organization     *Organization      `json:"organization" bson:"organization"`
	BookingResources []BookingResource  `json:"bookingResources" bson:"bookingResources"`
}

// function to get the location collection
func GetLocationCollection() (*mongo.Client, *mongo.Collection, error) {
	// Assuming you use a database connection named `db`
	// get the mongo client
	mongoClient := config.GetMongoClient()
	if mongoClient == nil {
		return nil, nil, errors.New("Mongo client is nil")
	}

	// now we will save the user to the database
	// get the database and collection
	clientCollection, err := config.GetMongoCollection(mongoClient, "locations")
	if err != nil {
		return nil, nil, errors.New("Error getting collection")
	}

	return mongoClient, clientCollection, nil
}

// The crud operations for the location model
// - CreateLocation
func CreateLocation(c *fiber.Ctx, organizationCollection *mongo.Collection, loc *Location) error { // we use the * operator to get the value of the pointer from the function argument
	// only one location can be created by a user with the same name
	// so we are checking if the location already exists for the user
	err := organizationCollection.FindOne(c.Context(), bson.M{"name": loc.Name, "ownerID": loc.OwnerID}).Err()
	if err == nil {
		return errors.New("Location already exists")
	}

	// insert the location into the database
	_, err = organizationCollection.InsertOne(c.Context(), loc)
	if err != nil {
		return errors.New("Error inserting location")
	}

	return nil
}

// - GetLocation
func GetLocation(c *fiber.Ctx, locationCollection *mongo.Collection, id string) (*Location, error) {
	var loc Location

	// parse the id string to a primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("Error parsing object id")
	}

	err = locationCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&loc)
	if err != nil {
		return nil, errors.New("Error getting location")
	}

	// add id string to the location
	loc.IDString = loc.ID.Hex()

	return &loc, nil
}

// - UpdateLocation
func UpdateLocation(c *fiber.Ctx, organizationCollection *mongo.Collection, loc *Location) error {
	// update the location in the database
	_, err := organizationCollection.UpdateOne(c.Context(), bson.M{"_id": loc.ID}, bson.M{"$set": loc})
	if err != nil {
		return errors.New("Error updating location")
	}

	return nil
}

// - DeleteLocation
func DeleteLocation(c *fiber.Ctx, organizationCollection *mongo.Collection, id string) error {
	// delete the location from the database
	_, err := organizationCollection.DeleteOne(c.Context(), bson.M{"_id": id})
	if err != nil {
		return errors.New("Error deleting location")
	}

	return nil
}

// - GetLocations with filters, pagination, and sorting for the user
func GetLocations(c *fiber.Ctx, locationCollection *mongo.Collection, ownerID string, filters map[string]string, offset int64, limit int64, sort string, sortOrder string) ([]Location, error) {
	if filters == nil {
		filters = map[string]string{}
	}

	// add owner id to the filters
	filters["ownerID"] = ownerID

	// make the *options.FindOptions object
	var getOptions *options.FindOptions = options.Find()

	// if the page is nil, set it to 1 and reference it to the page argument
	getOptions.SetSkip(offset)
	getOptions.SetLimit(limit)

	// find the locations, filter, sort, and paginate
	cursor, err := locationCollection.Find(context.TODO(), filters, getOptions)
	if err != nil {
		return nil, errors.New("Error getting locations")
	}

	// get the locations from the database
	var locations []Location

	// Iterate to add idString to the locations
	for cursor.Next(context.Background()) {
		var loc Location
		err := cursor.Decode(&loc)
		if err != nil {
			return nil, errors.New("Error decoding location")
		}

		// add id string to the location
		loc.IDString = loc.ID.Hex()
		locations = append(locations, loc)

	}

	cursor.Close(context.TODO())

	return locations, nil
}
