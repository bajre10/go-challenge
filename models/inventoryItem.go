package models

import (
	"gorm.io/gorm"
	"time"
)

type InventoryItem struct {
	gorm.Model  `swaggerignore:"true"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Quantity    uint      `json:"quantity" gorm:"not null"`
	LastRestock time.Time `json:"lastRestock" gorm:"default:null"`
	Restocks    []Restock `json:"-" gorm:"foreignKey:ItemID" swaggerignore:"true"`
}
