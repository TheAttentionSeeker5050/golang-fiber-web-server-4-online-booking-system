package controllers

import (
	"example/web-server/config"
	"example/web-server/models"
	"example/web-server/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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

func EditProfilePostController(c *fiber.Ctx) error {
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

	// if the user is not found
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Errors"] = []string{"User not found"}
	}

	// add the user data to the argumentsMap
	(*argumentsMap)["User"] = user

	// get the form values
	firstName := c.FormValue("first-name")
	lastName := c.FormValue("last-name")
	phoneNumber := c.FormValue("phone-number")

	// update the user data
	user.FirstName = firstName
	user.LastName = lastName
	user.Phone = phoneNumber

	// update the user data
	err = models.UpdateUser(c, userMongoCollection, user)

	// close the mongo connection
	config.CloseMongoClientConnection(mongoClient)

	// if the user data is not updated
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Errors"] = []string{"Error updating user data"}
		return utils.CustomRenderTemplate(c, "profile/edit", *argumentsMap)
	}

	// add success msg to the argumentsMap
	(*argumentsMap)["Success"] = "User data updated successfully"

	return utils.CustomRenderTemplate(c, "profile/edit", *argumentsMap)
}

func EditPasswordPostController(c *fiber.Ctx) error {
	argumentsMap := &fiber.Map{
		"Title": "Edit Password",
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

	// get the user password hash
	userPasswordHash, err := models.GetUserPasswordHash(c, userMongoCollection, userID)

	// if the user is not found
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Errors"] = []string{"User not found"}
	}

	// get the form values
	oldPassword := c.FormValue("current-password")
	newPassword := c.FormValue("new-password")
	newPasswordConfirm := c.FormValue("confirm-new-password")

	// compare the old password with PasswordHash using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(userPasswordHash), []byte(oldPassword))
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Errors"] = []string{"Old password is incorrect"}
		return utils.CustomRenderTemplate(c, "profile/edit-password", *argumentsMap)
	}

	// if the new password and new password confirm do not match
	if newPassword != newPasswordConfirm {
		// add error msg to the argumentsMap
		(*argumentsMap)["Errors"] = []string{"New password and confirm password do not match"}
		return utils.CustomRenderTemplate(c, "profile/edit-password", *argumentsMap)
	}

	// hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	// if there is an error hashing the password
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Errors"] = []string{"Error saving the new password"}
	}

	// update the user password
	err = models.UpdateUserPassword(c, userMongoCollection, userID, string(hashedPassword))

	// close the mongo connection
	config.CloseMongoClientConnection(mongoClient)

	// if the user data is not updated
	if err != nil {
		// add error msg to the argumentsMap
		(*argumentsMap)["Errors"] = []string{"Error updating user data"}
		return utils.CustomRenderTemplate(c, "profile/edit-password", *argumentsMap)
	}

	// add success msg to the argumentsMap
	(*argumentsMap)["Success"] = "Password updated successfully"

	return utils.CustomRenderTemplate(c, "profile/edit-password", *argumentsMap)
}
