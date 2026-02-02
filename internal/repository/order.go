package repository

import (
	"github.com/anzhy11/go-e-commerce/internal/models"
	"gorm.io/gorm"
)

type OrderRepositoryInterface interface {
	GetOrders(userID uint, page, limit int) ([]models.Order, error)
	GetOrderById(userID, orderID uint) (*models.Order, error)
	CountOrders(userID uint) int64
	CreateOrderTX(data *models.Order, tx *gorm.DB) error
	GetOrderByIdTx(orderID uint, tx *gorm.DB) (*models.Order, error)
	GetCartByUserIDTx(userID uint, tx *gorm.DB) (*models.Cart, error)
	BeginTx() *gorm.DB
	CommitTx(tx *gorm.DB)
	RollbackTx(tx *gorm.DB)
	UpdateProductStockTx(product *models.Product, tx *gorm.DB) error
	ClearCartTx(cartID uint, tx *gorm.DB) error
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepositoryInterface {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) GetOrders(userID uint, offset, limit int) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Preload("OrderItems.Product.Category").
		Where("user_id = ?", userID).
		Offset(offset).Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetOrderById(userID, orderID uint) (*models.Order, error) {
	var order models.Order
	if err := r.db.Preload("OrderItems.Product.Category").Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) CountOrders(userID uint) int64 {
	count := int64(0)

	r.db.Model(&models.Order{}).
		Where("user_id = ?", userID).
		Count(&count)

	return count
}

// Transactional methods
func (r *OrderRepository) CreateOrderTX(data *models.Order, tx *gorm.DB) error {
	return tx.Create(&data).Error
}

func (r *OrderRepository) GetOrderByIdTx(orderID uint, tx *gorm.DB) (*models.Order, error) {
	var order models.Order
	if err := tx.Preload("OrderItems.Product.Category").Where("id = ?", orderID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetCartByUserIDTx(userID uint, tx *gorm.DB) (*models.Cart, error) {
	var cart models.Cart
	if err := tx.Preload("CartItems.Product").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *OrderRepository) UpdateProductStockTx(product *models.Product, tx *gorm.DB) error {
	return tx.Save(product).Error
}

func (r *OrderRepository) ClearCartTx(cartID uint, tx *gorm.DB) error {
	return tx.Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}

// Transaction helper
// begin
func (r *OrderRepository) BeginTx() *gorm.DB {
	return r.db.Begin()
}

// commit
func (r *OrderRepository) CommitTx(tx *gorm.DB) {
	tx.Commit()
}

// rollback
func (r *OrderRepository) RollbackTx(tx *gorm.DB) {
	tx.Rollback()
}
