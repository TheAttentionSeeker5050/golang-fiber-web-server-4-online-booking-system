package controllers

import (
	"github.com/gofiber/fiber/v2"
)

// Make controller functions for the reservation routes
func DetailViewReservationController(c *fiber.Ctx) error {

	// For now return a dummy page
	return c.SendString("Reservation Details")
}

func AddReservationController(c *fiber.Ctx) error {

	// For now return a dummy page
	return c.SendString("Create Reservation")
}

func AddReservationPostController(c *fiber.Ctx) error {

	// For now return a redirect to the reservation details page
	return c.Redirect("/reservations")
}

func EditReservationController(c *fiber.Ctx) error {

	// For now return a dummy page
	return c.SendString("Edit Reservation")
}

func EditReservationPostController(c *fiber.Ctx) error {

	// For now return a redirect to the reservation details page
	return c.Redirect("/reservations")
}

func DeleteReservationController(c *fiber.Ctx) error {

	// For now return a redirect to the Reservation details page
	return c.Redirect("/reservations")
}

func DeleteReservationPostController(c *fiber.Ctx) error {

	// For now return a redirect to the Reservation details page
	return c.Redirect("/reservations")
}

func BulkDeleteReservationController(c *fiber.Ctx) error {

	// For now return a redirect to the Reservation details page
	return c.Redirect("/reservations")
}

func ListViewReservationsController(c *fiber.Ctx) error {

	// For now return a dummy page
	return c.SendString("List Reservations")
}
