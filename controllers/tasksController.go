package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"TODO/database"
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
