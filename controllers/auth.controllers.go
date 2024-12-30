package controllers

// the controllers for the auth routes
import (
	"example/web-server/config"
	"example/web-server/data"
	"example/web-server/utils"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

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
	c.Cookie(&fiber.Cookie{
		Name:     "Authenticated",
		Value:    "true",
		Secure:   true,
		HTTPOnly: true,
		// expires in 1 hour in data type time.Time and using our data.GetCookieExpirationTime() function
		Expires:     time.Now().Add(data.GetCookieExpirationTime() * time.Second),
		SameSite:    "Strict",
		SessionOnly: true,
	})

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

	// get the mongo client
	mongoClient := config.GetMongoClient()
	if mongoClient == nil {
		return utils.CustomRenderTemplate(c, "auth/register", data.GetFiberRenderMappingsAuthForms(email, password, nil, false))
	}

	// now we will save the user to the database
	// get the database and collection
	clientCollection, err := config.GetMongoCollection(mongoClient, "users")
	if err != nil {
		return utils.CustomRenderTemplate(c, "auth/register", data.GetFiberRenderMappingsAuthForms(email, password, nil, false))
	}

	// check if the user already exists with the same email
	err = clientCollection.FindOne(c.Context(), fiber.Map{"email": email}).Err()
	if err == nil {
		return utils.CustomRenderTemplate(c, "auth/register", data.GetFiberRenderMappingsAuthForms(email, password, &[]string{"User already exists"}, false))
	}

	// insert the user into the database
	_, err = clientCollection.InsertOne(c.Context(), fiber.Map{
		"email":        email,
		"passwordHash": string(hash),
	})
	if err != nil {
		return utils.CustomRenderTemplate(c, "auth/register", data.GetFiberRenderMappingsAuthForms(email, password, nil, false))
	}

	// close the mongo client connection
	config.CloseMongoClientConnection(mongoClient)

	// keep the same /register page but with a success flag
	return utils.CustomRenderTemplate(c, "auth/register", data.GetFiberRenderMappingsAuthForms(email, password, nil, true))
}

// AuthLogoutController handles the logout form submission
func AuthLogoutController(c *fiber.Ctx) error {
	// get the cookie before deleting it

	// delete the cookie
	c.Cookie(&fiber.Cookie{
		Name:     "Authenticated",
		Value:    "",
		Secure:   true,
		HTTPOnly: true,
		Expires:  time.Now().Add(-1 * time.Hour),
		SameSite: "Strict",
	})

	return c.Redirect("/auth/logout-success")
}
