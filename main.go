package main

import (
	"fmt"

	"examples.com/m/v2/endpoints"
	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("server is happy!")
	app := fiber.New()
	// test Get when server is up
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	// begining of this challenge
	appGroup := app.Group("/receipts")
	appGroup.Post("/process", endpoints.ProcessRecipts)
	appGroup.Get("/:id/points", endpoints.GetPoints)
	app.Listen(":3000")

}
