package authRoutes

import (
	"github.com/anzhy11/go-e-commerce/internal/config"
	authHandler "github.com/anzhy11/go-e-commerce/internal/server/handlers/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type authRoutes struct {
	ah authHandler.AuthHandlerInterface
}

func Setup(apiGroup *gin.RouterGroup, db *gorm.DB, cfg *config.Config, log *zerolog.Logger) {
	ah := authHandler.New(db, cfg, log)

	ar := &authRoutes{
		ah: ah,
	}

	arg := apiGroup.Group("/auth")
	arg.POST("/register", ar.ah.Register)
	arg.POST("/login", ar.ah.Login)
	arg.POST("/refresh-token", ar.ah.RefreshToken)
	arg.POST("/logout", ar.ah.Logout)
}
