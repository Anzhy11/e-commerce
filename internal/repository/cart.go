package repository

import (
	"github.com/anzhy11/go-e-commerce/internal/models"
	"gorm.io/gorm"
)

type CartRepositoryInterface interface {
	CreateCart(cart *models.Cart) error
	GetCartByUserID(userID uint) (*models.Cart, error)
	UpdateCart(cart *models.Cart) error
	GetCartItemByCartID(cartID, productID uint) (*models.CartItem, error)
	GetUserCartItem(userID, cartItemID uint) (*models.CartItem, error)
	CreateCartItem(cartItem *models.CartItem) error
	UpdateCartItem(cartItem *models.CartItem) error
	DeleteCartItem(userID, cartItemID uint) error
}

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepo(db *gorm.DB) CartRepositoryInterface {
	return &CartRepository{db: db}
}

// Cart
func (r *CartRepository) CreateCart(cart *models.Cart) error {
	return r.db.Create(&cart).Error
}

func (r *CartRepository) GetCartByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	if err := r.db.Preload("CartItems.Product.Category").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepository) UpdateCart(cart *models.Cart) error {
	return r.db.Save(&cart).Error
}

// CartItem
func (r *CartRepository) CreateCartItem(cartItem *models.CartItem) error {
	return r.db.Create(&cartItem).Error
}

func (r *CartRepository) GetCartItemByCartID(cartID, productID uint) (*models.CartItem, error) {
	var cartItem models.CartItem
	if err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&cartItem).Error; err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func (r *CartRepository) GetUserCartItem(userID, cartItemID uint) (*models.CartItem, error) {
	var cartItem models.CartItem
	if err := r.db.Joins("JOIN carts ON cart_items.cart_id = carts.id").
		Where("carts.user_id = ? AND cart_items.id = ?", userID, cartItemID).
		First(&cartItem).Error; err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func (r *CartRepository) UpdateCartItem(cartItem *models.CartItem) error {
	return r.db.Save(&cartItem).Error
}

func (r *CartRepository) DeleteCartItem(userID, cartItemID uint) error {
	var cartItem models.CartItem
	if err := r.db.Joins("JOIN carts ON cart_items.cart_id = carts.id").
		Where("carts.user_id = ? AND cart_items.id = ?", userID, cartItemID).
		First(&cartItem).Error; err != nil {
		return err
	}
	return r.db.Delete(&cartItem).Error
}
