package models

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

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

var (
	RegexRetailName  = regexp.MustCompile(`^[\w\s\-]+$`)
	RegexDescription = regexp.MustCompile(`^[\w\s\-]+$`)
	RegexPrice       = regexp.MustCompile(`^\d+\.\d{2}$`)
)

var receiptsMap map[string]Receipt = make(map[string]Receipt)

// payload validation, empty fields or entry with incorrect formats will be errored out
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
			} else if !RegexRetailName.MatchString(retailerName) {
				return errors.New("retailer name is not in correct format")
			}
		case "PurchaseDate":
			purchaseDate := field.Interface().(string)
			_, err := time.Parse("2006-01-02", purchaseDate)
			if err != nil {
				return errors.New("purchase date should be YYYY-MM-DD")
			}
		case "PurchaseTime":
			purchaseTime := field.Interface().(string)
			_, err := time.Parse("15:04", purchaseTime)
			if err != nil {
				return errors.New("purchase time should be HH:MM")
			}
		case "Total":
			total := field.Interface().(string)
			_, err := strconv.ParseFloat(total, 64)
			if err != nil {
				return errors.New("total should be entered as string in valid float format")
			} else if !RegexPrice.MatchString(total) {
				return errors.New("invalid total format")
			}

		case "Items":
			itemList := field.Interface().([]Item)
			if len(itemList) == 0 {
				return errors.New("please enter at least 1 item")
			}
			for _, item := range itemList {
				_, err := strconv.ParseFloat(item.Price, 64)
				if err != nil {
					return errors.New("at least one of items does not have valid price entry")
				} else if !RegexDescription.MatchString(item.ShortDescription) {
					return errors.New("at least one of items does not have valid short description")
				}
			}
		default:
			fmt.Printf("%s is a new field\n", fieldType.Name)
		}
	}

	return nil
}

// process receipts, validate payload in func checkPayload
// once a receipt is validated and processed an id will be generated accordingly
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
	fmt.Printf("receipt processed with id %s total %d receipts on file\n", idStr, len(receiptsMap))
	return idStr, nil
}

// calculate reward points. by this steps, all receipts stored in map should already contain valid entries
// rule functions is assigned to each field to calculate eligible reward points
func PointsCalculator(id string) (int, error) {
	receiptToCheck, ok := receiptsMap[id]
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
			purchaseDateStr := field.Interface().(string)
			purchaseDateVal, _ := time.Parse("2006-01-02", purchaseDateStr)
			purchaseDatePoints := rulePurchaseDate(purchaseDateVal)
			points += purchaseDatePoints
		case "PurchaseTime":
			purchaseTimeStr := field.Interface().(string)
			purchaseTimeVal, _ := time.Parse("15:04", purchaseTimeStr)
			purchaseTimePoints := ruleTimeOfPurchase(purchaseTimeVal)
			points += purchaseTimePoints
		case "Total":
			totalStr := field.Interface().(string)
			totalVal, _ := strconv.ParseFloat(totalStr, 64)
			points += ruleTotal(totalVal)
		case "Items":
			itemList := field.Interface().([]Item)
			points += ruleItemPurchased(itemList)
		default:
			fmt.Println("No Points!")
		}

	}
	fmt.Printf("%d points rewarded!\n", points)
	return points, nil
}

/*
	Begining of rule functions, 7 rules are implemented in following 5 rule functions
	1. One point for every alphanumeric character in the retailer name. (func ruleRetailerName)
	2. 50 points if the total is a round dollar amount with no cents. (func ruleTotal)
	3. 25 points if the total is a multiple of 0.25. (func ruleTotal)
	4. 5 points for every two items on the receipt. (func ruleItemPurchased)
	5. If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and
	   round up to the nearest integer. The result is the number of points earned. (func ruleItemPurchased)
	6. 6 points if the day in the purchase date is odd. (func rulePurchaseDate)
	7. 10 points if the time of purchase is after 2:00pm and before 4:00pm. (func ruleTimeOfPurchase)

*/

// One point for every alphanumeric character in the retailer name.
func ruleRetailerName(name string) int {
	nameRegExp := regexp.MustCompile(`[a-zA-Z0-9]+`)
	nameCompiled := nameRegExp.ReplaceAllString(name, "")
	pointsName := len(name) - len(nameCompiled)
	fmt.Printf("%d points for retailer name %s\n", pointsName, name)
	return pointsName
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

// 6 points if the day in the purchase date is odd.
func rulePurchaseDate(purchaseDate time.Time) int {
	if purchaseDate.Day()%2 == 1 {
		fmt.Println("6 points for day in purchase date is odd")
		return 6
	}
	fmt.Println("0 points for day in purchase date is even")
	return 0
}

// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
func ruleTimeOfPurchase(purchaseTime time.Time) int {
	if purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
		fmt.Println("10 points for time of purchase between 14:00 and 16:00")
		return 10
	}
	fmt.Println("0 points for time of purchase not within 14:00 to 16:00")
	return 0
}

// 5 points for every two items on the receipt.
// If the trimmed length of the item description is a multiple of 3,
// multiply the price by 0.2 and round up to the nearest integer.
// The result is the number of points earned.
func ruleItemPurchased(itemList []Item) int {

	totalNumOfItem := len(itemList)
	pointsTotalItem := totalNumOfItem / 2 * 5
	fmt.Printf("%d points for purchasing %d item(s)\n", pointsTotalItem, totalNumOfItem)
	pointsTrimmedLength := 0
	for _, item := range itemList {
		trimmedDescription := strings.TrimSpace(item.ShortDescription)
		trimmedDescriptionLength := len(trimmedDescription)
		if trimmedDescriptionLength%3 == 0 {
			itemPrice, _ := strconv.ParseFloat(item.Price, 64)
			pointsItemPrice := int(math.Ceil(itemPrice * 0.2))
			fmt.Printf("%d points for purchasing %s at $%s \n", pointsItemPrice, trimmedDescription, item.Price)
			pointsTrimmedLength += pointsItemPrice
		}
	}
	return pointsTotalItem + pointsTrimmedLength
}

func hashItem(itemToHash Item) string {

	s := []string{itemToHash.ShortDescription, itemToHash.Price}
	return strings.Join(s, ",")
}

func ruleAllUniqueItems(itemList []Item) int {
	// introduce new map to store the item
	// check whether the item isfound, if found -> not uniqure, no points
	var tmpMap map[string]Item = make(map[string]Item)
	// hashfunc for receiptVal
	for _, item := range itemList {
		hashVal := hashItem(item)
		if _, ok := tmpMap[hashVal]; ok {
			return 0
		}
		tmpMap[hashVal] = item
	}
	return 5 * len(itemList)
}
