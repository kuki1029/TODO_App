package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"

	"TODO/controllers"
)

// "github.com/joho/godotenv"

// This function will create all the needed routes for our different pages
func setupRoutes(app *fiber.App) {
	// Login page for user
	app.Post("/login", controller.Login)

	// Signup page for user
	app.Post("/signup", signupFunc)

	// Logout
	app.Post("/logout", logoutFunc)

	// Show tasks to user
	app.Get("/tasks", tasksFunc)
}

func main() {
	// Create a new engine

	engine := html.New("./views", ".html")

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// This serves the css files so it the HTML can render it
	app.Static("/static", "./static")

	setupRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("login", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	app.Listen(":8080")
}
