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

var receiptsMap map[string]Receipt = make(map[string]Receipt)

func checkPayload(receipt *Receipt) error {
	receiptStruct := reflect.Indirect(reflect.ValueOf(receipt))
	numFields := receiptStruct.NumField()
	for i := 0; i < numFields; i++ {
		if receiptStruct.Field(i).Interface() == "" {

			return errors.New("please check your payload, missing required fields")
		}
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
	idStr := id.String()
	receiptsMap[idStr] = *receipt
	for key, _ := range receiptsMap {
		fmt.Println(key)
	}
	return idStr, nil
}

func PointsCalculator(id string) (int64, error) {
	receiptToCheck, ok := receiptsMap[id]
	if !ok {
		return -1, errors.New("id not exists")
	}
	receiptStruct := reflect.ValueOf(receiptToCheck)
	numFields := receiptStruct.NumField()
	for i := 0; i < numFields; i++ {
		fieldName := receiptStruct.Type().Field(i).Name
		switch fieldName {
		case "Retailer":
			fmt.Println(fieldName)
		case "PurchaseDate":
			fmt.Println(fieldName)
		case "PurchaseTime":
			fmt.Println(fieldName)
		case "Total":
			fmt.Println(fieldName)
		case "Items":
			fmt.Println(fieldName)
		}

	}
	fmt.Println(receiptToCheck)
	return 100, nil

}

// func ruleRetailerName(name string) int64 {
// 	return 0
// }

// func ruleRoundDollar(name string) int64 {
// 	return 0
// }

// func ruleTotalMultipleOf25(name string) int64 {
// 	return 0
// }

// func ruleEveryTwoItems(name string) int64 {
// 	return 0
// }
// func ruleItemDescriptionLength(name string) int64 {

// 	return 0
// }
// func rulePurchaseDate(name string) int64 {
// 	return 0
// }

// func ruleTimeOfPurchase(name string) int64 {
// 	return 0
// }
