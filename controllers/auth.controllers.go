package controllers

// the controllers for the auth routes
import (
	"context"
	"encoding/json"
	"example/web-server/config"
	"example/web-server/data"
	"example/web-server/models"
	"example/web-server/utils"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var stateCache = sync.Map{}

// AuthLoginController handles the login form submission
func AuthLoginController(c *fiber.Ctx) error {

	// get the form values
	email := c.FormValue("email")
	password := c.FormValue("password")

	// get the mongo client
	mongoClient := config.GetMongoClient()

	// if the mongo client is nil
	if mongoClient == nil {
		return utils.CustomRenderTemplate(c, "auth/login", data.GetFiberRenderMappingsAuthForms(email, password, nil, false))
	}

	// now find the user with the email
	clientCollection, err := config.GetMongoCollection(mongoClient, "users")

	// if the collection is nil
	if clientCollection == nil {
		log.Println("Mongo collection is not available")

		return utils.CustomRenderTemplate(c, "auth/login", data.GetFiberRenderMappingsAuthForms(email, password, nil, false))
	}

	// find the user with the email
	var resultUser fiber.Map
	err = clientCollection.FindOne(c.Context(), fiber.Map{"email": email}).Decode(&resultUser)

	// if the user is not found
	if err != nil {
		log.Println("User not found")

		return utils.CustomRenderTemplate(c, "auth/login", data.GetFiberRenderMappingsAuthForms(email, password, &[]string{"User email or password is not correct"}, false))
	}

	// compare the password hash
	err = bcrypt.CompareHashAndPassword([]byte(resultUser["passwordHash"].(string)), []byte(password))

	// if the password is not correct
	if err != nil {
		log.Println("Password is not correct")

		return utils.CustomRenderTemplate(c, "auth/login", data.GetFiberRenderMappingsAuthForms(email, password, &[]string{"Password is not correct"}, false))
	}

	// close the mongo client connection
	config.CloseMongoClientConnection(mongoClient)

	// make a cookie with auth to true for now, as we will use something secure later
	// it should be http only and secure as we don't want to expose it to the client
	// use the utils function to add the cookie
	utils.AddToCookies(c, "Authenticated", "true", fiber.CookieSameSiteStrictMode)

	// keep the same /login page but with a success flag
	// do the same from now on
	return utils.CustomRenderTemplate(c, "auth/login", data.GetFiberRenderMappingsAuthForms(email, password, nil, true))
}

// AuthRegisterController handles the registration form submission
func AuthRegisterController(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm-password")

	// if passwords are the same
	if password != confirmPassword {
		return utils.CustomRenderTemplate(c, "auth/register", data.GetFiberRenderMappingsAuthForms(email, password, &[]string{"Passwords do not match"}, false))
	}

	// encrypt the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return utils.CustomRenderTemplate(c, "auth/register", data.GetFiberRenderMappingsAuthForms(email, password, nil, false))
	}

	mongoClient, userMongoCollection, err := models.GetUserCollection()
	if err != nil {
		return utils.CustomRenderTemplate(c, "auth/register", data.GetFiberRenderMappingsAuthForms(email, password, nil, false))
	}

	models.SaveUserToDBUsingLocalAuthProvider(c, userMongoCollection, models.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hash),
		AuthProvider: "Local",
		Picture:      "",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})

	// close the mongo client connection
	config.CloseMongoClientConnection(mongoClient)

	// keep the same /register page but with a success flag
	return utils.CustomRenderTemplate(c, "auth/register", data.GetFiberRenderMappingsAuthForms(email, password, nil, true))
}

// AuthLogoutController handles the logout form submission
func AuthLogoutController(c *fiber.Ctx) error {
	// delete the cookie using the utils function
	// utils.DeleteCookie(c, "Authenticated")
	utils.DeleteCookie(c, "Authenticated")
	utils.DeleteCookie(c, "AccessToken")
	utils.DeleteCookie(c, "TokenProvider")

	return c.Redirect("/auth/logout-success")
}

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
	utils.AddToCookies(c, "AccessToken", token.AccessToken, fiber.CookieSameSiteLaxMode)
	utils.AddToCookies(c, "TokenProvider", "Google", fiber.CookieSameSiteLaxMode)

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
		err = models.SaveUserToDBUsingGoogleProvider(c, userMongoCollection, &models.GoogleClaims{
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
