package orderService

import (
	"errors"
	"fmt"

	"github.com/anzhy11/go-e-commerce/internal/dto"
	"github.com/anzhy11/go-e-commerce/internal/models"
	"github.com/anzhy11/go-e-commerce/internal/repository"
	"github.com/anzhy11/go-e-commerce/internal/utils"
	"gorm.io/gorm"
)

type OrderServiceInterface interface {
	CreateOrder(userId uint) (*dto.OrderResponse, error)
	GetOrders(userId uint, page, limit int) ([]dto.OrderResponse, *utils.PaginatedMeta, error)
	GetOrder(userId, orderId uint) (*dto.OrderResponse, error)
}

type orderService struct {
	orderRepo repository.OrderRepositoryInterface
}

const dateFormat = "2006-01-02 15:04:05"

func New(db *gorm.DB) OrderServiceInterface {
	return &orderService{
		orderRepo: repository.NewOrderRepo(db),
	}
}

func (s *orderService) GetOrders(userId uint, page, limit int) ([]dto.OrderResponse, *utils.PaginatedMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	total := s.orderRepo.CountOrders(userId)

	orders, err := s.orderRepo.GetOrders(userId, offset, limit)
	if err != nil {
		return nil, nil, err
	}

	orderResponses := make([]dto.OrderResponse, len(orders))
	for i := range orders {
		order := &orders[i]
		orderResponses[i] = *s.generateOrderResponse(order)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	meta := utils.PaginatedMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return orderResponses, &meta, nil
}

func (s *orderService) GetOrder(userId, orderId uint) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.GetOrderById(userId, orderId)
	if err != nil {
		return nil, err
	}

	return s.generateOrderResponse(order), nil
}

func (s *orderService) CreateOrder(userId uint) (*dto.OrderResponse, error) {
	var orderResponse *dto.OrderResponse

	tx := s.orderRepo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			s.orderRepo.RollbackTx(tx)
			panic(r)
		}
	}()

	cartTx, createErr := s.orderRepo.GetCartByUserIDTx(userId, tx)
	if createErr != nil {
		return nil, createErr
	}

	if len(cartTx.CartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	var totalAmount float64
	var orderItems []models.OrderItem

	for i := range cartTx.CartItems {
		cartItem := &cartTx.CartItems[i]

		if cartItem.Product.Stock < cartItem.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %d", cartItem.Product.ID)
		}

		itemTotal := cartItem.Product.Price * float64(cartItem.Quantity)
		totalAmount += itemTotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Price:     cartItem.Product.Price,
		})

		cartItem.Product.Stock -= cartItem.Quantity
		if updateErr := s.orderRepo.UpdateProductStockTx(&cartItem.Product, tx); updateErr != nil {
			return nil, updateErr
		}
	}

	order := models.Order{
		UserID:      userId,
		Status:      models.OrderStatusPending,
		TotalAmount: totalAmount,
		OrderItems:  orderItems,
	}

	if err := s.orderRepo.CreateOrderTX(&order, tx); err != nil {
		return nil, err
	}

	if err := s.orderRepo.ClearCartTx(cartTx.ID, tx); err != nil {
		return nil, err
	}

	orderResponse, err := s.GetOrderByIdTx(order.ID, tx)
	if err != nil {
		return nil, err
	}

	s.orderRepo.CommitTx(tx)

	return orderResponse, nil
}

func (s *orderService) GetOrderByIdTx(orderID uint, tx *gorm.DB) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.GetOrderByIdTx(orderID, tx)
	if err != nil {
		return nil, err
	}
	return s.generateOrderResponse(order), nil
}

func (s *orderService) generateOrderResponse(order *models.Order) *dto.OrderResponse {
	orderItems := make([]dto.OrderItemResponse, len(order.OrderItems))
	for i := range order.OrderItems {
		orderItems[i] = dto.OrderItemResponse{
			ID:       order.OrderItems[i].ID,
			Quantity: order.OrderItems[i].Quantity,
			Price:    order.OrderItems[i].Price,
			Product: dto.ProductResponse{
				ID:          order.OrderItems[i].Product.ID,
				CategoryID:  order.OrderItems[i].Product.CategoryID,
				Name:        order.OrderItems[i].Product.Name,
				Description: order.OrderItems[i].Product.Description,
				Price:       order.OrderItems[i].Product.Price,
				Stock:       order.OrderItems[i].Product.Stock,
				SKU:         order.OrderItems[i].Product.SKU,
				IsActive:    order.OrderItems[i].Product.IsActive,
				Category: dto.CategoryResponse{
					ID:          order.OrderItems[i].Product.Category.ID,
					Name:        order.OrderItems[i].Product.Category.Name,
					Description: order.OrderItems[i].Product.Category.Description,
					IsActive:    order.OrderItems[i].Product.Category.IsActive,
				},
			},
		}
	}
	return &dto.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      string(order.Status),
		TotalAmount: order.TotalAmount,
		OrderItems:  orderItems,
		CreatedAt:   order.CreatedAt.Format(dateFormat),
	}
}
