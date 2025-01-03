package controllers

// the controllers for the auth routes
import (
	"example/web-server/config"
	"example/web-server/data"
	"example/web-server/models"
	"example/web-server/utils"
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
		return utils.CustomRenderTemplate(c, "auth/login", utils.GetFiberRenderMappingsAuthForms(email, password, "", "", "", nil, false))
	}

	// now find the user with the email
	clientCollection, err := config.GetMongoCollection(mongoClient, "users")

	// if the collection is nil
	if clientCollection == nil {
		return utils.CustomRenderTemplate(c, "auth/login", utils.GetFiberRenderMappingsAuthForms(email, password, "", "", "", nil, false))
	}

	// find the user with the email
	var resultUser fiber.Map
	err = clientCollection.FindOne(c.Context(), fiber.Map{"email": email}).Decode(&resultUser)

	// if the user is not found
	if err != nil {
		return utils.CustomRenderTemplate(c, "auth/login", utils.GetFiberRenderMappingsAuthForms(email, password, "", "", "", &[]string{"User email or password is not correct"}, false))
	}

	// compare the password hash
	err = bcrypt.CompareHashAndPassword([]byte(resultUser["passwordHash"].(string)), []byte(password))

	// if the password is not correct
	if err != nil {
		return utils.CustomRenderTemplate(c, "auth/login", utils.GetFiberRenderMappingsAuthForms(email, password, "", "", "", &[]string{"Password is not correct"}, false))
	}

	// close the mongo client connection
	config.CloseMongoClientConnection(mongoClient)

	// generate a jwt token
	token, err := utils.GenerateLocalAuthJWTToken(resultUser["_id"].(string), resultUser["email"].(string))
	if err != nil {
		return utils.CustomRenderTemplate(c, "auth/login", utils.GetFiberRenderMappingsAuthForms(email, password, "", "", "", &[]string{"Error authenticating your account"}, false))
	}

	// add token and provider to the cookies
	utils.AddToCookies(c, "AccessToken", token, fiber.CookieSameSiteStrictMode)
	utils.AddToCookies(c, "TokenProvider", data.AUTH_PROVIDER_LOCAL, fiber.CookieSameSiteStrictMode)

	// keep the same /login page but with a success flag
	// do the same from now on
	return utils.CustomRenderTemplate(c, "auth/login", utils.GetFiberRenderMappingsAuthForms(email, password, "", "", "", nil, true))
}

// AuthRegisterController handles the registration form submission
func AuthRegisterController(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm-password")
	firstName := c.FormValue("first-name")
	lastName := c.FormValue("last-name")
	phoneNumber := c.FormValue("phone-number")

	// if passwords are the same
	if password != confirmPassword {
		return utils.CustomRenderTemplate(c, "auth/register", utils.GetFiberRenderMappingsAuthForms(email, password, firstName, lastName, phoneNumber, &[]string{"Passwords do not match"}, false))
	}

	// encrypt the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return utils.CustomRenderTemplate(c, "auth/register", utils.GetFiberRenderMappingsAuthForms(email, password, firstName, lastName, phoneNumber, nil, false))
	}

	mongoClient, userMongoCollection, err := models.GetUserCollection()
	if err != nil {
		return utils.CustomRenderTemplate(c, "auth/register", utils.GetFiberRenderMappingsAuthForms(email, password, firstName, lastName, phoneNumber, nil, false))
	}

	models.SaveUserToDBUsingLocalAuthProvider(c, userMongoCollection, models.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hash),
		AuthProvider: data.AUTH_PROVIDER_LOCAL,
		Picture:      "",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		FirstName:    firstName,
		LastName:     lastName,
		Phone:        phoneNumber,
	})

	// close the mongo client connection
	config.CloseMongoClientConnection(mongoClient)

	// keep the same /register page but with a success flag
	return utils.CustomRenderTemplate(c, "auth/register", utils.GetFiberRenderMappingsAuthForms(email, password, firstName, lastName, phoneNumber, nil, true))
}

// AuthLogoutController handles the logout form submission
func AuthLogoutController(c *fiber.Ctx) error {
	// delete the cookie using the utils function
	// utils.DeleteCookie(c, "Authenticated")
	utils.DeleteCookie(c, "Authenticated")
	utils.DeleteCookie(c, "AccessToken")
	utils.DeleteCookie(c, "TokenProvider")

	return c.Redirect("/auth/logout/success")
}
