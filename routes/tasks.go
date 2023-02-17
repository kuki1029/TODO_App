package routes

import (
	"todo/app/controller"

	"github.com/gofiber/fiber/v2"
)

// This function will create all the needed routes for our different pages
func SetupTaskRoutes(app *fiber.App, tc *controller.TaskController) {
	// Show tasks to user
	app.Get("/tasks", tc.DisplayTasks)
	// Add task to database
	app.Post("/tasks", tc.AddTasks)
	// Delete tasks from database
	app.Delete("/tasks/:id", tc.DelTask)
	// Mark a task done given a certain id
	app.Post("/tasksDone/:id", tc.TaskDone)
	// Edit task in database
	app.Post("/tasksEdit/:id", tc.EditTask)
}
