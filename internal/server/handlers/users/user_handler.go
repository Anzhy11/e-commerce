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
	userService userService.UserServiceInterface
}

func New(db *gorm.DB) UserHandlerInterface {
	return &userHandler{
		userService: userService.New(db),
	}
}

// @Summary Get user profile
// @Description Get user profile
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.UserResponse} "User profile fetched successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /users/profile [get]
func (h *userHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		utils.NotFound(c, "User profile not found", err)
		return
	}

	utils.SuccessResponse(c, "User profile fetched successfully", user)
}

// @Summary Update user profile
// @Description Update user profile
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "User data"
// @Success 200 {object} utils.Response{data=dto.UserResponse} "User profile updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /users/profile [put]
func (h *userHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request", err)
		return
	}

	user, err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		utils.InternalServerError(c, "failed to update user profile", err)
		return
	}

	utils.SuccessResponse(c, "User profile updated successfully", user)
}
