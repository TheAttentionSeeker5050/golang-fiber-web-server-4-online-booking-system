package controllers

import (
	"example/web-server/data"

	"github.com/gofiber/fiber/v2"
)

// the error controllers callback function
func RouterErrorCallback(ctx *fiber.Ctx, err error) error {

	// Status code defaults to 500
	var code int = fiber.StatusInternalServerError

	// get the route path as string
	var path string = ctx.Path()

	// Retrieve the custom status code
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// if path starts with /api, return json error
	if len(path) >= 4 {
		if path[:4] == "/api" {
			if code == 404 {
				return APIErrorController(ctx, "Route not found", code)
			} else {
				return APIErrorController(ctx, err.Error(), code)
			}
		}
	}

	if code == 404 {
		// use handlerbars template to render 404 page
		return NotFoundErrorPageController(ctx)
	} else if code == 401 {
		// use handlerbars template to render 401 page
		return UnauthorizedErrorPageController(ctx)
	}

	// return handlebars error page in /views/error using the error message as page argument
	return OtherErrorPageController(ctx, err.Error())

}

// not found error page controller to load the /404 page
func NotFoundErrorPageController(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotFound).Render("404", fiber.Map{
		"Title":     "404 Not Found",
		"Paragraph": "The page you are looking for does not exist.",
	}, data.LAYOUT_PATH)
}

func UnauthorizedErrorPageController(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusUnauthorized).Render("401", fiber.Map{
		"Title":     "Unauthorized",
		"Paragraph": "You are not authorized to access this page.",
	}, data.LAYOUT_PATH)
}

// API error controller to return a json response with error from ctx
func APIErrorController(ctx *fiber.Ctx, errorMsg string, status int) error {
	return ctx.Status(status).JSON(fiber.Map{
		"error": errorMsg,
	})
}

// Other error page controller
func OtherErrorPageController(ctx *fiber.Ctx, errorMsg string) error {
	return ctx.Status(fiber.StatusInternalServerError).Render("error", fiber.Map{
		"Message": errorMsg,
	}, data.LAYOUT_PATH)
}
