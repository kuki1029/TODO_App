package controller

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"TODO/database"
	"TODO/models"
)

// This function will obtain the users tasks and then render them through fiber
// so that they can be displayed on the frontend
func DisplayTasks(ctx *fiber.Ctx) error {
	// Temp method to obtain userID.
	userID, _ := strconv.ParseUint(ctx.Cookies("userID"), 10, 64)
	ID := uint(userID)
	taskResponse, err := database.ReturnTasksWithID(ID)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	} else {
		return ctx.Render("tasks", fiber.Map{
			"Tasks": taskResponse,
		})
	}
}

// This function will add the task to the database
func AddTasks(ctx *fiber.Ctx) error {
	// Temp method to obtain userID.
	userID, _ := strconv.ParseUint(ctx.Cookies("userID"), 10, 64)
	ID := uint(userID)
	tempTask := models.Task{}
	// Get the details of the task
	err := ctx.BodyParser(&tempTask)
	if err != nil {
		fmt.Println("Error with parsing credentials")
	}
	// Add task to database
	task := models.Task{
		TaskName: tempTask.TaskName,
		UserID:   ID,
	}
	err = database.AddTask(task, ID)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	} else {
		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
		})
	}
}
