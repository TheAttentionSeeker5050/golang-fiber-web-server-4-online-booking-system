package routers

import (
	"example/web-server/controllers"

	"github.com/gofiber/fiber/v2"
)

// LocationRoutes defines the API endpoints
func LocationRoutes(group fiber.Router) {
	// CRUD for locations
	group.Get("/", controllers.ListViewLocationsController)

	group.Get("/add", controllers.AddLocationController)

	group.Post("/add", controllers.AddLocationPostController)

	group.Get("/view/:id", controllers.DetailViewLocationController)

	group.Get("/edit/:id", controllers.EditLocationController)

	group.Post("/edit/:id", controllers.EditLocationPostController)

	group.Get("/delete/:id", controllers.DeleteLocationController)

	group.Post("/delete/:id", controllers.DeleteLocationPostController)

	group.Post("/bulk-delete", controllers.BulkDeleteLocationController)
}
