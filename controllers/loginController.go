package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"TODO/database"
	"TODO/models"
	"fmt"
)

// This function will be called through the JS and handle any signup requirements
// the c variable here contains all the required credentials
func Signup(ctx *fiber.Ctx) error {
	var creds models.User
	// First we need to parse the variable ctx to receive the credentials
	err := ctx.BodyParser(&creds)
	if err != nil {
		fmt.Println("Error with parsing credentials")
	}

	// Once we have the required data, we need to make sure the user isn't a duplicate
	err = database.AddUser(creds)
	if err != nil {
		// In this case, we know the user is a duplicate, so we returen an error message
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "User already exists. Please login or use a different email address.",
		})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Account created.",
		})
	}
}
