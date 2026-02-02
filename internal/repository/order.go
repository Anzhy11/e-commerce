package repository

import (
	"github.com/anzhy11/go-e-commerce/internal/models"
	"gorm.io/gorm"
)

type OrderRepositoryInterface interface {
	CreateOrder(data *models.Order) error
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepositoryInterface {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) CreateOrder(data *models.Order) error {
	return r.db.Create(&data).Error
}
