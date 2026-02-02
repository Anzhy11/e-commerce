package orderHandler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	orderService "github.com/anzhy11/go-e-commerce/internal/services/orders"
	"github.com/anzhy11/go-e-commerce/internal/utils"
	"gorm.io/gorm"
)

type OrderHandlerInterface interface {
	CreateOrder(c *gin.Context)
	GetOrders(c *gin.Context)
	GetOrder(c *gin.Context)
}

type orderHandler struct {
	orderService orderService.OrderServiceInterface
}

func New(db *gorm.DB) OrderHandlerInterface {
	return &orderHandler{
		orderService: orderService.New(db),
	}
}

func (h *orderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	orderResponse, err := h.orderService.CreateOrder(userID)
	if err != nil {
		utils.InternalServerError(c, "failed to create order", err)
		return
	}

	utils.SuccessResponse(c, "Order created successfully", orderResponse)
}

func (h *orderHandler) GetOrders(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, meta, err := h.orderService.GetOrders(userID, page, limit)
	if err != nil {
		utils.InternalServerError(c, "failed to get orders", err)
		return
	}

	utils.SuccessResponse(c, "Orders fetched successfully", gin.H{
		"orders": orders,
		"meta":   meta,
	})
}

func (h *orderHandler) GetOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid order ID", err)
		return
	}

	orderResponse, err := h.orderService.GetOrder(userID, uint(orderID))
	if err != nil {
		utils.InternalServerError(c, "failed to get order", err)
		return
	}

	utils.SuccessResponse(c, "Order fetched successfully", orderResponse)
}
