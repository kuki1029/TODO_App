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
	// Obtain the ID by checking if the cookies sessionKey exists in cache or not
	key := ctx.Cookies("sessionKey")
	ID, err := database.GetFromRedis(client, key)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	}
	taskResponse, err := database.ReturnTasksWithID(ID)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	}
	name, err := database.ReturnName(ID)
	fmt.Println(name)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	} else {
		return ctx.Render("tasks", fiber.Map{
			"Tasks": taskResponse,
			"Name":  name,
		})
	}
}

// This function will add the task to the database
func AddTasks(ctx *fiber.Ctx) error {
	// Obtain the ID by checking if the cookies sessionKey exists in cache or not
	key := ctx.Cookies("sessionKey")
	ID, err := database.GetFromRedis(client, key)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	}
	tempTask := models.Task{}
	// Get the details of the task
	err = ctx.BodyParser(&tempTask)
	if err != nil {
		fmt.Println("Error with parsing credentials")
	}
	// Add task to database
	task := models.Task{
		TaskName: tempTask.TaskName,
		UserID:   ID,
	}
	primaryID, err := database.AddTask(task)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	} else {
		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"ID":      primaryID,
		})
	}
}

// This function will delete the tasks from the database
func DelTask(ctx *fiber.Ctx) error {
	// Obtain the ID by checking if the cookies sessionKey exists in cache or not
	key := ctx.Cookies("sessionKey")
	// We don't need the user ID here, but we still need to verify
	_, err := database.GetFromRedis(client, key)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	}
	// Parse the ID and convert it to int
	num, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	ID := uint(num)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	}
	// Call the database function to delete the task
	err = database.DelTask(ID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
		})
	}
}

// This function will mark the task as done in the database
func TaskDone(ctx *fiber.Ctx) error {
	// Obtain the ID by checking if the cookies sessionKey exists in cache or not
	key := ctx.Cookies("sessionKey")
	// We don't need the user ID here, but we still need to verify
	_, err := database.GetFromRedis(client, key)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	}
	// Parse the ID and convert it to int
	num, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	ID := uint(num)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	}
	// Call the database function to mark the task as done
	err = database.MarkTaskDone(ID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
		})
	}
}
