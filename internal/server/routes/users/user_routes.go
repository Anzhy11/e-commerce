package userRoutes

import (
	"github.com/anzhy11/go-e-commerce/internal/config"
	userHandler "github.com/anzhy11/go-e-commerce/internal/server/handlers/users"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type userRoutes struct {
	uh userHandler.UserHandlerInterface
}

func Setup(apiGroup *gin.RouterGroup, cfg *config.Config, log *zerolog.Logger) {
	uh := userHandler.New(cfg, log)

	ur := &userRoutes{
		uh: uh,
	}

	urg := apiGroup.Group("/users")
	urg.GET("/", ur.uh.GetUser)
}
