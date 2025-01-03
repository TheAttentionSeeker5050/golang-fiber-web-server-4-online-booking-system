package models

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
	ID                 string `json:"id" bson:"_id"`
	BookingResourceID  string `json:"bookingResourceId" bson:"bookingResourceId"`
	AssistantEmail     string `json:"assistantEmail" bson:"assistantEmail"`
	AssistantPhone     string `json:"assistantPhone" bson:"assistantPhone"`
	AssistantFirstName string `json:"assistantFirstName" bson:"assistantFirstName"`
	AssistantLastName  string `json:"assistantLastName" bson:"assistantLastName"`
	StartDate          string `json:"startDate" bson:"startDate"`
	EndDate            string `json:"endDate" bson:"endDate"`
	CreatedAt          string `json:"createdAt" bson:"createdAt"`
	UpdatedAt          string `json:"updatedAt" bson:"updatedAt"`
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
func GetReservations(bookingResourceID string, filters map[string]interface{}, page int, limit int, sort string) ([]Reservation, error) {
	return nil, nil
}
