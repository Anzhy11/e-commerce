package authHandler

import (
	"github.com/anzhy11/go-e-commerce/internal/config"
	"github.com/anzhy11/go-e-commerce/internal/dto"
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

type AuthHandler struct {
	as authService.AuthServiceInterface
}

func New(db *gorm.DB, cfg *config.Config, log *zerolog.Logger) AuthHandlerInterface {
	as := authService.New(db, cfg, log)
	return &AuthHandler{
		as: as,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
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

func (h *AuthHandler) Login(c *gin.Context) {
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

func (h *AuthHandler) RefreshToken(c *gin.Context) {
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

func (h *AuthHandler) Logout(c *gin.Context) {
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
