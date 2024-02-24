package models

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []Item `json:"items"`
}

var receiptsMap map[string]Receipt

func checkPayload(receipt *Receipt) error {

	// numFields := value.NumField()
	fmt.Println(receipt.PurchaseTime)
	// structType := value.Type()
	// for i := 0; i < numFields; i++ {
	// 	field := structType.Field(i)
	// 	fmt.Println(field.Name)
	// }

	return nil
}

func ReceiptsProcessor(c *fiber.Ctx) (string, error) {
	receipt := new(Receipt)
	if err := c.BodyParser(receipt); err != nil {
		return "N/A", err
	}
	checkPayload(receipt)
	id := uuid.New()
	// fmt.Println(id.String())
	return id.String(), nil
}
