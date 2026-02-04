package orderHandler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	_ "github.com/anzhy11/go-e-commerce/internal/dto"
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

// @Summary Create order
// @Description Create order
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Success 201 {object} utils.Response{data=dto.OrderResponse} "Order created successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /orders [post]
func (h *orderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	orderResponse, err := h.orderService.CreateOrder(userID)
	if err != nil {
		utils.InternalServerError(c, "failed to create order", err)
		return
	}

	utils.SuccessResponse(c, "Order created successfully", orderResponse)
}

// @Summary Get orders
// @Description Get orders
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.OrderResponse} "Orders fetched successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /orders [get]
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

// @Summary Get order
// @Description Get order
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param id path uint true "Order ID"
// @Success 200 {object} utils.Response{data=dto.OrderResponse} "Order fetched successfully"
// @Failure 404 {object} utils.Response "Order not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /orders/{id} [get]
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
