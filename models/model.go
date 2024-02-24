package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Item struct {
	shortDescription string
	price            string
}

type Receipt struct {
	retailer     string
	purchaseDate string
	purchaseTime string
	total        string
	items        []Item
}

var receiptsMap map[string]Receipt

func ReceiptsProcessor() (uuid.UUID, error) {
	id := uuid.New()
	fmt.Println(id.String())
	return id, nil
}
