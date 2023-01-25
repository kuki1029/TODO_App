package controller

import (
	"github.com/gofiber/fiber/v2"

	"TODO/database"
	"TODO/models"
	"fmt"
	"strconv"
	"time"
)

// This function will create cookies for the user and log them in
// so that they can view their tasks.
func Login(ctx *fiber.Ctx) error {
	var creds models.User
	// First we need to parse the variable ctx to receive the credentials
	err := ctx.BodyParser(&creds)
	if err != nil {
		fmt.Println("Error with parsing credentials")
	}
	// Once we have their credentials, we need to check if the user is in the database
	err = database.AuthenticateUser(creds)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "No account found. Please signup first.",
		})
	}
	// To be able to identify this user on other pages, we need to create a cookie for their browser
	cookie := new(fiber.Cookie)
	cookie.Name = "userID"
	// This will be changed. This is just here as a placeholder so that the tasks can be worked on.
	cookie.Value = strconv.FormatUint(uint64(database.ReturnUserID(creds)), 10)
	cookie.Expires = time.Now().Add(24 * time.Hour)
	ctx.Cookie(cookie)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login successful.",
	})

}

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
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
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
