package controller

import (
	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
