package userHandler

import (
	"net/http"

	"github.com/anzhy11/go-e-commerce/internal/config"
	userService "github.com/anzhy11/go-e-commerce/internal/services/users"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type UserHandlerInterface interface {
	GetUser(c *gin.Context)
}

type userHandler struct {
	us  userService.UserServiceInterface
	log *zerolog.Logger
}

func New(cfg *config.Config, log *zerolog.Logger) UserHandlerInterface {
	us := userService.New(cfg, log)
	return &userHandler{
		us:  us,
		log: log,
	}
}

func (h *userHandler) GetUser(c *gin.Context) {
	h.log.Info().Msg("Hit user handler")
	user := h.us.GetUser("test123")
	c.JSON(http.StatusOK, gin.H{"success": true, "message": user})
}
