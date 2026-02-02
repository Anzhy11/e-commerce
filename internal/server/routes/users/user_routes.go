package userRoutes

import (
	userHandler "github.com/anzhy11/go-e-commerce/internal/server/handlers/users"
	"github.com/anzhy11/go-e-commerce/internal/server/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userRoutes struct {
	userHandler userHandler.UserHandlerInterface
}

func Setup(routeGroup *gin.RouterGroup, mdw *middlewares.Middlewares, db *gorm.DB) {
	ur := &userRoutes{
		userHandler: userHandler.New(db),
	}

	urg := routeGroup.Group("/users")
	urg.Use(mdw.Authorization())
	urg.GET("/profile", ur.userHandler.GetProfile)
	urg.PUT("/profile", ur.userHandler.UpdateProfile)
}
