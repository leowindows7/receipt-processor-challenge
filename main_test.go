package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"sync"
	"testing"

	"examples.com/m/v2/endpoints"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var mockPointsStore sync.Map

func TestGetPoints(t *testing.T) {
	app := fiber.New()
	app.Get("/receipts/:id/points", endpoints.GetPoints)

	// Define test cases, using initializer func() to return expeceted ID/status code
	testCases := []struct {
		description  string
		initializer  func() string
		expectedCode int
	}{
		{
			description: "Valid receipt ID with points",
			initializer: func() string {
				// Setup: add points for a test receipt ID
				testReceiptID := uuid.New().String()
				mockPointsStore.Store(testReceiptID, 100)
				return testReceiptID
			},
			expectedCode: fiber.StatusBadRequest,
		},
		{
			description: "Valid UUID format but not in the store",
			initializer: func() string {
				// Return a valid but unknown UUID
				return uuid.New().String()
			},
			expectedCode: fiber.StatusBadRequest,
		},
	}

	// check each test case
	for _, eachCase := range testCases {
		t.Run(eachCase.description, func(t *testing.T) {
			receiptID := eachCase.initializer() // Execute setup function
			req := httptest.NewRequest("GET", fmt.Sprintf("/receipts/%s/points", receiptID), nil)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("%s: Test request failed: %v", eachCase.description, err)
			} else if resp.StatusCode != eachCase.expectedCode {
				t.Errorf("%s: Expected status code %d, got %d", eachCase.description, eachCase.expectedCode, resp.StatusCode)
			}

		})
	}
}

func TestProcessReceipt(t *testing.T) {
	app := fiber.New()
	app.Post("/receipts/process", endpoints.ProcessRecipts)

	// define test cases
	testCases := []struct {
		caseName       string
		receiptInput   map[string]interface{}
		expectedStatus int
	}{
		{
			caseName: "Valid receipt",
			receiptInput: map[string]interface{}{
				"retailer":     "Test",
				"purchaseDate": "2024-01-01",
				"purchaseTime": "17:00",
				"items": []map[string]string{
					{"shortDescription": "Item", "price": "10.00"},
				},
				"total": "10.00",
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			caseName: "Receipt with empty retailer name",
			receiptInput: map[string]interface{}{
				"retailer":     "",
				"purchaseDate": "2024-01-01",
				"purchaseTime": "17:00",
				"items": []map[string]string{
					{"shortDescription": "Item", "price": "10.00"},
				},
				"total": "10.00",
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			caseName: "Receipt with empty purchase date",
			receiptInput: map[string]interface{}{
				"retailer":     "Test Store",
				"purchaseDate": "",
				"purchaseTime": "15:00",
				"items": []map[string]interface{}{
					{"shortDescription": "Item 1", "price": "10.00"},
				},
				"total": "10.00",
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			caseName: "Receipt with invalid purchase date format",
			receiptInput: map[string]interface{}{
				"retailer":     "Test Store",
				"purchaseDate": "01-01-2024", // Invalid date format
				"purchaseTime": "15:00",
				"items": []map[string]interface{}{
					{"shortDescription": "Item 1", "price": "10.00"},
				},
				"total": "10.00",
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			caseName: "Receipt with invalid time format",
			receiptInput: map[string]interface{}{
				"retailer":     "Test Store",
				"purchaseDate": "2022-01-01",
				"purchaseTime": ":00",
				"items": []map[string]interface{}{
					{"shortDescription": "Item 1", "price": "10.00"},
				},
				"total": "10.00",
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			caseName: "Receipt with no items",
			receiptInput: map[string]interface{}{
				"retailer":     "Test Store",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "15:00",
				"total":        "10.00",
				"items":        []map[string]interface{}{},
			},

			expectedStatus: fiber.StatusBadRequest,
		},

		{
			caseName: "Receipt with invalid total format",
			receiptInput: map[string]interface{}{
				"retailer":     "Test Store",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "15:00",
				"items": []map[string]interface{}{
					{"shortDescription": "Item 1", "price": "10.00"},
				},
				"total": "1,0,00", // Invalid total format
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			caseName: "Invalid Recipet Payload",
			receiptInput: map[string]interface{}{
				"retailer":     "Test Store",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "15:00",
				"items":        []map[string]interface{}{}, // Empty items array
			},
			expectedStatus: fiber.StatusBadRequest,
		},
	}

	// iterate over test cases
	for _, eachCase := range testCases {
		t.Run(eachCase.caseName, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(eachCase.receiptInput)
			req := httptest.NewRequest("POST", "/receipts/process", bytes.NewReader(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("%s: Test request failed: %v", eachCase.caseName, err)
			}

			if resp.StatusCode != eachCase.expectedStatus {
				t.Errorf("%s: Expected status %d but got %d", eachCase.caseName, eachCase.expectedStatus, resp.StatusCode)
			}
		})
	}
}
