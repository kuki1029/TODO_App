package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"

	"todo/controller"
	"todo/repo"
)

// This function will create all the needed routes for our different pages
func setupRoutes(app *fiber.App) {
	// Signup page for user
	app.Post("/signup", controller.Signup)
	// Login page for user
	app.Post("/login", controller.Login)
	// Show tasks to user
	app.Get("/tasks", controller.DisplayTasks)
	// Add task to database
	app.Post("/tasks", controller.AddTasks)
	// Delete tasks from database
	app.Delete("/tasks/:id", controller.DelTask)
	// Mark a task done given a certain id
	app.Post("/tasksDone/:id", controller.TaskDone)
	// Logout
	app.Post("/logout", controller.Logout)
}

// Stop the Fiber application
func exit(app *fiber.App) {
	_ = app.Shutdown()
}

func main() {
	// Setup the database
	repo.ConnectToDB()
	repo.ConnectToDBTasks()
	// Create a new engine
	engine := html.New("./views", ".html")

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// This serves the css files so it the HTML can render it
	app.Static("/static", "./static")
	// Serves all the HTML filse
	app.Static("/", "./views", fiber.Static{
		Index: "login.html",
	})

	setupRoutes(app)

	// Close any connections on interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		exit(app)
	}()

	// Start listening on the specified address
	if err := app.Listen("127.0.0.1:8080"); err != nil {
		log.Panic(err)
	}
}
