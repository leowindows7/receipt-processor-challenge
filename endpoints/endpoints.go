package endpoints

import (
	"examples.com/m/v2/models"
	"github.com/gofiber/fiber/v2"
)

func ProcessRecipts(c *fiber.Ctx) error {

	receiptStruct := models.Receipt{}
	c.BodyParser(&receiptStruct)
	id, err := models.ReceiptsProcessor()
	if err != nil {
		return c.JSON(fiber.Map{"err": err.Error()})
	}
	return c.JSON(fiber.Map{"id": id})

}
