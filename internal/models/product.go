package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"unique;not null"`
	Description string         `json:"description" gorm:"not null"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relashionships
	Products []Product `json:"-" gorm:"foreignKey:CategoryID;references:ID"`
}

type Product struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CategoryID  uint           `json:"category_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"unique;not null"`
	Description string         `json:"description" gorm:"not null"`
	Price       float64        `json:"price" gorm:"not null"`
	Stock       int            `json:"stock" gorm:"default:0"`
	SKU         string         `json:"sku" gorm:"unique;not null"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relashionships
	Category   Category       `json:"-" gorm:"foreignKey:CategoryID;references:ID"`
	Images     []ProductImage `json:"-" gorm:"foreignKey:ProductID;references:ID"`
	OrderItems []OrderItem    `json:"-" gorm:"foreignKey:ProductID;references:ID"`
	CartItems  []CartItem     `json:"-" gorm:"foreignKey:ProductID;references:ID"`
}

type ProductImage struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ProductID uint           `json:"product_id" gorm:"not null"`
	URL       string         `json:"url" gorm:"not null"`
	AltText   string         `json:"alt_text"`
	IsPrimary bool           `json:"is_primary" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relashionships
	Product Product `json:"-" gorm:"foreignKey:ProductID;references:ID"`
}
