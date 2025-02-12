package controllers

import (
	"github.com/gofiber/fiber/v2"
)

// Make controller functions for the location routes
func DetailViewLocationController(c *fiber.Ctx) error {

	// For now return a dummy page
	return c.SendString("Location Details")
}

func AddLocationController(c *fiber.Ctx) error {

	// For now return a dummy page
	return c.SendString("Create Location")
}

func AddLocationPostController(c *fiber.Ctx) error {

	// For now return a redirect to the location details page
	return c.Redirect("/locations")
}

func EditLocationController(c *fiber.Ctx) error {

	// For now return a dummy page
	return c.SendString("Edit Location")
}

func EditLocationPostController(c *fiber.Ctx) error {

	// For now return a redirect to the location details page
	return c.Redirect("/locations")
}

func DeleteLocationController(c *fiber.Ctx) error {

	// For now return a redirect to the Location details page
	return c.Redirect("/locations")
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

	// For now return a dummy page
	return c.SendString("List Locations")
}
