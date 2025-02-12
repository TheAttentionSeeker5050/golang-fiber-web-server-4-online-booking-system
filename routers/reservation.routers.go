package routers

import (
	"example/web-server/controllers"

	"github.com/gofiber/fiber/v2"
)

// ReservationRoutes defines the API endpoints
func ReservationRoutes(group fiber.Router) {
	// CRUD for reservations
	group.Get("/", controllers.ListViewReservationsController)

	group.Get("/add", controllers.AddReservationController)

	group.Post("/add", controllers.AddReservationPostController)

	group.Get("/view/:id", controllers.DetailViewReservationController)

	group.Get("/edit/:id", controllers.EditReservationController)

	group.Post("/edit/:id", controllers.EditReservationPostController)

	group.Get("/delete/:id", controllers.DeleteReservationController)

	group.Post("/delete/:id", controllers.DeleteReservationPostController)

	group.Post("/bulk-delete", controllers.BulkDeleteReservationController)
}
