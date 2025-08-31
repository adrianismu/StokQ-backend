package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	SKU   string  `gorm:"unique;not null" json:"sku"`
	Name  string  `gorm:"not null" json:"name"`
	Stock int     `gorm:"default:0" json:"stock"`
	Price float64 `gorm:"not null" json:"price"`
}
