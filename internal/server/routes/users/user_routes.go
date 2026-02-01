package userRoutes

import (
	userHandler "github.com/anzhy11/go-e-commerce/internal/server/handlers/users"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userRoutes struct {
	uh userHandler.UserHandlerInterface
}

func Setup(routeGroup *gin.RouterGroup, db *gorm.DB) {
	uh := userHandler.New(db)

	ur := &userRoutes{
		uh: uh,
	}

	urg := routeGroup.Group("/users")
	urg.GET("/profile", ur.uh.GetProfile)
	urg.PUT("/profile", ur.uh.UpdateProfile)
}
