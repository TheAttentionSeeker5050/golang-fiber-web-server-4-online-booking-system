package controllers

import (
	"example/web-server/config"
	"example/web-server/data"
	"example/web-server/models"
	"example/web-server/utils"

	"github.com/gofiber/fiber/v2"
)

func ViewDashboardController(c *fiber.Ctx) error {
	argumentsMap := &fiber.Map{
		"Title": "Dashboard",
	}

	// add the Authenticated flag to the argumentsMap using http headers
	(*argumentsMap)["IsAuthenticated"] = string(c.Request().Header.Peek("Authenticated")) == "true"

	userID := c.Locals("userID").(string)

	// start the mongo client
	mongoClient := config.GetMongoClient()
	if mongoClient == nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Error"] = "Error starting mongo client"
		return utils.CustomRenderTemplate(c, "dashboard/index", *argumentsMap)
	}

	// get the Organizations collection
	organizationsCollection, err := config.GetMongoCollection(mongoClient, "organizations")
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Error"] = "Error getting Organizations collection"
		return utils.CustomRenderTemplate(c, "dashboard/index", *argumentsMap)
	}

	// add a number of result limit, which will be used for each of the collections
	// filter the collection by the owner id
	organizations, err := models.GetOrganizations(c, organizationsCollection, userID, nil, 0, 5, "", "")
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Error"] = "Error getting Organizations"
		return utils.CustomRenderTemplate(c, "dashboard/index", *argumentsMap)
	}

	// make dummy data just to display in the dashboard page
	(*argumentsMap)["Organizations"] = organizations

	// get the locations collection
	locationsCollection, err := config.GetMongoCollection(mongoClient, "locations")
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Error"] = "Error getting Locations collection"
		return utils.CustomRenderTemplate(c, "dashboard/index", *argumentsMap)
	}

	// get the locations using models function
	locations, err := models.GetLocations(c, locationsCollection, userID, nil, 0, 5, "", "")
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Error"] = "Error getting Locations"
		return utils.CustomRenderTemplate(c, "dashboard/index", *argumentsMap)
	}

	(*argumentsMap)["Locations"] = locations

	// get the booking resources collection
	bookingResourcesCollection, err := config.GetMongoCollection(mongoClient, "booking_resources")
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Error"] = "Error getting Booking Resources collection"
		return utils.CustomRenderTemplate(c, "dashboard/index", *argumentsMap)
	}

	// get the booking resources using models function
	bookingResources, err := models.GetBookingResources(c, bookingResourcesCollection, userID, nil, 0, 5, "", "")

	// add the booking resources to the argumentsMap
	(*argumentsMap)["BookingResources"] = bookingResources

	// get the reservations collection
	reservationsCollection, err := config.GetMongoCollection(mongoClient, "reservations")
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Error"] = "Error getting Reservations collection"
		return utils.CustomRenderTemplate(c, "dashboard/index", *argumentsMap)
	}

	// get the reservations using models function
	reservations, err := models.GetReservations(c, reservationsCollection, userID, nil, 0, 5, "", data.SORT_ORDER_DESC)
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Error"] = "Error getting Reservations"
		return utils.CustomRenderTemplate(c, "dashboard/index", *argumentsMap)
	}

	(*argumentsMap)["Reservations"] = reservations

	// close the mongo client connection
	config.CloseMongoClientConnection(mongoClient)

	return utils.CustomRenderTemplate(c, "dashboard/index", *argumentsMap)
}
