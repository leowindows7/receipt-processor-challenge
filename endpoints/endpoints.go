package endpoints

import (
	"examples.com/m/v2/models"
	"github.com/gofiber/fiber/v2"
)

func ProcessRecipts(c *fiber.Ctx) error {

	id, err := models.ReceiptsProcessor(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"id": id})

}
