package controller

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"todo/middleware"
	"todo/models"
	"todo/repo"
)

// This function will obtain the users tasks and then render them through fiber
// so that they can be displayed on the frontend
func DisplayTasks(ctx *fiber.Ctx) error {
	db := repo.DB.DbConn
	// Obtain the ID by checking if the cookies sessionKey exists in cache or not
	key := ctx.Cookies("sessionKey")
	ID, err := middleware.GetFromRedis(client, key)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	}
	taskResponse, err := repo.ReturnTasksWithID(ID, db)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	}
	name, err := repo.ReturnName(ID, db)

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
	db := repo.DB.DbConn
	// Obtain the ID by checking if the cookies sessionKey exists in cache or not
	key := ctx.Cookies("sessionKey")
	ID, err := middleware.GetFromRedis(client, key)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	}
	tempTask := models.TaskDTO{}
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
	primaryID, err := repo.AddTask(task, db)
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
	db := repo.DB.DbConn
	// Obtain the ID by checking if the cookies sessionKey exists in cache or not
	key := ctx.Cookies("sessionKey")
	// We don't need the user ID here, but we still need to verify
	_, err := middleware.GetFromRedis(client, key)
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
	err = repo.DelTask(ID, db)
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
	db := repo.DB.DbConn
	// Obtain the ID by checking if the cookies sessionKey exists in cache or not
	key := ctx.Cookies("sessionKey")
	// We don't need the user ID here, but we still need to verify
	_, err := middleware.GetFromRedis(client, key)
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
	err = repo.MarkTaskDone(ID, db)
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

// The function will edit the task by calling the appropraite database function
func EditTask(ctx *fiber.Ctx) error {
	db := repo.DB.DbConn
	// Obtain the ID by checking if the cookies sessionKey exists in cache or not
	key := ctx.Cookies("sessionKey")
	// We don't need the user ID here, but we still need to verify
	_, err := middleware.GetFromRedis(client, key)
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
	tempTask := models.TaskDTO{}
	// Get the details of the task
	err = ctx.BodyParser(&tempTask)
	if err != nil {
		fmt.Println("Error with parsing credentials")
	}
	// Call the database function to mark the task as done
	err = repo.EditTask(ID, tempTask.TaskName, db)
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
