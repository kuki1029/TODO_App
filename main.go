package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Get("/:name?", func(c *fiber.Ctx) error {
		if c.Params("name") != "" {
			return c.SendString("Hello " + c.Params("name"))
			// => Hello john
		}
		return c.SendString("Where is john?")
	})

	app.Listen(":8080")
}
