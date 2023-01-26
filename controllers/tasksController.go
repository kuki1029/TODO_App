package controller

import (
	"github.com/gofiber/fiber/v2"

	"TODO/models"
)

// This function will obtain the users tasks and then render them through fiber
// so that they can be displayed on the frontend
func DisplayTasks(ctx *fiber.Ctx) error {
	// First we obtain the models struct for tasks
	taskResponse := make([]models.TaskResponse, 0)
}
