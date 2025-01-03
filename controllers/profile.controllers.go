package controllers

import (
	"example/web-server/config"
	"example/web-server/models"
	"example/web-server/utils"

	"github.com/gofiber/fiber/v2"
)

// ViewProfileDataController
func ViewProfileDataController(c *fiber.Ctx) error {
	argumentsMap := &fiber.Map{
		"Title": "Profile",
	}

	// add the Authenticated flag to the argumentsMap using http headers
	(*argumentsMap)["IsAuthenticated"] = string(c.Request().Header.Peek("Authenticated")) == "true"

	// get the user id from the local context
	userID := c.Locals("userID").(string)

	// get the user data using the user id
	mongoClient, userMongoCollection, err := models.GetUserCollection()
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Error"] = "Error getting user data"
		return utils.CustomRenderTemplate(c, "auth/profile", *argumentsMap)
	}

	// get the user data
	var user *models.User
	user, err = models.GetUserData(c, userMongoCollection, userID)

	// close the mongo client connection
	config.CloseMongoClientConnection(mongoClient)

	// if the user is not found
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Error"] = "User not found"
	}

	// add the user data to the argumentsMap
	(*argumentsMap)["User"] = user

	return utils.CustomRenderTemplate(c, "profile/index", *argumentsMap)
}

func EditProfileController(c *fiber.Ctx) error {
	argumentsMap := &fiber.Map{
		"Title": "Edit Profile",
	}

	// add the Authenticated flag to the argumentsMap using http headers
	(*argumentsMap)["IsAuthenticated"] = string(c.Request().Header.Peek("Authenticated")) == "true"

	// get the old user data
	userID := c.Locals("userID").(string)

	// get the user data using the user id
	mongoClient, userMongoCollection, err := models.GetUserCollection()
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Errors"] = []string{"Error getting user data"}
		return utils.CustomRenderTemplate(c, "auth/profile", *argumentsMap)
	}

	// get the user data
	var user *models.User
	user, err = models.GetUserData(c, userMongoCollection, userID)

	// close the mongo client connection
	config.CloseMongoClientConnection(mongoClient)

	// if the user is not found
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Errors"] = []string{"User not found"}
	}

	// add the user data to the argumentsMap
	(*argumentsMap)["User"] = user

	return utils.CustomRenderTemplate(c, "profile/edit", *argumentsMap)
}
