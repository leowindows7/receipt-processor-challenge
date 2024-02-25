package models

import (
	"errors"
	"fmt"
	"reflect"

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
	receiptStruct := reflect.Indirect(reflect.ValueOf(receipt))
	types := receiptStruct.Type()
	numFields := receiptStruct.NumField()
	for i := 0; i < numFields; i++ {
		if receiptStruct.Field(i).Interface() == "" {

			return errors.New("please check your payload, missing required fields")
		}

		fmt.Println(types.Field(i).Name, receiptStruct.Field(i))
	}

	return nil
}

func ReceiptsProcessor(c *fiber.Ctx) (string, error) {
	receipt := new(Receipt)
	if err := c.BodyParser(receipt); err != nil {
		return "N/A", err
	} else if checkPayload(receipt) != nil {

		return "N/A", checkPayload(receipt)
	}

	id := uuid.New()

	return id.String(), nil
}
