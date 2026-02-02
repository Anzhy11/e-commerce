package cartService

import (
	"errors"

	"github.com/anzhy11/go-e-commerce/internal/dto"
	"github.com/anzhy11/go-e-commerce/internal/models"
	"github.com/anzhy11/go-e-commerce/internal/repository"
	"gorm.io/gorm"
)

type CartServiceInterface interface {
	GetCartByUserID(userID uint) (*dto.CartResponse, error)
	AddToCart(userID uint, cart *dto.AddToCartRequest) (*dto.CartResponse, error)
	UpdateCartItem(userID, cartItemID uint, cart *dto.UpdateCartRequest) (*dto.CartResponse, error)
	RemoveFromCart(userID uint, cartItemID uint) error
}

type cartService struct {
	db          *gorm.DB
	cartRepo    repository.CartRepositoryInterface
	productRepo repository.ProductRepositoryInterface
}

func New(db *gorm.DB) CartServiceInterface {
	return &cartService{
		db:          db,
		cartRepo:    repository.NewCartRepo(db),
		productRepo: repository.NewProductRepo(db),
	}
}

func (s *cartService) GetCartByUserID(userID uint) (*dto.CartResponse, error) {
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	return s.generateCartResponse(cart), nil
}

func (s *cartService) AddToCart(userID uint, data *dto.AddToCartRequest) (*dto.CartResponse, error) {
	product, err := s.productRepo.GetProductById(data.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if product.Stock < data.Quantity {
		return nil, errors.New("stock not enough")
	}

	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		createCart := models.Cart{UserID: userID}
		if createErr := s.cartRepo.CreateCart(&createCart); createErr != nil {
			return nil, createErr
		}
		cart = &createCart
	}

	cartItem, err := s.cartRepo.GetCartItemByCartID(cart.ID, data.ProductID)
	if err != nil {
		cartItem := models.CartItem{
			CartID:    cart.ID,
			ProductID: data.ProductID,
			Quantity:  data.Quantity,
		}

		if err := s.cartRepo.CreateCartItem(&cartItem); err != nil {
			return nil, err
		}
	} else {
		cartItem.Quantity += data.Quantity
		if cartItem.Quantity > product.Stock {
			return nil, errors.New("stock not enough")
		}
		if err := s.cartRepo.UpdateCartItem(cartItem); err != nil {
			return nil, err
		}
	}

	return s.generateCartResponse(cart), nil
}

func (s *cartService) UpdateCartItem(userID, cartItemID uint, data *dto.UpdateCartRequest) (*dto.CartResponse, error) {
	cartItem, err := s.cartRepo.GetUserCartItem(userID, cartItemID)
	if err != nil {
		return nil, errors.New("cart item not found")
	}

	product, err := s.productRepo.GetProductById(cartItem.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if product.Stock < data.Quantity {
		return nil, errors.New("stock not enough")
	}

	cartItem.Quantity = data.Quantity
	if err := s.cartRepo.UpdateCartItem(cartItem); err != nil {
		return nil, err
	}

	return s.GetCartByUserID(userID)
}

func (s *cartService) RemoveFromCart(userID, carttemID uint) error {
	if err := s.cartRepo.DeleteCartItem(userID, carttemID); err != nil {
		return err
	}
	return nil
}

// Helper
func (s *cartService) generateCartResponse(cart *models.Cart) *dto.CartResponse {
	cartItems := make([]dto.CartItemResponse, len(cart.CartItems))
	var total float64

	for i := range cart.CartItems {
		subtotal := float64(cart.CartItems[i].Quantity) * cart.CartItems[i].Product.Price
		total += subtotal
		cartItems[i] = dto.CartItemResponse{
			ID:       cart.CartItems[i].ID,
			Quantity: cart.CartItems[i].Quantity,
			Subtotal: subtotal,
			Product: dto.ProductResponse{
				ID:          cart.CartItems[i].Product.ID,
				CategoryID:  cart.CartItems[i].Product.CategoryID,
				Name:        cart.CartItems[i].Product.Name,
				Description: cart.CartItems[i].Product.Description,
				Price:       cart.CartItems[i].Product.Price,
				Stock:       cart.CartItems[i].Product.Stock,
				SKU:         cart.CartItems[i].Product.SKU,
				IsActive:    cart.CartItems[i].Product.IsActive,
				Category: dto.CategoryResponse{
					ID:          cart.CartItems[i].Product.Category.ID,
					Name:        cart.CartItems[i].Product.Category.Name,
					Description: cart.CartItems[i].Product.Category.Description,
					IsActive:    cart.CartItems[i].Product.Category.IsActive,
				},
			},
		}
	}

	return &dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CartItems: cartItems,
		Total:     total,
	}
}
