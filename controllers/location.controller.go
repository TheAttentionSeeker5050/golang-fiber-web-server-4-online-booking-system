package controllers

import (
	"example/web-server/config"
	"example/web-server/models"
	"example/web-server/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Make controller functions for the location routes
func DetailViewLocationController(c *fiber.Ctx) error {

	// get the location id from the url param :id
	locationID := c.Params("id")
	// print the location id

	// create the arguments map
	argumentsMap := &fiber.Map{
		"Title": "location Details",
	}

	// add the Authenticated flag to the argumentsMap using http headers
	(*argumentsMap)["IsAuthenticated"] = string(c.Request().Header.Peek("Authenticated")) == "true"

	// get the user id from the Locals
	userID := c.Locals("userID").(string)

	// get the mongo client and collection for locations
	mongoClient, locationCollection, err := models.GetLocationCollection()
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Error"] = "Error getting location collection"
	} else {
		// get the location from the model
		location, err := models.GetLocation(c, locationCollection, locationID)
		if err != nil {
			// add error msg to the argumentsMap
			(*argumentsMap)["Error"] = "Error getting Location"
			(*argumentsMap)["Title"] = "Error - Resource Not Found"
		} else {
			// check if the user is the owner of the location
			if location.OwnerID != userID {
				// add error msg to the argumentsMap
				(*argumentsMap)["Error"] = "You are not the owner of the location"
				(*argumentsMap)["Title"] = "Error - Unauthorized"
			} else {
				// add the location to the argumentsMap
				(*argumentsMap)["Location"] = location
				// add title to the argumentsMap to location name
				(*argumentsMap)["Title"] = location.Name + " - Details"
			}

		}
	}

	// close the mongo client connection
	config.CloseMongoClientConnection(mongoClient)

	// For now return the page
	return utils.CustomRenderTemplate(c, "location/view", *argumentsMap)
}

func AddLocationController(c *fiber.Ctx) error {

	// create the arguments map
	argumentsMap := &fiber.Map{
		"Title": "Add Location",
	}

	// add the Authenticated flag to the argumentsMap using http headers
	(*argumentsMap)["IsAuthenticated"] = string(c.Request().Header.Peek("Authenticated")) == "true"

	// get the user id from the Locals
	userID := c.Locals("userID").(string)

	// Add user id to the argumentsMap
	(*argumentsMap)["UserID"] = userID

	// get the organizations collection
	mongoClient, organizationCollection, err := models.GetOrganizationCollection()

	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Error"] = "Error getting organization collection"
	} else {
		// get the organizations for the user
		organizations, err := models.GetOrganizations(c, organizationCollection, userID, nil, 0, 1000, "", "")

		if err != nil {
			// add error msg to the argumentsMap
			(*argumentsMap)["Error"] = "Error getting Organizations"
		} else {
			// add the organizations to the argumentsMap
			(*argumentsMap)["Organizations"] = organizations
		}
	}

	// close the mongo client connection
	config.CloseMongoClientConnection(mongoClient)

	return utils.CustomRenderTemplate(c, "location/add", *argumentsMap)
}

func AddLocationPostController(c *fiber.Ctx) error {

	// get the user id from the Locals
	userID := c.Locals("userID").(string)

	// Get the form values
	name := c.FormValue("name")
	organizationID := c.FormValue("organizationID")
	address := c.FormValue("address")
	city := c.FormValue("city")
	state := c.FormValue("state")
	zip := c.FormValue("zip")
	country := c.FormValue("country")
	phone := c.FormValue("phone")
	email := c.FormValue("email")

	// get the mongo client and collection for organizations
	mongoClient, organizationCollection, err := models.GetOrganizationCollection()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	} else {
		// get the organization from the model
		organization, err := models.GetOrganization(c, organizationCollection, organizationID)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		} else {
			// check if the user is the owner of the organization
			if organization.OwnerID != userID {
				return c.SendStatus(fiber.StatusUnauthorized)
			}
		}
	}

	// Now create a new location instance, we use the & operator to create a pointer to the location struct
	location := &models.Location{
		ID:             primitive.NewObjectID(),
		Name:           name,
		OrganizationID: organizationID,
		Address:        address,
		City:           city,
		State:          state,
		Zip:            zip,
		Country:        country,
		Phone:          phone,
		Email:          email,
		OwnerID:        userID,
	}

	// get the mongo client and collection for locations
	mongoClient, locationCollection, err := models.GetLocationCollection()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	} else {
		// add the location to the model
		err = models.CreateLocation(c, locationCollection, location)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	// close the mongo client connection
	config.CloseMongoClientConnection(mongoClient)

	// For now return a redirect to the location details page
	return c.Redirect("/locations")
}

func EditLocationController(c *fiber.Ctx) error {

	// create the arguments map
	argumentsMap := &fiber.Map{
		"Title": "Edit Location",
	}

	// For now return the page
	return utils.CustomRenderTemplate(c, "location/edit", *argumentsMap)
}

func EditLocationPostController(c *fiber.Ctx) error {

	// For now return a redirect to the location details page
	return c.Redirect("/locations")
}

func DeleteLocationController(c *fiber.Ctx) error {
	// create the arguments map
	argumentsMap := &fiber.Map{
		"Title": "Delete Location",
	}

	// For now return the page
	return utils.CustomRenderTemplate(c, "location/delete", *argumentsMap)
}

func DeleteLocationPostController(c *fiber.Ctx) error {
	// For now return a redirect to the Location details page
	return c.Redirect("/locations")
}

func BulkDeleteLocationController(c *fiber.Ctx) error {

	// For now return a redirect to the Location details page
	return c.Redirect("/locations")
}

func ListViewLocationsController(c *fiber.Ctx) error {

	argumentsMap := &fiber.Map{
		"Title": "Locations",
	}

	// add the Authenticated flag to the argumentsMap using http headers
	(*argumentsMap)["IsAuthenticated"] = string(c.Request().Header.Peek("Authenticated")) == "true"

	// get the user id from the Locals
	userID := c.Locals("userID").(string)

	// start the mongo client, the collection and the error
	mongoClient, locationCollection, err := models.GetLocationCollection()

	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Error"] = "Error getting location collection"
		return utils.CustomRenderTemplate(c, "locations/index", *argumentsMap)
	} else {
		locations, err := models.GetLocations(c, locationCollection, userID, nil, 0, 5, "", "")

		if err != nil {
			// add error msg to the argumentsMap
			(*argumentsMap)["Error"] = "Error getting Locations"
			return utils.CustomRenderTemplate(c, "locations/index", *argumentsMap)
		} else {
			// add the locations to the argumentsMap
			(*argumentsMap)["Locations"] = locations
		}
	}

	// close the mongo client connection
	config.CloseMongoClientConnection(mongoClient)

	// For now return the page
	return utils.CustomRenderTemplate(c, "location/index", *argumentsMap)
}
