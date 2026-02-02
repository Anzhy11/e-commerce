package cartHandler

import (
	"strconv"

	"github.com/anzhy11/go-e-commerce/internal/dto"
	cartService "github.com/anzhy11/go-e-commerce/internal/services/cart"
	"github.com/anzhy11/go-e-commerce/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartHandlerInterface interface {
	GetCartByUserID(c *gin.Context)
	AddToCart(c *gin.Context)
	UpdateCartItem(c *gin.Context)
	RemoveFromCart(c *gin.Context)
}

type cartHandler struct {
	cartService cartService.CartServiceInterface
}

func NewCartHandler(db *gorm.DB) CartHandlerInterface {
	return &cartHandler{
		cartService: cartService.New(db),
	}
}

func (h *cartHandler) GetCartByUserID(c *gin.Context) {
	userID := c.GetUint("user_id")

	cart, err := h.cartService.GetCartByUserID(userID)
	if err != nil {
		utils.NotFound(c, "cart not found", err)
		return
	}

	utils.SuccessResponse(c, "cart initialized successfully", cart)
}

func (h *cartHandler) AddToCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request", err)
		return
	}

	cart, err := h.cartService.AddToCart(userID, &req)
	if err != nil {
		utils.InternalServerError(c, "failed to add to cart", err)
		return
	}

	utils.SuccessResponse(c, "added to cart", cart)
}

func (h *cartHandler) UpdateCartItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "invalid item id", err)
		return
	}

	var req dto.UpdateCartRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request", err)
		return
	}

	cart, err := h.cartService.UpdateCartItem(userID, uint(id), &req)
	if err != nil {
		utils.InternalServerError(c, "failed to update cart item", err)
		return
	}

	utils.SuccessResponse(c, "cart item updated", cart)
}

func (h *cartHandler) RemoveFromCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "invalid request", err)
		return
	}

	if err := h.cartService.RemoveFromCart(userID, uint(id)); err != nil {
		utils.InternalServerError(c, "failed to remove from cart", err)
		return
	}

	utils.SuccessResponse(c, "removed from cart", nil)
}
