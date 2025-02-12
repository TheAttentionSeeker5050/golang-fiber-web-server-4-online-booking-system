package controllers

import (
	"github.com/gofiber/fiber/v2"
)

// Make controller functions for the booking resource routes
func DetailViewBookingResourceController(c *fiber.Ctx) error {

	// For now return a dummy page
	return c.SendString("Booking Resource Details")
}

func AddBookingResourceController(c *fiber.Ctx) error {

	// For now return a dummy page
	return c.SendString("Create Booking Resource")
}

func AddBookingResourcePostController(c *fiber.Ctx) error {

	// For now return a redirect to the booking resource details page
	return c.Redirect("/booking-resources")
}

func EditBookingResourceController(c *fiber.Ctx) error {

	// For now return a dummy page
	return c.SendString("Edit Booking Resource")
}

func EditBookingResourcePostController(c *fiber.Ctx) error {

	// For now return a redirect to the booking resource details page
	return c.Redirect("/booking-resources")
}

func DeleteBookingResourceController(c *fiber.Ctx) error {

	// For now return a redirect to the Booking Resource details page
	return c.Redirect("/booking-resources")
}

func DeleteBookingResourcePostController(c *fiber.Ctx) error {

	// For now return a redirect to the Booking Resource details page
	return c.Redirect("/booking-resources")
}

func BulkDeleteBookingResourceController(c *fiber.Ctx) error {

	// For now return a redirect to the Booking Resource details page
	return c.Redirect("/booking-resources")
}

func ListViewBookingResourcesController(c *fiber.Ctx) error {

	// For now return a dummy page
	return c.SendString("List Booking Resources")
}
