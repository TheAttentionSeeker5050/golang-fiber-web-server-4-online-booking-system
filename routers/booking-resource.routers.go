package routers

import (
	"example/web-server/controllers"

	"github.com/gofiber/fiber/v2"
)

// BookingResourceRoutes defines the API endpoints
func BookingResourceRoutes(group fiber.Router) {
	// CRUD for booking resources
	group.Get("/", controllers.ListViewBookingResourcesController)

	group.Get("/add", controllers.AddBookingResourceController)

	group.Post("/add", controllers.AddBookingResourcePostController)

	group.Get("/view/:id", controllers.DetailViewBookingResourceController)

	group.Get("/edit/:id", controllers.EditBookingResourceController)

	group.Post("/edit/:id", controllers.EditBookingResourcePostController)

	group.Get("/delete/:id", controllers.DeleteBookingResourceController)

	group.Post("/delete/:id", controllers.DeleteBookingResourcePostController)

	group.Post("/bulk-delete", controllers.BulkDeleteBookingResourceController)
}
