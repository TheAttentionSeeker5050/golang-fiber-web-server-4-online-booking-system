package controllers

import (
	"context"
	"encoding/json"
	"example/web-server/config"
	"example/web-server/data"
	"example/web-server/models"
	"example/web-server/utils"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Oauth2 providers controllers
func SignInWithGoogleController(c *fiber.Ctx) error {

	// generate unique state code
	var state string = uuid.New().String()
	authCodeUrl := config.OauthConfig.GoogleLoginConfig.AuthCodeURL(state)

	c.Status(fiber.StatusSeeOther)
	c.Redirect(authCodeUrl)

	return c.JSON(authCodeUrl)
}

func SignWithGoogleCallbackController(c *fiber.Ctx) error {
	// to prevent double execution of the callback, we will check the state
	state := c.Query("state")
	if state == "" {
		return c.SendString("Invalid state parameter")
	}

	// Mark the state as used
	stateCache.Store(state, true)

	code := c.Query("code")

	if code == "" {
		return c.SendString("Invalid code parameter")
	}

	googleCon := config.InitGoogleConfig()

	token, err := googleCon.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Code-Token Exchange Failed")
	}

	// Fetch user info
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken
	response, err := http.Get(userInfoURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch user info")
	}

	defer response.Body.Close()

	// _, err = ioutil.ReadAll(response.Body)
	userData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read user info")
	}

	// add the token to cookies
	utils.AddToCookies(c, "AccessToken", token.AccessToken, data.COOKIE_SAME_SITE_LAX)
	utils.AddToCookies(c, "TokenProvider", "Google", data.COOKIE_SAME_SITE_LAX)

	// get user email from the token and save it to the database if it doesn't exist
	// get the mongo client
	mongoClient, userMongoCollection, err := models.GetUserCollection()

	// if the mongo client is nil
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user collection")
	}

	userDataMap := fiber.Map{}
	err = json.Unmarshal(userData, &userDataMap)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to unmarshal user data")
	}

	// get the user by email
	userFromDB := models.User{}
	err = userMongoCollection.FindOne(c.Context(), fiber.Map{
		"email": userDataMap["email"],
	}).Decode(&userFromDB) // this should return an error if the user is not found
	if err != nil {
		// save the user to the database
		err = models.SaveUserToDBUsingGoogleProvider(c, userMongoCollection, &utils.GoogleClaims{
			ID:            userDataMap["id"].(string),
			Email:         userDataMap["email"].(string),
			EmailVerified: true,
			Sub:           userDataMap["id"].(string),
			Name:          userDataMap["name"].(string),
			Picture:       userDataMap["picture"].(string),
		}) // this should return an error if the user could not be saved
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save user to database")
		}
	}

	// close the mongo client connection
	config.CloseMongoClientConnection(mongoClient)

	// redirect to the success page
	return c.Redirect("/auth/google/success")
}

func SignWithGoogleSuccessController(c *fiber.Ctx) error {
	// make a fiber mapping for displaying the success page template
	argumentsMap := &fiber.Map{
		"Title":    "Google Login Success",
		"Provider": "Google",
	}

	// render using utility function
	return utils.CustomRenderTemplate(c, "auth/oauth-success", *argumentsMap)
}
