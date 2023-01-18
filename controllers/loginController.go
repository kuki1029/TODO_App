package controller

import (
	"TODO/models"

	"github.com/gofiber/fiber/v2"
)

// This function will be called through the JS and handle any signup requirements
// the c variable here contains all the required credentials
func Signup(c *fiber.Ctx) error {
	var creds models.User
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
