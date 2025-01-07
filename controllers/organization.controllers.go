package controllers

import (
	"example/web-server/config"
	"example/web-server/models"
	"example/web-server/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func DetailViewOrganizationController(c *fiber.Ctx) error {

	return nil
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
		ID:        uuid.New().String(),
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

	return nil
}

func EditOrganizationPostController(c *fiber.Ctx) error {

	return nil
}

func DeleteOrganizationController(c *fiber.Ctx) error {

	return nil
}

func DeleteOrganizationPostController(c *fiber.Ctx) error {

	return nil
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
