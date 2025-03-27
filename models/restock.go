package models

import (
	"fmt"
	"gorm.io/gorm"
	"shortcut-challenge/api/appErrors"
)

const MIN_QUANTITY = 10
const MAX_QUANTITY = 1000

type Restock struct {
	gorm.Model    `swaggerignore:"true"`
	ItemID        uint          `json:"name" gorm:"not null"`
	InventoryItem InventoryItem `json:"inventoryItem" gorm:"foreignKey:ItemID"`
	Quantity      uint          `json:"quantity" gorm:"not null" swaggerignore:"true"`
}

func (r *Restock) BeforeSave(*gorm.DB) (err error) {
	if r.Quantity <= MIN_QUANTITY || r.Quantity >= MAX_QUANTITY {
		return &appErrors.DBError{
			Field:   "Quantity",
			Type:    appErrors.Validation,
			Message: fmt.Sprintf("Quantity must be between %d and %d", MIN_QUANTITY, MAX_QUANTITY),
		}
	}
	return nil
}
