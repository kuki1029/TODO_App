package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"

	"todo/app/controller"
	"todo/app/repo"
)

// This function will create all the needed routes for our different pages
func setupRoutes(app *fiber.App, uc *controller.UserController, tc *controller.TaskController) {
	// Signup page for user
	app.Post("/signup", uc.Signup)
	// Login page for user
	app.Post("/login", uc.Login)
	// Show tasks to user
	app.Get("/tasks", tc.DisplayTasks)
	// Add task to database
	app.Post("/tasks", tc.AddTasks)
	// Delete tasks from database
	app.Delete("/tasks/:id", tc.DelTask)
	// Mark a task done given a certain id
	app.Post("/tasksDone/:id", tc.TaskDone)
	// Logout
	app.Post("/logout", uc.Logout)
	// Edit task in database
	app.Post("/tasksEdit/:id", tc.EditTask)
}

// Stop the Fiber application
func exit(app *fiber.App) {
	_ = app.Shutdown()
}

func main() {
	// Setup the database
	repo.ConnectToDB()
	// Setup repo and controller interfaces
	ur := repo.NewUserRepo(repo.DB.DbConn)
	uc := controller.NewUserController(*ur)
	tr := repo.NewTaskRepo(repo.DB.DbConn)
	tc := controller.NewTaskController(*tr)
	// Create a new engine
	engine := html.New("./resources/views", ".html")

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// This serves the css files so it the HTML can render it
	app.Static("/static", "./static")
	app.Static("/resources/JS", "./resources/JS")
	// Serves all the HTML filse
	app.Static("/", "./resources/views", fiber.Static{
		Index: "login.html",
	})

	setupRoutes(app, uc, tc)

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
