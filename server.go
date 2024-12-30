package main

import (
	"fmt"
	"log"

	"example/web-server/config"
	"example/web-server/controllers"
	"example/web-server/data"
	"example/web-server/middlewares"
	"example/web-server/routers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/handlebars/v2"
)

func main() {

	// call the GetMongoConfig function from the config package
	mongoClient := config.GetMongoClient()
	if mongoClient == nil {
		log.Fatal("Could not connect to MongoDB")
	}

	// call the defer function to close the mongo client connection
	config.CloseMongoClientConnection(mongoClient)

	// Create a new engine
	engine := handlebars.New("./views", ".hbs")

	// Initialize the golang fiber app
	app := fiber.New(fiber.Config{
		// Prefork: true, // for production to support spawning multiple processes
		// Concurrency: 256 * 1024, // set the maximum number of concurrent connections
		AppName: "Online Booking System",
		// use controllers.routerErrorCallback as the error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return controllers.RouterErrorCallback(ctx, err)
		},
		Views: engine,
		// PassLocalsToViews: true, // this is done so we can pass the render condition to the views from the auth middleware, in the case for example we want to show a different header if the user is authenticated than if not
		ViewsLayout: data.LAYOUT_PATH,
	})

	// use cors
	app.Use(cors.New(cors.Config{
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	// add middleware that validates cookie and flags the user for authentication valid or not so we can render custom things from the controllers routes contents
	app.Use(middlewares.UserIsAuthorized)

	app.Static("/public", "./public") // => Serve static files from ./public

	// Register routes
	homeGroup := app.Group("/")
	routers.HomeRoutes(homeGroup)
	apiGroup := app.Group("/api/v1")
	routers.APIRoutes(apiGroup)
	authGroup := app.Group("/auth")
	routers.AuthRoutes(authGroup)

	fmt.Println(app.Stack())
	log.Fatal(app.Listen(":8080"))
}
