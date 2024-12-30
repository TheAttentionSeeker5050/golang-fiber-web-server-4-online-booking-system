package data

import "github.com/gofiber/fiber/v2"

// the path to the layout file
// errorMsg is optional
func GetFiberRenderMappingsAuthForms(email string, password string, errorMsgs *[]string, success bool) fiber.Map {
	var resultFiberMap fiber.Map = fiber.Map{
		"Errors":   []string{},
		"Email":    email,
		"Password": password,
		"Success":  success,
	}

	// if the size of the errorMsgs is 0
	if errorMsgs == nil {
		if !success {
			resultFiberMap["Errors"] = []string{"Internal server error"}
		}
	} else {
		resultFiberMap["Errors"] = *errorMsgs
	}

	return resultFiberMap
}
