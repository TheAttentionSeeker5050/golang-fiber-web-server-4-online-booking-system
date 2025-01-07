package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// helper function to create a token and add it to cookies, sameSite is optional
func AddToCookies(c *fiber.Ctx, key string, value string, sameSite string) {

	// if cookies sameSite is not provided or not allowed, return error
	if sameSite != "" && sameSite != fiber.CookieSameSiteLaxMode && sameSite != fiber.CookieSameSiteStrictMode && sameSite != fiber.CookieSameSiteNoneMode && sameSite != fiber.CookieSameSiteDisabled {
		// throw error if the sameSite is not allowed
		return // return nothing
	}
	// add the token to cookies
	c.Cookie(&fiber.Cookie{
		Name:     key,
		Value:    value,
		Secure:   true,
		HTTPOnly: true,
		// expires in 1 hour in data type time.Time and using our data.GetCookieExpirationTime() function
		Expires:     time.Now().Add(GetCookieExpirationTime() * time.Second),
		SameSite:    sameSite,
		SessionOnly: true,
	})
}

// helper function to read a cookie
func ReadCookie(c *fiber.Ctx, key string) string {
	cookie := c.Cookies(key)
	return cookie
}

// helper function to delete a cookie
func DeleteCookie(c *fiber.Ctx, key string) {
	// delete the cookie
	c.Cookie(&fiber.Cookie{
		Name:     key,
		Value:    "",
		Secure:   true,
		HTTPOnly: true,
		Expires:  time.Now().Add(-1 * time.Hour),
	})
}

// GetCookieExpirationTime returns time in seconds
func GetCookieExpirationTime() time.Duration {
	// get number of days using days from environment variable to int
	var numberOfDays string = os.Getenv("COOKIE_EXPIRES_IN_DAYS")
	// if can parse number of days to int
	if days, err := strconv.Atoi(numberOfDays); err == nil {
		// return time in seconds
		return time.Duration(60 * 60 * 24 * days)
	}

	return time.Duration(60 * 60 * 24 * 7)
}
