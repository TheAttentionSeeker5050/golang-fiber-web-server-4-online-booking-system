package controllers

import (
	"example/web-server/config"
	"example/web-server/models"
	"example/web-server/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DetailViewOrganizationController(c *fiber.Ctx) error {

	// get the organization id from the url param :id
	organizationID := c.Params("id")

	// create the arguments map
	argumentsMap := &fiber.Map{
		"Title": "Organization Details",
	}

	// add the Authenticated flag to the argumentsMap using http headers
	(*argumentsMap)["IsAuthenticated"] = string(c.Request().Header.Peek("Authenticated")) == "true"

	// get the user id from the Locals
	userID := c.Locals("userID").(string)

	// get the mongo client and collection for organizations
	mongoClient, organizationCollection, err := models.GetOrganizationCollection()
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Error"] = "Error getting organization collection"
	} else {
		// get the organization from the model
		organization, err := models.GetOrganization(c, organizationCollection, organizationID)
		if err != nil {
			// add error msg to the argumentsMap
			(*argumentsMap)["Error"] = "Error getting Organization"
			(*argumentsMap)["Title"] = "Error - Resource Not Found"
		} else {
			// check if the user is the owner of the organization
			if organization.OwnerID != userID {
				// add error msg to the argumentsMap
				(*argumentsMap)["Error"] = "You are not the owner of the organization"
				(*argumentsMap)["Title"] = "Error - Unauthorized"
			} else {
				// make a filters map for all the locations to match the organization id
				filters := map[string]string{}

				// add the organization to the argumentsMap
				(*argumentsMap)["Organization"] = organization
				// add title to the argumentsMap to organization name
				(*argumentsMap)["Title"] = organization.Name + " - Details"

				// get the resources for the organization, the collections: locations, booking resources, and reservations
				locationCollection, err := config.GetMongoCollection(mongoClient, "locations")
				if err != nil {
					// add error msg to the argumentsMap
					(*argumentsMap)["Error"] = "Error getting location collection"
				} else {
					// add the organizationID to the filters
					filters["organizationID"] = organizationID

					// get the locations for the organization
					locations, err := models.GetLocations(c, locationCollection, userID, filters, 0, 5, "", "")
					if err != nil {
						// add error msg to the argumentsMap
						(*argumentsMap)["Error"] = "Error getting Locations"
					} else {
						// add the locations to the argumentsMap
						(*argumentsMap)["Locations"] = locations
					}
				}

				// get the booking resources collection
				bookingResourceCollection, err := config.GetMongoCollection(mongoClient, "booking_resources")
				if err != nil {
					// add error msg to the argumentsMap
					(*argumentsMap)["Error"] = "Error getting booking resource collection"
				} else {
					// get the booking resources for the organization
					bookingResources, err := models.GetBookingResources(c, bookingResourceCollection, userID, filters, 0, 5, "", "")
					if err != nil {
						// add error msg to the argumentsMap
						(*argumentsMap)["Error"] = "Error getting Booking Resources"
					} else {
						// add the booking resources to the argumentsMap
						(*argumentsMap)["BookingResources"] = bookingResources
					}
				}

				// get the reservations collection
				reservationCollection, err := config.GetMongoCollection(mongoClient, "reservations")
				if err != nil {
					// add error msg to the argumentsMap
					(*argumentsMap)["Error"] = "Error getting reservation collection"
				} else {
					// get the reservations for the organization
					reservations, err := models.GetReservations(c, reservationCollection, userID, filters, 0, 5, "", "")
					if err != nil {
						// add error msg to the argumentsMap
						(*argumentsMap)["Error"] = "Error getting Reservations"
					} else {
						// add the reservations to the argumentsMap
						(*argumentsMap)["Reservations"] = reservations
					}
				}
			}
		}
	}

	// close the mongo client connection
	config.CloseMongoClientConnection(mongoClient)

	return utils.CustomRenderTemplate(c, "organization/view", *argumentsMap)
}

func AddOrganizationController(c *fiber.Ctx) error {
	return utils.CustomRenderTemplate(c, "organization/add", fiber.Map{
		"Title": "Add Organization",
	})
}

func AddOrganizationPostController(c *fiber.Ctx) error {
	// return nil
	// Get the form data from api json
	var name string

	// convert the reqBody to a map
	formData := make(map[string]string)
	err := c.BodyParser(&formData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing form data",
		})
	}

	name = formData["name"]

	// get the owner id from the Local
	ownerID := c.Locals("userID").(string)

	// if name or ownerID is empty
	if name == "" || ownerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Name or ownerID is empty",
		})
	}

	// get the mongo collection and client for organizations
	mongoClient, organizationCollection, err := models.GetOrganizationCollection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting organization collection",
		})
	}

	// run the model function create organization
	err = models.CreateOrganization(c, organizationCollection, &models.Organization{
		ID:        primitive.NewObjectID(),
		Name:      name,
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// close the mongo client connection
	config.CloseMongoClientConnection(mongoClient)

	// Add organization with json http code 200
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Organization added successfully",
	})
}

func EditOrganizationController(c *fiber.Ctx) error {

	// get edit slug parameter id
	organizationID := c.Params("id")

	// make the arguments map structure
	argumentsMap := &fiber.Map{
		"Title": "Edit Organization",
	}

	// add the Authenticated flag to the argumentsMap using http headers
	(*argumentsMap)["IsAuthenticated"] = string(c.Request().Header.Peek("Authenticated")) == "true"

	// get the user id from the Locals
	userID := c.Locals("userID").(string)

	// get the mongo client and collection for organizations
	mongoClient, organizationCollection, err := models.GetOrganizationCollection()
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Error"] = "Error getting organization collection"
	} else {
		// get the organization from the model
		organization, err := models.GetOrganization(c, organizationCollection, organizationID)
		if err != nil {
			// add error msg to the argumentsMap
			(*argumentsMap)["Error"] = "Error getting Organization"
			(*argumentsMap)["Title"] = "Error - Resource Not Found"
		} else {
			// check if the user is the owner of the organization
			if organization.OwnerID != userID {
				// add error msg to the argumentsMap
				(*argumentsMap)["Error"] = "You are not the owner of the organization"
				(*argumentsMap)["Title"] = "Error - Unauthorized"
			} else {
				// add the organization to the argumentsMap
				(*argumentsMap)["Organization"] = organization
			}
		}
	}

	// close the mongo client
	config.CloseMongoClientConnection(mongoClient)

	return utils.CustomRenderTemplate(c, "organization/edit", *argumentsMap)
}

func EditOrganizationPostController(c *fiber.Ctx) error {

	// return nil
	// Get the form data from api json
	var newName string

	// get edit slug parameter id
	organizationID := c.Params("id")
	if organizationID == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Could not find the organization",
		})
	}

	// convert the reqBody to a map
	formData := make(map[string]string)
	err := c.BodyParser(&formData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid form data",
		})
	}

	newName = formData["name"]
	if newName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid new name",
		})
	}

	// get the owner id from the Local
	ownerID := c.Locals("userID").(string)

	// get the mongo collection and client for organizations
	mongoClient, organizationCollection, err := models.GetOrganizationCollection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting organization collection",
		})
	}

	// get the organization from the model
	organization, err := models.GetOrganization(c, organizationCollection, organizationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting organization",
		})
	}

	// if ownerID is not the owner of the organization
	if organization.OwnerID != ownerID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not the owner of the organization",
		})
	}

	// add the new name to the organization
	organization.Name = newName

	// run the model function update organization
	err = models.UpdateOrganization(c, organizationCollection, organization)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating organization",
		})
	}

	// close the client connection
	config.CloseMongoClientConnection(mongoClient)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Organization updated successfully",
	})
}

func DeleteOrganizationController(c *fiber.Ctx) error {

	// get edit slug parameter id
	organizationID := c.Params("id")

	// make the arguments map structure
	argumentsMap := &fiber.Map{
		"Title": "Edit Organization",
	}

	// add the Authenticated flag to the argumentsMap using http headers
	(*argumentsMap)["IsAuthenticated"] = string(c.Request().Header.Peek("Authenticated")) == "true"

	// get the user id from the Locals
	userID := c.Locals("userID").(string)

	// get the mongo client and collection for organizations
	mongoClient, organizationCollection, err := models.GetOrganizationCollection()
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Error"] = "Error getting organization collection"
	} else {
		// get the organization from the model
		organization, err := models.GetOrganization(c, organizationCollection, organizationID)
		if err != nil {
			// add error msg to the argumentsMap
			(*argumentsMap)["Error"] = "Error getting Organization"
			(*argumentsMap)["Title"] = "Error - Resource Not Found"
		} else {
			// check if the user is the owner of the organization
			if organization.OwnerID != userID {
				// add error msg to the argumentsMap
				(*argumentsMap)["Error"] = "You are not the owner of the organization"
				(*argumentsMap)["Title"] = "Error - Unauthorized"
			} else {
				// add the organization to the argumentsMap
				(*argumentsMap)["Organization"] = organization
			}
		}
	}

	// close the mongo client
	config.CloseMongoClientConnection(mongoClient)

	return utils.CustomRenderTemplate(c, "organization/delete", *argumentsMap)
}

func DeleteOrganizationPostController(c *fiber.Ctx) error {

	// get edit slug parameter id
	organizationID := c.Params("id")
	if organizationID == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Could not find the organization",
		})
	}

	// get the owner id from the Local
	ownerID := c.Locals("userID").(string)

	// get the mongo collection and client for organizations
	mongoClient, organizationCollection, err := models.GetOrganizationCollection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting organization collection",
		})
	}

	// get the organization from the model
	organization, err := models.GetOrganization(c, organizationCollection, organizationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting organization",
		})
	}

	// if ownerID is not the owner of the organization
	if organization.OwnerID != ownerID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not the owner of the organization",
		})
	}

	// delete the organization
	err = models.DeleteOrganization(c, organizationCollection, organization)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting organization",
		})
	}

	// close the client connection
	config.CloseMongoClientConnection(mongoClient)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Organization updated successfully",
	})
}

func ListViewOrganizationsController(c *fiber.Ctx) error {

	argumentsMap := &fiber.Map{
		"Title": "Organizations",
	}

	// add the Authenticated flag to the argumentsMap using http headers
	(*argumentsMap)["IsAuthenticated"] = string(c.Request().Header.Peek("Authenticated")) == "true"

	userID := c.Locals("userID").(string)

	// start the mongo client, the collection and the error
	mongoClient, organizationCollection, err := models.GetOrganizationCollection()
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Error"] = "Error getting organization collection"
	} else {
		organizations, err := models.GetOrganizations(c, organizationCollection, userID, nil, 0, 5, "", "")
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

	return utils.CustomRenderTemplate(c, "organization/index", *argumentsMap)
}
