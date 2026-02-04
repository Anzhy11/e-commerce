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

// @Summary Get cart by user ID
// @Description Get cart by user ID
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Param user_id path uint true "User ID"
// @Success 200 {object} utils.Response{data=dto.CartResponse} "Cart retrieved successfully"
// @Failure 404 {object} utils.Response "Cart not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /cart/{user_id} [get]
func (h *cartHandler) GetCartByUserID(c *gin.Context) {
	userID := c.GetUint("user_id")

	cart, err := h.cartService.GetCartByUserID(userID)
	if err != nil {
		utils.NotFound(c, "cart not found", err)
		return
	}

	utils.SuccessResponse(c, "cart initialized successfully", cart)
}

// @Summary Add item to cart
// @Description Add item to cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AddToCartRequest true "Cart data"
// @Success 201 {object} utils.Response{data=dto.CartResponse} "Item added to cart successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /cart [post]
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

// @Summary Update cart item
// @Description Update cart item
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "Cart item ID"
// @Param request body dto.UpdateCartRequest true "Cart data"
// @Success 200 {object} utils.Response{data=dto.CartResponse} "Cart item updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /cart/{id} [put]
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

// @Summary Remove item from cart
// @Description Remove item from cart
// @Tags Cart
// @Security BearerAuth
// @Param id path uint true "Cart item ID"
// @Success 200 {object} utils.Response "Item removed from cart successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /cart/{id} [delete]
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
