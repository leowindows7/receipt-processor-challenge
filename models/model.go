package models

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"

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

	receiptVal := reflect.Indirect(reflect.ValueOf(receipt))
	receiptType := receiptVal.Type()
	for i := 0; i < receiptType.NumField(); i++ {
		field := receiptVal.Field(i)
		fieldType := receiptType.Field(i)
		switch fieldType.Name {
		case "Retailer":
			retailerName := field.Interface().(string)
			if retailerName == "" {
				return errors.New("retailer name should not be empty")
			}
		case "PurchaseDate":
			purchaseDate := field.Interface().(string)
			if purchaseDate == "" {
				return errors.New("purchase date should not be empty")
			}
		case "PurchaseTime":
			purchaseTime := field.Interface().(string)
			if purchaseTime == "" {
				return errors.New("purchase time should not be empty")
			}
		case "Total":
			total := field.Interface().(string)
			_, err := strconv.ParseFloat(total, 64)
			if total == "" {
				return errors.New("total should not be empty")
			} else if err != nil {
				return errors.New("please check the format of your total")
			}

		case "Items":
			item := field.Interface().([]Item)
			for _, val := range item {
				_, err := strconv.ParseFloat(val.Price, 64)
				if err != nil {
					return errors.New("item does not have valid price entry")
				}
			}
		default:
			fmt.Printf("%s is a new field\n", fieldType.Name)
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

func PointsCalculator(id string) (int, error) {
	receiptToCheck, ok := receiptsMap[id]
	// fmt.Println(receiptToCheck)
	if !ok {
		return -1, errors.New("id not exists")
	}
	receiptVal := reflect.ValueOf(receiptToCheck)
	receiptType := receiptVal.Type()
	points := 0
	for i := 0; i < receiptType.NumField(); i++ {
		field := receiptVal.Field(i)
		fieldType := receiptType.Field(i)
		switch fieldType.Name {
		case "Retailer":
			retailerName := field.Interface().(string)
			retailerPoints := ruleRetailerName(retailerName)
			points += retailerPoints
		case "PurchaseDate":
			fmt.Println(field.Interface())
		case "PurchaseTime":
			fmt.Println(field.Interface())
		case "Total":
			totalStr := field.Interface().(string)
			totalVal, _ := strconv.ParseFloat(totalStr, 64)
			points += ruleTotal(totalVal)
		case "Items":
			// fmt.Println(fieldType.Name)
			fmt.Println(field.Interface())
		default:
			fmt.Println("No Points!")
		}

	}
	// fmt.Println(receiptToCheck)
	return points, nil

}

// One point for every alphanumeric character in the retailer name.
func ruleRetailerName(name string) int {
	nameLength := len(name)
	fmt.Printf("%d points from retailer name: %s\n", nameLength, name)
	return nameLength
}

// 50 points if the total is a round dollar amount with no cents.
// 25 points if the total is a multiple of 0.25.
func ruleTotal(total float64) int {
	pointReturn := 0
	if total == math.Round(total) {
		fmt.Println("50 points for total is a round dollar amount")
		pointReturn += 50
	}
	if math.Mod(total*100, 25) == 0 {
		fmt.Println("25 points for total is a multiple of 0.25")
		pointReturn += 25
	}
	return pointReturn
}

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
