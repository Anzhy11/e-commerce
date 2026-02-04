package authHandler

import (
	"github.com/anzhy11/go-e-commerce/internal/config"
	"github.com/anzhy11/go-e-commerce/internal/dto"
	"github.com/anzhy11/go-e-commerce/internal/events"
	authService "github.com/anzhy11/go-e-commerce/internal/services/auth"
	"github.com/anzhy11/go-e-commerce/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type AuthHandlerInterface interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
	Logout(c *gin.Context)
}

type authHandler struct {
	as authService.AuthServiceInterface
}

func New(db *gorm.DB, cfg *config.Config, log *zerolog.Logger, eventPub events.PublisherInterface) AuthHandlerInterface {
	return &authHandler{
		as: authService.New(db, cfg, log, eventPub),
	}
}

// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration data"
// @Success 201 {object} utils.Response{data=dto.AuthResponse} "User registered successfully"
// @Failure 400 {object} utils.Response "Invalid request data or user already exists"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /auth/register [post]
func (h *authHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request", err)
		return
	}

	resp, err := h.as.Register(&req)
	if err != nil {
		utils.InternalServerError(c, "failed to register user", err)
		return
	}

	utils.SuccessResponse(c, "user registered successfully", resp)
}

// @Summary Login user
// @Description Login user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login data"
// @Success 200 {object} utils.Response{data=dto.AuthResponse} "User logged in successfully"
// @Failure 400 {object} utils.Response "Invalid request data or user not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /auth/login [post]
func (h *authHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request", err)
		return
	}

	resp, err := h.as.Login(&req)
	if err != nil {
		utils.InternalServerError(c, "failed to login user", err)
		return
	}

	utils.SuccessResponse(c, "user logged in successfully", resp)
}

// @Summary Refresh token
// @Description Refresh access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token data"
// @Success 200 {object} utils.Response{data=dto.AuthResponse} "Token refreshed successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /auth/refresh [post]
func (h *authHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request", err)
		return
	}

	resp, err := h.as.RefreshToken(&req)
	if err != nil {
		utils.InternalServerError(c, "failed to refresh token", err)
		return
	}

	utils.SuccessResponse(c, "token refreshed successfully", resp)
}

// @Summary Logout user
// @Description Logout user by revoking refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token data"
// @Success 200 {object} utils.Response "User logged out successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /auth/logout [post]
func (h *authHandler) Logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request", err)
		return
	}

	err := h.as.Logout(req.RefreshToken)
	if err != nil {
		utils.InternalServerError(c, "failed to logout user", err)
		return
	}

	utils.SuccessResponse(c, "user logged out successfully", nil)
}
