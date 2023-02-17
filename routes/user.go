package routes

import (
	"todo/app/controller"

	"github.com/gofiber/fiber/v2"
)

// This function will create all the needed routes for our different pages
func SetupUserRoutes(app *fiber.App, uc *controller.UserController) {
	// Signup page for user
	app.Post("/signup", uc.Signup)
	// Login page for user
	app.Post("/login", uc.Login)
	// Logout
	app.Post("/logout", uc.Logout)

}
