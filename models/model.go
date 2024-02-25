package models

import (
	"errors"
	"fmt"
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
			fmt.Printf("%d points from Retailer Name: %s\n", retailerPoints, retailerName)
			points += retailerPoints
		case "PurchaseDate":
			// fmt.Println(fieldType.Name)
			fmt.Println(field.Interface())
		case "PurchaseTime":
			// fmt.Println(fieldType.Name)
			fmt.Println(field.Interface())
		case "Total":
			// fmt.Println(fieldType.Name)

			// points += ruleRoundDollar(field.Interface().(string))
			fmt.Println(field.Interface())
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
	return len(name)
}

// 50 points if the total is a round dollar amount with no cents.
func ruleRoundDollar(total float64) int {

	// if totalInFloat == math.Round(totalInFloat) {
	// 	return 50
	// }
	return 0
}

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
