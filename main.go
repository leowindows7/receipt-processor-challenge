package main

import (
	"examples.com/m/v2/endpoints"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	appGroup := app.Group("/receipts")
	appGroup.Post("/process", endpoints.ProcessRecipts)
	app.Listen(":3000")
}
