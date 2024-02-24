package endpoints

import "github.com/gofiber/fiber/v2"

func ProcessRecipts(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{"id": 123})

}
