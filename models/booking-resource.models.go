package models

import "time"

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
	ID          string    `json:"id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	Type        string    `json:"type" bson:"type"`
	LocationID  string    `json:"locationId" bson:"locationId"`
	OwnerID     string    `json:"ownerId" bson:"ownerId"`
	Description string    `json:"description" bson:"description"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
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
func GetBookingResources(locationID string, filters map[string]interface{}, page int, limit int, sort string) ([]BookingResource, error) {
	return nil, nil
}
