package userHandler

import (
	"github.com/anzhy11/go-e-commerce/internal/dto"
	userService "github.com/anzhy11/go-e-commerce/internal/services/users"
	"github.com/anzhy11/go-e-commerce/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandlerInterface interface {
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
}

type userHandler struct {
	us userService.UserServiceInterface
	db *gorm.DB
}

func New(db *gorm.DB) UserHandlerInterface {
	return &userHandler{
		db: db,
		us: userService.New(db),
	}
}

func (h *userHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := h.us.GetProfile(userID)
	if err != nil {
		utils.NotFound(c, "User profile not found", err)
		return
	}

	utils.SuccessResponse(c, "User profile fetched successfully", user)
}

func (h *userHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request", err)
		return
	}

	user, err := h.us.UpdateProfile(userID, &req)
	if err != nil {
		utils.InternalServerError(c, "failed to update user profile", err)
		return
	}

	utils.SuccessResponse(c, "User profile updated successfully", user)
}
