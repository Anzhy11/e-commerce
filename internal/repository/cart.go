package repository

import (
	"github.com/anzhy11/go-e-commerce/internal/models"
	"gorm.io/gorm"
)

type CartRepositoryInterface interface {
	CreateCart(cart *models.Cart) error
	GetCartByUserID(userID uint) (*models.Cart, error)
	UpdateCart(cart *models.Cart) error
	DeleteCart(cart *models.Cart) error
}

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepositoryInterface {
	return &CartRepository{db: db}
}

// Cart
func (r *CartRepository) CreateCart(cart *models.Cart) error {
	return r.db.Create(cart).Error
}

func (r *CartRepository) GetCartByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	if err := r.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepository) UpdateCart(cart *models.Cart) error {
	return r.db.Save(cart).Error
}

func (r *CartRepository) DeleteCart(cart *models.Cart) error {
	return r.db.Delete(cart).Error
}
