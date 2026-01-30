package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	FirstName string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name" gorm:"not null"`
	Phone     string         `json:"phone" gorm:"unique;not null"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"-"`
	Role      string         `json:"role" gorm:"default:customer"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relashionships
	RefreshTokens []RefreshToken `json:"-" gorm:"foreignKey:UserID;references:ID"`
	Orders        []Order        `json:"-" gorm:"foreignKey:UserID;references:ID"`
	Cart          Cart           `json:"-" gorm:"foreignKey:UserID;references:ID"`
}

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleCustomer UserRole = "customer"
)

type RefreshToken struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Token     string         `json:"token" gorm:"unique;not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relashionships
	User User `json:"-" gorm:"foreignKey:UserID;references:ID"`
}
