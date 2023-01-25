package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"

	//"github.com/joho/godotenv"

	controller "TODO/controllers"
	"TODO/database"
)

// This function will create all the needed routes for our different pages
func setupRoutes(app *fiber.App) {
	// Signup page for user
	app.Post("/signup", controller.Signup)

	// Login page for user
	app.Post("/login", controller.Login)
	/*
		// Logout
		app.Post("/logout", logoutFunc)

		// Show tasks to user
		app.Get("/tasks", tasksFunc) */
}

func main() {
	// Setup the database
	database.ConnectToDB()
	// Create a new engine
	engine := html.New("./views", ".html")

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// This serves the css files so it the HTML can render it
	app.Static("/static", "./static")
	// Servers all the HTML filse
	app.Static("/", "./views", fiber.Static{
		Index: "login.html",
	})

	setupRoutes(app)
	/*
		app.Get("/", func(c *fiber.Ctx) error {
			// Render index
			return c.Render("login", fiber.Map{
				"Title": "Login",
			})
		})*/

	app.Listen("127.0.0.1:8080")
}
